package usecase

import (
	"context"
	"errors"
	"fmt"

	"golang-clean-architecture/internal/entity"
	"golang-clean-architecture/internal/gateway/mikrotik"
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/internal/model/converter"
	"golang-clean-architecture/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CustomerUseCase struct {
	DB                        *gorm.DB
	Log                       *logrus.Logger
	Validate                  *validator.Validate
	CustomerRepository        *repository.CustomerRepository
	PackageRepository         *repository.PackageRepository
	RouterRepository          *repository.RouterRepository
	RadiusRepository          *repository.RadiusRepository
	CustomerHistoryRepository *repository.CustomerHistoryRepository
	MikrotikClient            *mikrotik.MikrotikClient
}

func NewCustomerUseCase(
	db *gorm.DB,
	log *logrus.Logger,
	validate *validator.Validate,
	custRepo *repository.CustomerRepository,
	pkgRepo *repository.PackageRepository,
	routerRepo *repository.RouterRepository,
	radRepo *repository.RadiusRepository,
	histRepo *repository.CustomerHistoryRepository,
	mkClient *mikrotik.MikrotikClient,
) *CustomerUseCase {
	return &CustomerUseCase{
		DB:                        db,
		Log:                       log,
		Validate:                  validate,
		CustomerRepository:        custRepo,
		PackageRepository:         pkgRepo,
		RouterRepository:          routerRepo,
		RadiusRepository:          radRepo,
		CustomerHistoryRepository: histRepo,
		MikrotikClient:            mkClient,
	}
}

func (c *CustomerUseCase) List(ctx context.Context) ([]model.CustomerResponse, error) {
	tx := c.DB.WithContext(ctx)
	customers, err := c.CustomerRepository.FindAll(tx)
	if err != nil {
		return nil, fiber.ErrInternalServerError
	}

	var responses []model.CustomerResponse
	for _, cust := range customers {
		responses = append(responses, converter.CustomerToResponse(&cust))
	}
	return responses, nil
}

func (c *CustomerUseCase) Get(ctx context.Context, id string) (model.CustomerResponse, error) {
	tx := c.DB.WithContext(ctx)
	customer := new(entity.Customer)
	if err := c.CustomerRepository.FindByIdWithDetails(tx, customer, id); err != nil {
		return model.CustomerResponse{}, fiber.ErrNotFound
	}
	return converter.CustomerToResponse(customer), nil
}

func (c *CustomerUseCase) GetHistory(ctx context.Context, id string) ([]model.CustomerHistoryResponse, error) {
	tx := c.DB.WithContext(ctx)
	histories, err := c.CustomerHistoryRepository.FindAllByCustomerId(tx, id)
	if err != nil {
		return nil, fiber.ErrInternalServerError
	}

	var responses []model.CustomerHistoryResponse
	for _, h := range histories {
		responses = append(responses, converter.CustomerHistoryToResponse(&h))
	}
	return responses, nil
}

