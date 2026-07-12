package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"golang-clean-architecture/internal/entity"
	"golang-clean-architecture/internal/gateway/midtrans"
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/internal/model/converter"
	"golang-clean-architecture/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type InvoiceUseCase struct {
	DB                 *gorm.DB
	Log                *logrus.Logger
	Validate           *validator.Validate
	InvoiceRepository  *repository.InvoiceRepository
	CustomerRepository *repository.CustomerRepository
	PaymentRepository  *repository.PaymentRepository
	MidtransClient     *midtrans.MidtransClient
	CustomerUseCase    *CustomerUseCase
}

func NewInvoiceUseCase(
	db *gorm.DB,
	log *logrus.Logger,
	validate *validator.Validate,
	invRepo *repository.InvoiceRepository,
	custRepo *repository.CustomerRepository,
	payRepo *repository.PaymentRepository,
	mtClient *midtrans.MidtransClient,
	custUseCase *CustomerUseCase,
) *InvoiceUseCase {
	return &InvoiceUseCase{
		DB:                 db,
		Log:                log,
		Validate:           validate,
		InvoiceRepository:  invRepo,
		CustomerRepository: custRepo,
		PaymentRepository:  payRepo,
		MidtransClient:     mtClient,
		CustomerUseCase:    custUseCase,
	}
}

func (c *InvoiceUseCase) Create(ctx context.Context, request *model.CreateInvoiceRequest) (model.InvoiceResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		return model.InvoiceResponse{}, fiber.ErrBadRequest
	}

	customer := new(entity.Customer)
	if err := c.CustomerRepository.FindByIdWithDetails(tx, customer, request.CustomerID); err != nil {
		return model.InvoiceResponse{}, fiber.ErrNotFound
	}

	invoiceID := fmt.Sprintf("INV-%d%02d%s", request.PeriodYear, request.PeriodMonth, customer.ID[5:])

	// Check if invoice for this period already exists
	var count int64
	tx.Model(&entity.Invoice{}).Where("customer_id = ? AND period_month = ? AND period_year = ?", customer.ID, request.PeriodMonth, request.PeriodYear).Count(&count)
	if count > 0 {
		return model.InvoiceResponse{}, fiber.NewError(fiber.StatusConflict, "Invoice for this period already exists")
	}

	subTotal := customer.Package.Price
	tax := subTotal * customer.Package.TaxRate
	total := subTotal + tax

	invoice := &entity.Invoice{
		ID:              invoiceID,
		CustomerID:      customer.ID,
		DueDate:         time.Date(request.PeriodYear, time.Month(request.PeriodMonth), customer.DueDateDay, 23, 59, 59, 0, time.Local).UnixMilli(),
		PeriodMonth:     request.PeriodMonth,
		PeriodYear:      request.PeriodYear,
		Amount:          subTotal,
		TaxAmount:       tax,
		InstallationFee: 0,
		TotalAmount:     total,
		Status:          "pending",
	}

	if err := c.InvoiceRepository.Create(tx, invoice); err != nil {
		return model.InvoiceResponse{}, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		return model.InvoiceResponse{}, fiber.ErrInternalServerError
	}

	invoice.Customer = *customer
	return converter.InvoiceToResponse(invoice), nil
}

func (c *InvoiceUseCase) GetSnapToken(ctx context.Context, invoiceID string) (string, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	invoice := new(entity.Invoice)
	if err := c.InvoiceRepository.FindByIdWithDetails(tx, invoice, invoiceID); err != nil {
		return "", fiber.ErrNotFound
	}

	if invoice.SnapToken != "" {
		return invoice.SnapToken, nil
	}

	token, err := c.MidtransClient.CreateSnapToken(invoice.ID, invoice.TotalAmount, invoice.Customer.User.Name, invoice.Customer.User.Email)
	if err != nil {
		c.Log.Errorf("Midtrans token creation failed: %v", err)
		return "", fiber.ErrInternalServerError
	}

	invoice.SnapToken = token
	if err := c.InvoiceRepository.Update(tx, invoice); err != nil {
		return "", fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		return "", fiber.ErrInternalServerError
	}

	return token, nil
}

