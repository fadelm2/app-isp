package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"golang-clean-architecture/internal/entity"
	"golang-clean-architecture/internal/gateway/midtrans"
	"golang-clean-architecture/internal/gateway/notification"
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/internal/model/converter"
	"golang-clean-architecture/internal/repository"
	"golang-clean-architecture/internal/util"

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
	NotificationClient *notification.NotificationClient
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
	notifClient *notification.NotificationClient,
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
		NotificationClient: notifClient,
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

func (c *InvoiceUseCase) ListPublicCustomerInvoices(ctx context.Context, customerID string) ([]model.InvoiceResponse, error) {
	tx := c.DB.WithContext(ctx)

	customer := new(entity.Customer)
	if err := c.CustomerRepository.FindById(tx, customer, customerID); err != nil {
		c.Log.Warnf("Public invoice query failed: customer %s not found", customerID)
		return nil, fiber.ErrNotFound
	}

	var invoices []entity.Invoice
	err := tx.Preload("Customer.User").
		Preload("Customer.Package").
		Where("customer_id = ? AND (status = ? OR status = ?)", customerID, "pending", "owed").
		Find(&invoices).Error

	if err != nil {
		c.Log.Errorf("Failed to list public invoices: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	var responses []model.InvoiceResponse
	for _, inv := range invoices {
		responses = append(responses, converter.InvoiceToResponse(&inv))
	}
	return responses, nil
}

func (c *InvoiceUseCase) AutoGenerateInvoices(ctx context.Context) error {
	tx := c.DB.WithContext(ctx)

	// Load all active customers
	var customers []entity.Customer
	if err := tx.Preload("Package").Where("status = ?", "active").Find(&customers).Error; err != nil {
		c.Log.Errorf("Failed to find active customers: %v", err)
		return err
	}

	now := time.Now()
	for _, cust := range customers {
		// 1. Check if invoice for current month and year exists
		var count int64
		tx.Model(&entity.Invoice{}).Where("customer_id = ? AND period_month = ? AND period_year = ?", cust.ID, int(now.Month()), now.Year()).Count(&count)

		if count == 0 {
			// Generate invoice for current period
			_, err := c.Create(ctx, &model.CreateInvoiceRequest{
				CustomerID:  cust.ID,
				PeriodMonth: int(now.Month()),
				PeriodYear:  now.Year(),
			})
			if err != nil {
				c.Log.Errorf("Failed to auto-generate current invoice for customer %s: %v", cust.ID, err)
			} else {
				c.Log.Infof("Auto-generated current period invoice for customer %s", cust.ID)
			}
			continue // move to next customer
		}

		// 2. If current invoice exists, check if due date is within 5 days (or past due date)
		var currentInvoice entity.Invoice
		err := tx.Where("customer_id = ? AND period_month = ? AND period_year = ?", cust.ID, int(now.Month()), now.Year()).First(&currentInvoice).Error
		if err == nil {
			dueDate := time.UnixMilli(currentInvoice.DueDate)
			// 5 days before due date means: now is >= due_date - 5 days
			fiveDaysBefore := dueDate.AddDate(0, 0, -5)

			if now.After(fiveDaysBefore) || now.Equal(fiveDaysBefore) {
				// Check if invoice for the NEXT month already exists
				nextMonth := int(now.Month()) + 1
				nextYear := now.Year()
				if nextMonth > 12 {
					nextMonth = 1
					nextYear = now.Year() + 1
				}

				var nextCount int64
				tx.Model(&entity.Invoice{}).Where("customer_id = ? AND period_month = ? AND period_year = ?", cust.ID, nextMonth, nextYear).Count(&nextCount)

				if nextCount == 0 {
					_, err := c.Create(ctx, &model.CreateInvoiceRequest{
						CustomerID:  cust.ID,
						PeriodMonth: nextMonth,
						PeriodYear:  nextYear,
					})
					if err != nil {
						c.Log.Errorf("Failed to auto-generate next invoice for customer %s: %v", cust.ID, err)
					} else {
						c.Log.Infof("Auto-generated next period invoice for customer %s", cust.ID)
					}
				}
			}
		}
	}
	return nil
}

func (c *InvoiceUseCase) SendLatestInvoiceWhatsApp(ctx context.Context, customerId string, baseURL string) error {
	tx := c.DB.WithContext(ctx)

	customer := new(entity.Customer)
	if err := c.CustomerRepository.FindByIdWithDetails(tx, customer, customerId); err != nil {
		return fiber.ErrNotFound
	}

	phone := ""
	if customer.Registration != nil && customer.Registration.Phone != "" {
		phone = customer.Registration.Phone
	} else {
		var contacts []entity.Contact
		if err := tx.Where("user_id = ?", customer.UserID).Find(&contacts).Error; err == nil && len(contacts) > 0 {
			phone = contacts[0].Phone
		}
	}

	if phone == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Customer phone number not found")
	}

	var latestInvoice entity.Invoice
	if err := tx.Where("customer_id = ?", customer.ID).Order("created_at DESC").First(&latestInvoice).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "No invoice found for this customer")
	}

	ispName := util.GetISPName()
	pdfLink := fmt.Sprintf("%s/api/public/invoices/%s/pdf", baseURL, latestInvoice.ID)

	message := fmt.Sprintf("Halo %s, berikut tagihan internet %s Anda untuk periode %02d/%d sebesar Rp %s. Silakan bayar sebelum %s.\nLink PDF: %s",
		customer.User.Name,
		ispName,
		latestInvoice.PeriodMonth,
		latestInvoice.PeriodYear,
		formatPriceStr(latestInvoice.TotalAmount),
		time.UnixMilli(latestInvoice.DueDate).Format("02-01-2006"),
		pdfLink,
	)

	return c.NotificationClient.SendWhatsApp(phone, message)
}

