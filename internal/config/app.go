package config

import (
	"golang-clean-architecture/internal/delivery/http"
	"golang-clean-architecture/internal/delivery/http/middleware"
	"golang-clean-architecture/internal/delivery/http/route"
	"golang-clean-architecture/internal/gateway/messaging"
	"golang-clean-architecture/internal/gateway/midtrans"
	"golang-clean-architecture/internal/gateway/mikrotik"
	"golang-clean-architecture/internal/gateway/notification"
	"golang-clean-architecture/internal/repository"
	"golang-clean-architecture/internal/usecase"
	"golang-clean-architecture/internal/util"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB        *gorm.DB
	App       *fiber.App
	Log       *logrus.Logger
	Validate  *validator.Validate
	Config    *viper.Viper
	Producer  *kafka.Producer
	SecretKey string
}

func Bootstrap(config *BootstrapConfig) {
	// setup repositories
	userRepository := repository.NewUserRepository(config.Log)
	contactRepository := repository.NewContactRepository(config.Log)
	addressRepository := repository.NewAddressRepository(config.Log)
	tokenSecretKey := config.SecretKey

	packageRepo := repository.NewPackageRepository(config.Log)
	regRepo := repository.NewRegistrationRepository(config.Log)
	custRepo := repository.NewCustomerRepository(config.Log)
	routerRepo := repository.NewRouterRepository(config.Log)
	invRepo := repository.NewInvoiceRepository(config.Log)
	payRepo := repository.NewPaymentRepository(config.Log)
	radRepo := repository.NewRadiusRepository(config.Log)
	histRepo := repository.NewCustomerHistoryRepository(config.Log)

	// setup gateways
	midtransServerKey := config.Config.GetString("midtrans.server_key")
	midtransIsProd := config.Config.GetBool("midtrans.is_production")

	mkClient := mikrotik.NewMikrotikClient(config.Log)
	mtClient := midtrans.NewMidtransClient(midtransServerKey, midtransIsProd, config.Log)
	notifClient := notification.NewNotificationClient(config.Log)

	//setup producer
	userProducer := messaging.NewUserProducer(config.Producer, config.Log)
	contactProducer := messaging.NewContactProducer(config.Producer, config.Log)
	addressProducer := messaging.NewAddressProducer(config.Producer, config.Log)

	tokenUtil := util.NewTokenUtil(tokenSecretKey)

	//setup use cases
	userUseCase := usecase.NewUserCase(config.DB, config.Log, config.Validate, userRepository, userProducer, tokenUtil)
	ContactUseCase := usecase.NewContactUseCase(config.DB, config.Log, config.Validate, contactRepository, contactProducer)
	addressUseCase := usecase.NewAddressUseCase(config.DB, config.Log, config.Validate, contactRepository, addressRepository, addressProducer)

	pkgUseCase := usecase.NewPackageUseCase(config.DB, config.Log, config.Validate, packageRepo)
	custUseCase := usecase.NewCustomerUseCase(config.DB, config.Log, config.Validate, custRepo, packageRepo, routerRepo, radRepo, histRepo, mkClient)
	invUseCase := usecase.NewInvoiceUseCase(config.DB, config.Log, config.Validate, invRepo, custRepo, payRepo, mtClient, custUseCase, notifClient)
	regUseCase := usecase.NewRegistrationUseCase(config.DB, config.Log, config.Validate, regRepo, packageRepo, userRepository, custRepo, radRepo, routerRepo, invRepo, mkClient, notifClient)
	routerUseCase := usecase.NewRouterUseCase(config.DB, config.Log, config.Validate, routerRepo, mkClient)
	dashboardUseCase := usecase.NewDashboardUseCase(config.DB, config.Log, custRepo, invRepo, routerRepo, radRepo)

	//setup controller
	userController := http.NewUserController(config.Log, userUseCase)
	ContactController := http.NewContactController(ContactUseCase, config.Log)
	addressController := http.NewAddressController(addressUseCase, config.Log)

	packageController := http.NewPackageController(config.Log, pkgUseCase)
	regController := http.NewRegistrationController(config.Log, regUseCase)
	custController := http.NewCustomerController(config.Log, custUseCase)
	routerController := http.NewRouterController(config.Log, routerUseCase)
	invoiceController := http.NewInvoiceController(config.Log, invUseCase)
	dashboardController := http.NewDashboardController(config.Log, dashboardUseCase)

	requestLogMiddleware := middleware.NewRequestLogger(userUseCase)
	authMiddlewareAdmin := middleware.NewAuthAdmin(userUseCase, tokenUtil)
	authMiddlewareCustomer := middleware.NewAuthCustomer(userUseCase, tokenUtil)
	authMiddlewareSuperAdmin := middleware.NewAuthSuperAdmin(userUseCase, tokenUtil)
	authMiddlewareDriver := middleware.NewAuthDriver(userUseCase, tokenUtil)

	routeConfig := route.RouteConfig{
		App:                      config.App,
		UserController:           userController,
		ContactController:        ContactController,
		AddressController:        addressController,
		PackageController:        packageController,
		RegistrationController:   regController,
		CustomerController:       custController,
		RouterController:         routerController,
		InvoiceController:        invoiceController,
		DashboardController:      dashboardController,
		AuthAdminMiddleware:      authMiddlewareAdmin,
		AuthCustomerMiddleware:   authMiddlewareCustomer,
		AuthSuperAdminMiddleware: authMiddlewareSuperAdmin,
		AuthDriverMiddleware:     authMiddlewareDriver,
		RequestLoggerMiddleware:  requestLogMiddleware,
	}
	routeConfig.Setup()
}
