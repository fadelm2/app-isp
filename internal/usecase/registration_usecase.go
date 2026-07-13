package usecase

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"golang-clean-architecture/internal/entity"
	"golang-clean-architecture/internal/gateway/mikrotik"
	"golang-clean-architecture/internal/gateway/notification"
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/internal/model/converter"
	"golang-clean-architecture/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type RegistrationUseCase struct {
	DB                     *gorm.DB
	Log                    *logrus.Logger
	Validate               *validator.Validate
	RegistrationRepository *repository.RegistrationRepository
	PackageRepository      *repository.PackageRepository
	UserRepository         *repository.UserRepository
	CustomerRepository     *repository.CustomerRepository
	RadiusRepository       *repository.RadiusRepository
	RouterRepository       *repository.RouterRepository
	InvoiceRepository      *repository.InvoiceRepository
	MikrotikClient         *mikrotik.MikrotikClient
	NotificationClient     *notification.NotificationClient
}

func NewRegistrationUseCase(
	db *gorm.DB,
	log *logrus.Logger,
	validate *validator.Validate,
	regRepo *repository.RegistrationRepository,
	pkgRepo *repository.PackageRepository,
	userRepo *repository.UserRepository,
	custRepo *repository.CustomerRepository,
	radRepo *repository.RadiusRepository,
	routerRepo *repository.RouterRepository,
	invRepo *repository.InvoiceRepository,
	mkClient *mikrotik.MikrotikClient,
	notifClient *notification.NotificationClient,
) *RegistrationUseCase {
	return &RegistrationUseCase{
		DB:                     db,
		Log:                    log,
		Validate:               validate,
		RegistrationRepository: regRepo,
		PackageRepository:      pkgRepo,
		UserRepository:         userRepo,
		CustomerRepository:     custRepo,
		RadiusRepository:       radRepo,
		RouterRepository:       routerRepo,
		InvoiceRepository:      invRepo,
		MikrotikClient:         mkClient,
		NotificationClient:     notifClient,
	}
}

func (c *RegistrationUseCase) Create(ctx context.Context, request *model.CreateRegistrationRequest) (model.RegistrationResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body: %+v", err)
		return model.RegistrationResponse{}, fiber.ErrBadRequest
	}

	pkg := new(entity.InternetPackage)
	if err := c.PackageRepository.FindById(tx, pkg, request.PackageID); err != nil {
		c.Log.Warnf("Package not found: %s", request.PackageID)
		return model.RegistrationResponse{}, fiber.ErrBadRequest
	}

	reg := &entity.Registration{
		ID:                  uuid.NewString(),
		FullName:            request.FullName,
		NIK:                 request.NIK,
		BirthPlace:          request.BirthPlace,
		BirthDate:           request.BirthDate,
		Gender:              request.Gender,
		Email:               request.Email,
		Phone:               request.Phone,
		InstallationAddress: request.InstallationAddress,
		BillingAddress:      request.BillingAddress,
		PackageID:           request.PackageID,
		Latitude:            request.Latitude,
		Longitude:           request.Longitude,
		Notes:               request.Notes,
		Status:              "pending",
		KtpPath:             request.KtpPath,
		SelfiePath:          request.SelfiePath,
		HousePath:           request.HousePath,
		InstallationPath:    request.InstallationPath,
		SupportingDocPath:   request.SupportingDocPath,
	}

	if err := c.RegistrationRepository.Create(tx, reg); err != nil {
		c.Log.Errorf("Failed to save registration: %+v", err)
		return model.RegistrationResponse{}, err
	}

	if err := tx.Commit().Error; err != nil {
		return model.RegistrationResponse{}, err
	}

	reg.Package = *pkg
	return converter.RegistrationToResponse(reg), nil
}

func (c *RegistrationUseCase) List(ctx context.Context) ([]model.RegistrationResponse, error) {
	tx := c.DB.WithContext(ctx)
	regs, err := c.RegistrationRepository.FindAll(tx)
	if err != nil {
		return nil, fiber.ErrInternalServerError
	}

	var responses []model.RegistrationResponse
	for _, r := range regs {
		responses = append(responses, converter.RegistrationToResponse(&r))
	}
	return responses, nil
}

func (c *RegistrationUseCase) UpdateStatus(ctx context.Context, request *model.UpdateRegistrationStatusRequest, adminUserID string) (model.RegistrationResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		return model.RegistrationResponse{}, fiber.ErrBadRequest
	}

	reg := new(entity.Registration)
	if err := c.RegistrationRepository.FindById(tx, reg, request.ID); err != nil {
		return model.RegistrationResponse{}, fiber.ErrNotFound
	}

	// Fetch Package
	pkg := new(entity.InternetPackage)
	if err := c.PackageRepository.FindById(tx, pkg, reg.PackageID); err != nil {
		return model.RegistrationResponse{}, fiber.ErrInternalServerError
	}
	reg.Package = *pkg

	oldStatus := reg.Status
	reg.Status = request.Status

	if err := c.RegistrationRepository.Update(tx, reg); err != nil {
		return model.RegistrationResponse{}, fiber.ErrInternalServerError
	}

	// If approved, trigger customer creation flow
	if request.Status == "approved" && oldStatus != "approved" {
		err := c.processApproval(tx, reg, adminUserID)
		if err != nil {
			c.Log.Errorf("Approval processing failed: %+v", err)
			return model.RegistrationResponse{}, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed processing customer activation: %v", err))
		}
	}

	if err := tx.Commit().Error; err != nil {
		return model.RegistrationResponse{}, fiber.ErrInternalServerError
	}

	return converter.RegistrationToResponse(reg), nil
}

