package route

import (
	"golang-clean-architecture/internal/delivery/http"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App                      *fiber.App
	UserController           *http.UserController
	ContactController        *http.ContactController
	AddressController        *http.AddressController
	PackageController        *http.PackageController
	RegistrationController   *http.RegistrationController
	CustomerController       *http.CustomerController
	RouterController         *http.RouterController
	InvoiceController        *http.InvoiceController
	DashboardController      *http.DashboardController
	AuthCustomerMiddleware   fiber.Handler
	AuthAdminMiddleware      fiber.Handler
	AuthSuperAdminMiddleware fiber.Handler
	AuthDriverMiddleware     fiber.Handler
	AuthMiddleWare           fiber.Handler
	RequestLoggerMiddleware  fiber.Handler
}

func (c *RouteConfig) Setup() {
	c.SetupGuestRoute()
	c.SetupAuthRoute()
	c.SetupAuthAdminRoute()
	c.SetupAuthSuperAdminRoute()
	c.SetupAuthDriverRoute()
}

func (c *RouteConfig) SetupGuestRoute() {
	c.App.Post("/api/users/_Login", c.UserController.Login)
	c.App.Post("/api/registrations", c.RegistrationController.Create)
	c.App.Post("/api/webhook/midtrans", c.InvoiceController.MidtransWebhook)
	c.App.Get("/api/packages", c.PackageController.List)
}

func (c *RouteConfig) SetupAuthAdminRoute() {
	admin := c.App.Group("/api/admin", c.AuthAdminMiddleware)
	admin.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"message": "Admin API is accessible",
		})
	})

	// Dashboard
	admin.Get("/dashboard", c.DashboardController.GetStats)

	// Packages
	admin.Get("/packages", c.PackageController.List)
	admin.Post("/packages", c.PackageController.Create)
	admin.Get("/packages/:packageId", c.PackageController.Get)
	admin.Patch("/packages/:packageId", c.PackageController.Update)
	admin.Delete("/packages/:packageId", c.PackageController.Delete)

	// Registrations
	admin.Get("/registrations", c.RegistrationController.List)
	admin.Patch("/registrations/:registrationId/status", c.RegistrationController.UpdateStatus)

	// Customers
	admin.Get("/customers", c.CustomerController.List)
	admin.Get("/customers/:customerId", c.CustomerController.Get)
	admin.Patch("/customers/:customerId", c.CustomerController.Update)
	admin.Post("/customers/:customerId/_suspend", c.CustomerController.Suspend)
	admin.Post("/customers/:customerId/_unsuspend", c.CustomerController.Unsuspend)
	admin.Post("/customers/:customerId/_terminate", c.CustomerController.Terminate)
	admin.Get("/customers/:customerId/history", c.CustomerController.GetHistory)

	// Routers
	admin.Get("/routers", c.RouterController.List)
	admin.Post("/routers", c.RouterController.Create)
	admin.Get("/routers/:routerId", c.RouterController.Get)
	admin.Delete("/routers/:routerId", c.RouterController.Delete)

	// Invoices
	admin.Get("/invoices", c.InvoiceController.List)
	admin.Post("/invoices", c.InvoiceController.Create)
	admin.Get("/invoices/:invoiceId", c.InvoiceController.Get)
}

func (c *RouteConfig) SetupAuthDriverRoute() {
	driver := c.App.Group("/api/driver", c.AuthDriverMiddleware)
	driver.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"message": "Driver API is accessible",
		})
	})
}

func (c *RouteConfig) SetupAuthSuperAdminRoute() {
	superAdmin := c.App.Group("/api/superadmin", c.AuthSuperAdminMiddleware)
	superAdmin.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"message": "Superadmin API is accessible",
		})
	})
}

func (c *RouteConfig) SetupAuthRoute() {
	// Base user/customer routes
	customer := c.App.Group("/api/customer", c.AuthCustomerMiddleware)
	customer.Get("/me", c.CustomerController.GetCurrentCustomer)
	customer.Get("/invoices/:invoiceId/pay", c.InvoiceController.GetSnapToken)

	admin := c.App.Group("/api/admin", c.AuthAdminMiddleware)

	admin.Delete("/users", c.UserController.Logout)
	admin.Patch("/users/_current", c.UserController.Update)
	admin.Get("/users/_current", c.UserController.Current)

	admin.Get("/contacts", c.ContactController.List)
	admin.Post("/contacts", c.ContactController.Create)
	admin.Put("/contacts/:contactId", c.ContactController.Update)
	admin.Get("/contacts/:contactId", c.ContactController.Get)
	admin.Delete("/contacts/:contactId", c.ContactController.Delete)

	admin.Get("/contacts/:contactId/addresses", c.AddressController.List)
	admin.Post("/contacts/:contactId/addresses", c.AddressController.Create)
	admin.Put("/contacts/:contactId/addresses/:addressId", c.AddressController.Update)
	admin.Get("/contacts/:contactId/addresses/:addressId", c.AddressController.Get)
	admin.Delete("/contacts/:contactId/addresses/:addressId", c.AddressController.Delete)
}
