package http

import (
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type InvoiceController struct {
	Log     *logrus.Logger
	UseCase *usecase.InvoiceUseCase
}

func NewInvoiceController(log *logrus.Logger, useCase *usecase.InvoiceUseCase) *InvoiceController {
	return &InvoiceController{
		Log:     log,
		UseCase: useCase,
	}
}

func (c *InvoiceController) Create(ctx *fiber.Ctx) error {
	request := new(model.CreateInvoiceRequest)
	if err := ctx.BodyParser(request); err != nil {
		return fiber.ErrBadRequest
	}

	response, err := c.UseCase.Create(ctx.UserContext(), request)
	if err != nil {
		return err
	}

	return ctx.JSON(model.WebResponse[model.InvoiceResponse]{Data: response})
}

func (c *InvoiceController) List(ctx *fiber.Ctx) error {
	response, err := c.UseCase.List(ctx.UserContext())
	if err != nil {
		return err
	}

	return ctx.JSON(model.WebResponse[[]model.InvoiceResponse]{Data: response})
}

func (c *InvoiceController) Get(ctx *fiber.Ctx) error {
	id := ctx.Params("invoiceId")
	response, err := c.UseCase.Get(ctx.UserContext(), id)
	if err != nil {
		return err
	}

	return ctx.JSON(model.WebResponse[model.InvoiceResponse]{Data: response})
}

func (c *InvoiceController) GetSnapToken(ctx *fiber.Ctx) error {
	id := ctx.Params("invoiceId")
	token, err := c.UseCase.GetSnapToken(ctx.UserContext(), id)
	if err != nil {
		return err
	}

	return ctx.JSON(model.WebResponse[string]{Data: token})
}

func (c *InvoiceController) MidtransWebhook(ctx *fiber.Ctx) error {
	var payload map[string]interface{}
	if err := ctx.BodyParser(&payload); err != nil {
		return fiber.ErrBadRequest
	}

	err := c.UseCase.ProcessWebhook(ctx.UserContext(), payload)
	if err != nil {
		c.Log.Errorf("Midtrans Webhook process failed: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(fiber.Map{"status": "ok"})
}