func (c *RegistrationUseCase) processApproval(tx *gorm.DB, reg *entity.Registration, adminUserID string) error {
	// 1. Create a User for the Customer
	userPassword := fmt.Sprintf("Greenet%s", reg.NIK[len(reg.NIK)-4:]) // Last 4 digits of NIK
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	customerUserID := uuid.NewString()
	user := &entity.User{
		ID:          customerUserID,
		Password:    string(hashedPassword),
		Name:        reg.FullName,
		Email:       reg.Email,
		RoleID:      "2", // Customers
		CompanyName: "GREENET",
	}

	if err := c.UserRepository.Create(tx, user); err != nil {
		return err
	}

	// 2. Select first router if any
	var routers []entity.Router
	var assignedRouterID *string
	var router entity.Router
	c.DB.Find(&routers)
	if len(routers) > 0 {
		assignedRouterID = &routers[0].ID
		router = routers[0]
	}

	// Generate PPP / RADIUS Creds
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	customerCode := fmt.Sprintf("CUST-%05d", r.Intn(100000))
	pppUser := fmt.Sprintf("%s@greenet", reg.FullName)
	pppPass := fmt.Sprintf("pass%d", r.Intn(90000)+10000)

	customer := &entity.Customer{
		ID:             customerCode,
		RegistrationID: &reg.ID,
		UserID:         customerUserID,
		Status:         "suspended", // Suspended until initial invoice is paid
		PackageID:      reg.PackageID,
		RouterID:       assignedRouterID,
		PppUsername:    pppUser,
		PppPassword:    pppPass,
		RadiusUsername: pppUser,
		RadiusPassword: pppPass,
		DueDateDay:     10,
	}

	if err := c.CustomerRepository.Create(tx, customer); err != nil {
		return err
	}

	// 3. Create initial customer history
	history := &entity.CustomerHistory{
		ID:         uuid.NewString(),
		CustomerID: customer.ID,
		Action:     "register",
		Notes:      "Customer created from approved registration form.",
		CreatedBy:  adminUserID,
	}
	if err := tx.Create(history).Error; err != nil {
		return err
	}

	// 4. RADIUS tables setup (but check entries set as disabled/unauthenticated by default because status is suspended)
	radCheck := &entity.RadCheck{
		Username:  pppUser,
		Attribute: "Cleartext-Password",
		Op:        ":=",
		Value:     pppPass,
	}
	if err := c.RadiusRepository.CreateOrUpdateCheck(tx, radCheck); err != nil {
		return err
	}

	// Rate limit entry
	speedVal := fmt.Sprintf("%dM/%dM", reg.Package.SpeedMbps, reg.Package.SpeedMbps)
	radReply := &entity.RadReply{
		Username:  pppUser,
		Attribute: "Mikrotik-Rate-Limit",
		Op:        "=",
		Value:     speedVal,
	}
	if err := c.RadiusRepository.CreateOrUpdateReply(tx, radReply); err != nil {
		return err
	}

	// 5. MikroTik router provisioning
	if assignedRouterID != nil {
		go func() {
			err := c.MikrotikClient.CreatePPPoESecret(router.Host, router.Port, router.Username, router.Password, pppUser, pppPass, speedVal)
			if err != nil {
				c.Log.Errorf("Failed provisioning to Mikrotik: %v", err)
			}
			// Keep it disabled initially
			err = c.MikrotikClient.DisablePPPoESecret(router.Host, router.Port, router.Username, router.Password, pppUser)
			if err != nil {
				c.Log.Errorf("Failed disabling Mikrotik secret initially: %v", err)
			}
		}()
	}

	// 6. Generate initial invoice (package price + installation fee)
	now := time.Now()
	invoiceID := fmt.Sprintf("INV-%d%02d%s", now.Year(), int(now.Month()), customerCode[5:])
	subTotal := reg.Package.Price
	tax := subTotal * reg.Package.TaxRate
	total := subTotal + tax + reg.Package.InstallationFee

	invoice := &entity.Invoice{
		ID:              invoiceID,
		CustomerID:      customer.ID,
		DueDate:         now.AddDate(0, 0, 7).UnixMilli(), // Due in 7 days
		PeriodMonth:     int(now.Month()),
		PeriodYear:      now.Year(),
		Amount:          subTotal,
		TaxAmount:       tax,
		InstallationFee: reg.Package.InstallationFee,
		TotalAmount:     total,
		Status:          "pending",
	}

	if err := c.InvoiceRepository.Create(tx, invoice); err != nil {
		return err
	}

	// Send notification
	notifMsg := fmt.Sprintf("Halo %s, pendaftaran Greenet Anda disetujui. Akun Anda: %s, Password: %s. Silakan bayar tagihan pertama sebesar Rp %d melalui portal untuk mengaktifkan internet.", reg.FullName, reg.Email, userPassword, int(total))
	go c.NotificationClient.SendWhatsApp(reg.Phone, notifMsg)

	return nil
}
