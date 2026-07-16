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

func (c *CustomerUseCase) List(ctx context.Context, request *model.SearchCustomerRequest) ([]model.CustomerResponse, int64, error) {
	tx := c.DB.WithContext(ctx)

	if err := c.Validate.Struct(request); err != nil {
		return nil, 0, fiber.ErrBadRequest
	}

	customers, total, err := c.CustomerRepository.Search(tx, request)
	if err != nil {
		return nil, 0, fiber.ErrInternalServerError
	}

	var responses []model.CustomerResponse
	for _, cust := range customers {
		responses = append(responses, converter.CustomerToResponse(&cust))
	}
	return responses, total, nil
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

	oldPppUsername := customer.PppUsername
	oldRadiusUsername := customer.RadiusUsername
	var oldRouter *entity.Router
	if customer.Router != nil {
		r := *customer.Router
		oldRouter = &r
	}

	notes := "Customer details updated: "
	packageChanged := false
	if request.PackageID != "" && request.PackageID != customer.PackageID {
		pkg := new(entity.InternetPackage)
		if err := c.PackageRepository.FindById(tx, pkg, request.PackageID); err != nil {
			return model.CustomerResponse{}, fiber.NewError(fiber.StatusBadRequest, "Invalid internet package")
		}
		notes += fmt.Sprintf("Package changed from %s to %s. ", customer.Package.Name, pkg.Name)
		customer.PackageID = request.PackageID
		customer.Package = *pkg
		packageChanged = true
	}

	credsChanged := false
	if request.PppUsername != "" && request.PppUsername != customer.PppUsername {
		notes += fmt.Sprintf("PPP username changed from %s to %s. ", customer.PppUsername, request.PppUsername)
		customer.PppUsername = request.PppUsername
		customer.RadiusUsername = request.PppUsername // keep in sync
		credsChanged = true
	}

	if request.PppPassword != "" && request.PppPassword != customer.PppPassword {
		notes += fmt.Sprintf("PPP password changed. ")
		customer.PppPassword = request.PppPassword
		customer.RadiusPassword = request.PppPassword // keep in sync
		credsChanged = true
	}

	// 1. Sync to RADIUS if package or credentials changed
	if credsChanged || packageChanged {
		if oldRadiusUsername != customer.RadiusUsername {
			_ = c.RadiusRepository.DeleteCheck(tx, oldRadiusUsername)
			_ = c.RadiusRepository.DeleteReply(tx, oldRadiusUsername)
		}

		radPass := customer.RadiusPassword
		if customer.Status == "suspended" {
			radPass += "_SUSPENDED"
		}
		radCheck := &entity.RadCheck{
			Username:  customer.RadiusUsername,
			Attribute: "Cleartext-Password",
			Op:        ":=",
			Value:     radPass,
		}
		if err := c.RadiusRepository.CreateOrUpdateCheck(tx, radCheck); err != nil {
			return model.CustomerResponse{}, err
		}

		speedVal := fmt.Sprintf("%dM/%dM", customer.Package.SpeedMbps, customer.Package.SpeedMbps)
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

	routerChanged := false
	if request.RouterID != "" && (customer.RouterID == nil || request.RouterID != *customer.RouterID) {
		router := new(entity.Router)
		if err := c.RouterRepository.FindById(tx, router, request.RouterID); err != nil {
			return model.CustomerResponse{}, fiber.NewError(fiber.StatusBadRequest, "Invalid router")
		}
		notes += fmt.Sprintf("Router changed. ")
		customer.RouterID = &request.RouterID
		customer.Router = router
		routerChanged = true
	}

	// 2. Sync to MikroTik Router if router, credentials, or package changed
	if routerChanged || credsChanged || packageChanged {
		speedVal := fmt.Sprintf("%dM/%dM", customer.Package.SpeedMbps, customer.Package.SpeedMbps)

		// Delete secret on old router if router changed
		if routerChanged && oldRouter != nil {
			_ = c.MikrotikClient.DeletePPPoESecret(oldRouter.Host, oldRouter.Port, oldRouter.Username, oldRouter.Password, oldPppUsername)
			_ = c.MikrotikClient.DisconnectActiveSession(oldRouter.Host, oldRouter.Port, oldRouter.Username, oldRouter.Password, oldPppUsername)
		}

		// Delete old secret on current router if credentials/package changed but router didn't change
		if !routerChanged && (credsChanged || packageChanged) && oldRouter != nil {
			_ = c.MikrotikClient.DeletePPPoESecret(oldRouter.Host, oldRouter.Port, oldRouter.Username, oldRouter.Password, oldPppUsername)
			_ = c.MikrotikClient.DisconnectActiveSession(oldRouter.Host, oldRouter.Port, oldRouter.Username, oldRouter.Password, oldPppUsername)
		}

		// Provision secret on new/current router
		if customer.RouterID != nil {
			r := customer.Router
			err := c.MikrotikClient.CreatePPPoESecret(r.Host, r.Port, r.Username, r.Password, customer.PppUsername, customer.PppPassword, speedVal)
			if err != nil {
				return model.CustomerResponse{}, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("MikroTik sync failed: %v", err))
			}

			if customer.Status == "suspended" || customer.Status == "terminated" {
				_ = c.MikrotikClient.DisablePPPoESecret(r.Host, r.Port, r.Username, r.Password, customer.PppUsername)
				_ = c.MikrotikClient.DisconnectActiveSession(r.Host, r.Port, r.Username, r.Password, customer.PppUsername)
			}
		}
	}

	if request.DueDateDay != 0 {
		customer.DueDateDay = request.DueDateDay
	}

	if request.OdpNumber != "" {
		notes += fmt.Sprintf("ODP Number changed from %s to %s. ", customer.OdpNumber, request.OdpNumber)
		customer.OdpNumber = request.OdpNumber
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

func (c *CustomerUseCase) RecreatePPPoE(ctx context.Context, id string, adminUserID string) (model.CustomerResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	customer := new(entity.Customer)
	if err := c.CustomerRepository.FindByIdWithDetails(tx, customer, id); err != nil {
		return model.CustomerResponse{}, fiber.ErrNotFound
	}

	// 1. Recreate RADIUS records
	_ = c.RadiusRepository.DeleteCheck(tx, customer.RadiusUsername)
	_ = c.RadiusRepository.DeleteReply(tx, customer.RadiusUsername)

	radPass := customer.RadiusPassword
	if customer.Status == "suspended" {
		radPass += "_SUSPENDED"
	}
	radCheck := &entity.RadCheck{
		Username:  customer.RadiusUsername,
		Attribute: "Cleartext-Password",
		Op:        ":=",
		Value:     radPass,
	}
	if err := c.RadiusRepository.CreateOrUpdateCheck(tx, radCheck); err != nil {
		return model.CustomerResponse{}, err
	}

	// 2. Recreate MikroTik PPPoE Secret
	if customer.RouterID != nil {
		r := customer.Router
		speedVal := fmt.Sprintf("%dM/%dM", customer.Package.SpeedMbps, customer.Package.SpeedMbps)
		
		_ = c.MikrotikClient.DeletePPPoESecret(r.Host, r.Port, r.Username, r.Password, customer.PppUsername)
		_ = c.MikrotikClient.DisconnectActiveSession(r.Host, r.Port, r.Username, r.Password, customer.PppUsername)

		err := c.MikrotikClient.CreatePPPoESecret(r.Host, r.Port, r.Username, r.Password, customer.PppUsername, customer.PppPassword, speedVal)
		if err != nil {
			return model.CustomerResponse{}, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("MikroTik recreation failed: %v", err))
		}

		if customer.Status == "suspended" || customer.Status == "terminated" {
			_ = c.MikrotikClient.DisablePPPoESecret(r.Host, r.Port, r.Username, r.Password, customer.PppUsername)
			_ = c.MikrotikClient.DisconnectActiveSession(r.Host, r.Port, r.Username, r.Password, customer.PppUsername)
		}
	}

	// Audit Log
	history := &entity.CustomerHistory{
		ID:         uuid.NewString(),
		CustomerID: customer.ID,
		Action:     "recreate_pppoe",
		Notes:      "Recreated PPPoE secret on MikroTik and RADIUS settings.",
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