func (c *CustomerUseCase) Update(ctx context.Context, request *model.UpdateCustomerRequest, adminUserID string) (model.CustomerResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		return model.CustomerResponse{}, fiber.ErrBadRequest
	}

	customer := new(entity.Customer)
	if err := c.CustomerRepository.FindByIdWithDetails(tx, customer, request.ID); err != nil {
		return model.CustomerResponse{}, fiber.ErrNotFound
	}

	notes := "Customer details updated: "
	if request.PackageID != "" && request.PackageID != customer.PackageID {
		pkg := new(entity.InternetPackage)
		if err := c.PackageRepository.FindById(tx, pkg, request.PackageID); err != nil {
			return model.CustomerResponse{}, fiber.NewError(fiber.StatusBadRequest, "Invalid internet package")
		}
		notes += fmt.Sprintf("Package changed from %s to %s. ", customer.Package.Name, pkg.Name)
		customer.PackageID = request.PackageID
		customer.Package = *pkg

		// Sync speed to RADIUS
		speedVal := fmt.Sprintf("%dM/%dM", pkg.SpeedMbps, pkg.SpeedMbps)
		radReply := &entity.RadReply{
			Username:  customer.RadiusUsername,
			Attribute: "Mikrotik-Rate-Limit",
			Op:        "=",
			Value:     speedVal,
		}
		if err := c.RadiusRepository.CreateOrUpdateReply(tx, radReply); err != nil {
			return model.CustomerResponse{}, err
		}
	}

	if request.RouterID != "" && (customer.RouterID == nil || request.RouterID != *customer.RouterID) {
		router := new(entity.Router)
		if err := c.RouterRepository.FindById(tx, router, request.RouterID); err != nil {
			return model.CustomerResponse{}, fiber.NewError(fiber.StatusBadRequest, "Invalid router")
		}
		notes += fmt.Sprintf("Router changed. ")
		customer.RouterID = &request.RouterID
		customer.Router = router
	}

	if request.DueDateDay != 0 {
		customer.DueDateDay = request.DueDateDay
	}

	if err := c.CustomerRepository.Update(tx, customer); err != nil {
		return model.CustomerResponse{}, fiber.ErrInternalServerError
	}

	// Save audit log
	history := &entity.CustomerHistory{
		ID:         uuid.NewString(),
		CustomerID: customer.ID,
		Action:     "update",
		Notes:      notes,
		CreatedBy:  adminUserID,
	}
	if err := tx.Create(history).Error; err != nil {
		return model.CustomerResponse{}, err
	}

	if err := tx.Commit().Error; err != nil {
		return model.CustomerResponse{}, fiber.ErrInternalServerError
	}

	return converter.CustomerToResponse(customer), nil
}

func (c *CustomerUseCase) Suspend(ctx context.Context, id string, notes string, adminUserID string) (model.CustomerResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	customer := new(entity.Customer)
	if err := c.CustomerRepository.FindByIdWithDetails(tx, customer, id); err != nil {
		return model.CustomerResponse{}, fiber.ErrNotFound
	}

	if customer.Status == "suspended" {
		return converter.CustomerToResponse(customer), nil
	}

	customer.Status = "suspended"
	if err := c.CustomerRepository.Update(tx, customer); err != nil {
		return model.CustomerResponse{}, fiber.ErrInternalServerError
	}

	// RADIUS: change check to fail auth
	radCheck := &entity.RadCheck{
		Username:  customer.RadiusUsername,
		Attribute: "Cleartext-Password",
		Op:        ":=",
		Value:     customer.RadiusPassword + "_SUSPENDED", // Invalidate password
	}
	if err := c.RadiusRepository.CreateOrUpdateCheck(tx, radCheck); err != nil {
		return model.CustomerResponse{}, err
	}

	// MikroTik: Disable PPP secret and drop active session
	if customer.RouterID != nil {
		go func() {
			r := customer.Router
			_ = c.MikrotikClient.DisablePPPoESecret(r.Host, r.Port, r.Username, r.Password, customer.PppUsername)
			_ = c.MikrotikClient.DisconnectActiveSession(r.Host, r.Port, r.Username, r.Password, customer.PppUsername)
		}()
	}

	// Audit Log
	history := &entity.CustomerHistory{
		ID:         uuid.NewString(),
		CustomerID: customer.ID,
		Action:     "suspend",
		Notes:      notes,
		CreatedBy:  adminUserID,
	}
	if err := tx.Create(history).Error; err != nil {
		return model.CustomerResponse{}, err
	}

	if err := tx.Commit().Error; err != nil {
		return model.CustomerResponse{}, fiber.ErrInternalServerError
	}

	return converter.CustomerToResponse(customer), nil
}