func (c *InvoiceUseCase) SendLatestInvoiceEmail(ctx context.Context, customerId string, baseURL string) error {
	tx := c.DB.WithContext(ctx)

	customer := new(entity.Customer)
	if err := c.CustomerRepository.FindByIdWithDetails(tx, customer, customerId); err != nil {
		return fiber.ErrNotFound
	}

	email := customer.User.Email
	if email == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Customer email not found")
	}

	var latestInvoice entity.Invoice
	if err := tx.Where("customer_id = ?", customer.ID).Order("created_at DESC").First(&latestInvoice).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "No invoice found for this customer")
	}

	ispName := util.GetISPName()
	pdfLink := fmt.Sprintf("%s/api/public/invoices/%s/pdf", baseURL, latestInvoice.ID)

	subject := fmt.Sprintf("Tagihan Internet %s - Invoice %s", ispName, latestInvoice.ID)
	body := fmt.Sprintf("Halo %s,\n\nTerima kasih telah menggunakan layanan internet %s.\n\nBerikut rincian tagihan Anda:\n- Invoice ID: %s\n- Periode: %02d/%d\n- Total Tagihan: Rp %s\n- Jatuh Tempo: %s\n\nAnda dapat mengunduh invoice PDF Anda di sini: %s\n\nSilakan lakukan pembayaran melalui portal pelanggan.\n\nSalam hangat,\n%s",
		customer.User.Name,
		ispName,
		latestInvoice.ID,
		latestInvoice.PeriodMonth,
		latestInvoice.PeriodYear,
		formatPriceStr(latestInvoice.TotalAmount),
		time.UnixMilli(latestInvoice.DueDate).Format("02-01-2006"),
		pdfLink,
		ispName,
	)

	return c.NotificationClient.SendEmail(email, subject, body)
}

func (c *InvoiceUseCase) SendInvoiceWhatsApp(ctx context.Context, invoiceId string, baseURL string) error {
	tx := c.DB.WithContext(ctx)

	invoice := new(entity.Invoice)
	if err := c.InvoiceRepository.FindByIdWithDetails(tx, invoice, invoiceId); err != nil {
		return fiber.ErrNotFound
	}

	phone := ""
	if invoice.Customer.Registration != nil && invoice.Customer.Registration.Phone != "" {
		phone = invoice.Customer.Registration.Phone
	} else {
		var contacts []entity.Contact
		if err := tx.Where("user_id = ?", invoice.Customer.UserID).Find(&contacts).Error; err == nil && len(contacts) > 0 {
			phone = contacts[0].Phone
		}
	}

	if phone == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Customer phone number not found")
	}

	ispName := util.GetISPName()
	pdfLink := fmt.Sprintf("%s/api/public/invoices/%s/pdf", baseURL, invoice.ID)

	message := fmt.Sprintf("Halo %s, berikut tagihan internet %s Anda untuk periode %02d/%d sebesar Rp %s. Silakan bayar sebelum %s.\nLink PDF: %s",
		invoice.Customer.User.Name,
		ispName,
		invoice.PeriodMonth,
		invoice.PeriodYear,
		formatPriceStr(invoice.TotalAmount),
		time.UnixMilli(invoice.DueDate).Format("02-01-2006"),
		pdfLink,
	)

	return c.NotificationClient.SendWhatsApp(phone, message)
}

func (c *InvoiceUseCase) SendInvoiceEmail(ctx context.Context, invoiceId string, baseURL string) error {
	tx := c.DB.WithContext(ctx)

	invoice := new(entity.Invoice)
	if err := c.InvoiceRepository.FindByIdWithDetails(tx, invoice, invoiceId); err != nil {
		return fiber.ErrNotFound
	}

	email := invoice.Customer.User.Email
	if email == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Customer email not found")
	}

	ispName := util.GetISPName()
	pdfLink := fmt.Sprintf("%s/api/public/invoices/%s/pdf", baseURL, invoice.ID)

	subject := fmt.Sprintf("Tagihan Internet %s - Invoice %s", ispName, invoice.ID)
	body := fmt.Sprintf("Halo %s,\n\nTerima kasih telah menggunakan layanan internet %s.\n\nBerikut rincian tagihan Anda:\n- Invoice ID: %s\n- Periode: %02d/%d\n- Total Tagihan: Rp %s\n- Jatuh Tempo: %s\n\nAnda dapat mengunduh invoice PDF Anda di sini: %s\n\nSilakan lakukan pembayaran melalui portal pelanggan.\n\nSalam hangat,\n%s",
		invoice.Customer.User.Name,
		ispName,
		invoice.ID,
		invoice.PeriodMonth,
		invoice.PeriodYear,
		formatPriceStr(invoice.TotalAmount),
		time.UnixMilli(invoice.DueDate).Format("02-01-2006"),
		pdfLink,
		ispName,
	)

	return c.NotificationClient.SendEmail(email, subject, body)
}

func formatPriceStr(val float64) string {
	str := fmt.Sprintf("%.0f", val)
	var result []string
	length := len(str)
	for i, ch := range str {
		result = append(result, string(ch))
		if (length-i-1)%3 == 0 && i != length-1 {
			result = append(result, ".")
		}
	}
	return strings.Join(result, "")
}

func (c *InvoiceUseCase) GetPDFData(ctx context.Context, id string) ([]byte, error) {
	tx := c.DB.WithContext(ctx)
	invoice := new(entity.Invoice)
	if err := c.InvoiceRepository.FindByIdWithDetails(tx, invoice, id); err != nil {
		return nil, fiber.ErrNotFound
	}

	ispName := util.GetISPName()
	return util.GenerateInvoicePDF(invoice, ispName)
}