func (c *InvoiceUseCase) ProcessWebhook(ctx context.Context, payload map[string]interface{}) error {
	orderID, ok := payload["order_id"].(string)
	if !ok {
		return errors.New("missing order_id")
	}

	statusCode, _ := payload["status_code"].(string)
	grossAmount, _ := payload["gross_amount"].(string)
	signatureKey, _ := payload["signature_key"].(string)

	// Verify signature
	if c.MidtransClient.ServerKey != "" {
		valid := c.MidtransClient.VerifySignature(orderID, statusCode, grossAmount, signatureKey)
		if !valid {
			c.Log.Warnf("Invalid Midtrans signature webhook for order: %s", orderID)
			return errors.New("invalid signature key")
		}
	}

	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	invoice := new(entity.Invoice)
	err := c.InvoiceRepository.FindByIdWithDetails(tx, invoice, orderID)
	if err != nil {
		c.Log.Warnf("Invoice not found for webhook order: %s", orderID)
		return nil // Return nil so we don't retry from Midtrans side
	}

	transactionStatus, _ := payload["transaction_status"].(string)
	paymentType, _ := payload["payment_type"].(string)
	transactionID, _ := payload["transaction_id"].(string)

	paidAmount := invoice.TotalAmount
	if gross, err := strconv.ParseFloat(grossAmount, 64); err == nil {
		paidAmount = gross
	}

	status := "pending"
	isPaid := false

	switch transactionStatus {
	case "capture", "settlement":
		status = "paid"
		isPaid = true
	case "pending":
		status = "pending"
	case "deny", "expire", "cancel":
		status = "cancelled"
	}

	invoice.Status = status
	if isPaid {
		now := time.Now().UnixMilli()
		invoice.PaidAt = &now
	}

	if err := c.InvoiceRepository.Update(tx, invoice); err != nil {
		return err
	}

	// Save payment record
	rawJSON, _ := json.Marshal(payload)
	payment := &entity.Payment{
		ID:            uuid.NewString(),
		InvoiceID:     invoice.ID,
		TransactionID: transactionID,
		PaymentType:   paymentType,
		PaidAmount:    paidAmount,
		Status:        transactionStatus,
		PaidAt:        time.Now().UnixMilli(),
		RawResponse:   string(rawJSON),
	}
	if err := c.PaymentRepository.Create(tx, payment); err != nil {
		return err
	}

	// Trigger automatic reactivation
	if isPaid {
		c.Log.Infof("Invoice %s PAID. Unsuspending customer %s.", invoice.ID, invoice.CustomerID)
		_, err = c.CustomerUseCase.Unsuspend(ctx, invoice.CustomerID, fmt.Sprintf("Auto unsuspended due to payment of Invoice %s", invoice.ID), "SYSTEM")
		if err != nil {
			c.Log.Errorf("Reactivation failed: %v", err)
			return err
		}
	}

	return tx.Commit().Error
}

func (c *InvoiceUseCase) List(ctx context.Context) ([]model.InvoiceResponse, error) {
	tx := c.DB.WithContext(ctx)
	invoices, err := c.InvoiceRepository.FindAll(tx)
	if err != nil {
		return nil, fiber.ErrInternalServerError
	}

	var responses []model.InvoiceResponse
	for _, inv := range invoices {
		responses = append(responses, converter.InvoiceToResponse(&inv))
	}
	return responses, nil
}

func (c *InvoiceUseCase) Get(ctx context.Context, id string) (model.InvoiceResponse, error) {
	tx := c.DB.WithContext(ctx)
	invoice := new(entity.Invoice)
	if err := c.InvoiceRepository.FindByIdWithDetails(tx, invoice, id); err != nil {
		return model.InvoiceResponse{}, fiber.ErrNotFound
	}
	return converter.InvoiceToResponse(invoice), nil
}

func (c *InvoiceUseCase) ListByCustomer(ctx context.Context, customerId string) ([]model.InvoiceResponse, error) {
	tx := c.DB.WithContext(ctx)
	invoices, err := c.InvoiceRepository.FindByCustomerId(tx, customerId)
	if err != nil {
		return nil, fiber.ErrInternalServerError
	}

	var responses []model.InvoiceResponse
	for _, inv := range invoices {
		responses = append(responses, converter.InvoiceToResponse(&inv))
	}
	return responses, nil
}