func (c *CustomerUseCase) Unsuspend(ctx context.Context, id string, notes string, adminUserID string) (model.CustomerResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	customer := new(entity.Customer)
	if err := c.CustomerRepository.FindByIdWithDetails(tx, customer, id); err != nil {
		return model.CustomerResponse{}, fiber.ErrNotFound
	}

	if customer.Status == "active" {
		return converter.CustomerToResponse(customer), nil
	}

	customer.Status = "active"
	if err := c.CustomerRepository.Update(tx, customer); err != nil {
		return model.CustomerResponse{}, fiber.ErrInternalServerError
	}

	// RADIUS restore password
	radCheck := &entity.RadCheck{
		Username:  customer.RadiusUsername,
		Attribute: "Cleartext-Password",
		Op:        ":=",
		Value:     customer.RadiusPassword,
	}
	if err := c.RadiusRepository.CreateOrUpdateCheck(tx, radCheck); err != nil {
		return model.CustomerResponse{}, err
	}

	// MikroTik Enable
	if customer.RouterID != nil {
		go func() {
			r := customer.Router
			_ = c.MikrotikClient.EnablePPPoESecret(r.Host, r.Port, r.Username, r.Password, customer.PppUsername)
		}()
	}

	// Audit Log
	history := &entity.CustomerHistory{
		ID:         uuid.NewString(),
		CustomerID: customer.ID,
		Action:     "unsuspend",
		Notes:      notes,
		CreatedBy:  adminUserID,
	}
	if err := tx.Create(history).Error; err != nil {
		return model.CustomerResponse{}, err
	}

	if err := tx.Commit().Error; err != nil {
		return model.CustomerResponse{}, fiber.ErrInternalServerError
	}

	return converter.CustomerToResponse(customer), nil
}

func (c *CustomerUseCase) Terminate(ctx context.Context, id string, notes string, adminUserID string) (model.CustomerResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	customer := new(entity.Customer)
	if err := c.CustomerRepository.FindByIdWithDetails(tx, customer, id); err != nil {
		return model.CustomerResponse{}, fiber.ErrNotFound
	}

	customer.Status = "terminated"
	if err := c.CustomerRepository.Update(tx, customer); err != nil {
		return model.CustomerResponse{}, fiber.ErrInternalServerError
	}

	// RADIUS delete check & reply
	_ = c.RadiusRepository.DeleteCheck(tx, customer.RadiusUsername)
	_ = c.RadiusRepository.DeleteReply(tx, customer.RadiusUsername)

	// MikroTik disable/delete PPP
	if customer.RouterID != nil {
		go func() {
			r := customer.Router
			_ = c.MikrotikClient.DisablePPPoESecret(r.Host, r.Port, r.Username, r.Password, customer.PppUsername)
			_ = c.MikrotikClient.DisconnectActiveSession(r.Host, r.Port, r.Username, r.Password, customer.PppUsername)
		}()
	}

	// Audit Log
	history := &entity.CustomerHistory{
		ID:         uuid.NewString(),
		CustomerID: customer.ID,
		Action:     "terminate",
		Notes:      notes,
		CreatedBy:  adminUserID,
	}
	if err := tx.Create(history).Error; err != nil {
		return model.CustomerResponse{}, err
	}

	// Deactivate user role/status
	user := new(entity.User)
	if err := c.DB.Where("id = ?", customer.UserID).First(user).Error; err == nil {
		user.RoleID = "50" // nonaktif
		c.DB.Save(user)
	}

	if err := tx.Commit().Error; err != nil {
		return model.CustomerResponse{}, fiber.ErrInternalServerError
	}

	return converter.CustomerToResponse(customer), nil
}

func (c *CustomerUseCase) FindByUserId(ctx context.Context, userId string) (model.CustomerResponse, error) {
	tx := c.DB.WithContext(ctx)
	customer := new(entity.Customer)
	if err := c.CustomerRepository.FindByUserId(tx, customer, userId); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.CustomerResponse{}, fiber.ErrNotFound
		}
		return model.CustomerResponse{}, fiber.ErrInternalServerError
	}
	return converter.CustomerToResponse(customer), nil
}
