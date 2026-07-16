package http

import (
	"golang-clean-architecture/internal/delivery/http/middleware"
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/internal/usecase"
	"math"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type CustomerController struct {
	Log     *logrus.Logger
	UseCase *usecase.CustomerUseCase
}

func NewCustomerController(log *logrus.Logger, useCase *usecase.CustomerUseCase) *CustomerController {
	return &CustomerController{
		Log:     log,
		UseCase: useCase,
	}
}

func (c *CustomerController) List(ctx *fiber.Ctx) error {
	request := &model.SearchCustomerRequest{
		Search: ctx.Query("search", ""),
		Status: ctx.Query("status", ""),
		Page:   ctx.QueryInt("page", 1),
		Size:   ctx.QueryInt("size", 10),
	}

	response, total, err := c.UseCase.List(ctx.UserContext(), request)
	if err != nil {
		return err
	}

	paging := &model.PageMetadata{
		Page:      request.Page,
		Size:      request.Size,
		TotalItem: total,
		TotalPage: int64(math.Ceil(float64(total) / float64(request.Size))),
	}

	return ctx.JSON(model.WebResponse[[]model.CustomerResponse]{
		Data:   response,
		Paging: paging,
	})
}

func (c *CustomerController) Get(ctx *fiber.Ctx) error {
	id := ctx.Params("customerId")
	response, err := c.UseCase.Get(ctx.UserContext(), id)
	if err != nil {
		return err
	}
	return ctx.JSON(model.WebResponse[model.CustomerResponse]{Data: response})
}

func (c *CustomerController) GetHistory(ctx *fiber.Ctx) error {
	id := ctx.Params("customerId")
	response, err := c.UseCase.GetHistory(ctx.UserContext(), id)
	if err != nil {
		return err
	}
	return ctx.JSON(model.WebResponse[[]model.CustomerHistoryResponse]{Data: response})
}

func (c *CustomerController) Update(ctx *fiber.Ctx) error {
	request := new(model.UpdateCustomerRequest)
	if err := ctx.BodyParser(request); err != nil {
		return fiber.ErrBadRequest
	}
	request.ID = ctx.Params("customerId")

	auth := middleware.GetUser(ctx)
	response, err := c.UseCase.Update(ctx.UserContext(), request, auth.ID)
	if err != nil {
		return err
	}
	return ctx.JSON(model.WebResponse[model.CustomerResponse]{Data: response})
}

func (c *CustomerController) Suspend(ctx *fiber.Ctx) error {
	id := ctx.Params("customerId")
	var body struct {
		Notes string `json:"notes"`
	}
	_ = ctx.BodyParser(&body)

	auth := middleware.GetUser(ctx)
	response, err := c.UseCase.Suspend(ctx.UserContext(), id, body.Notes, auth.ID)
	if err != nil {
		return err
	}
	return ctx.JSON(model.WebResponse[model.CustomerResponse]{Data: response})
}

func (c *CustomerController) Unsuspend(ctx *fiber.Ctx) error {
	id := ctx.Params("customerId")
	var body struct {
		Notes string `json:"notes"`
	}
	_ = ctx.BodyParser(&body)

	auth := middleware.GetUser(ctx)
	response, err := c.UseCase.Unsuspend(ctx.UserContext(), id, body.Notes, auth.ID)
	if err != nil {
		return err
	}
	return ctx.JSON(model.WebResponse[model.CustomerResponse]{Data: response})
}

func (c *CustomerController) Terminate(ctx *fiber.Ctx) error {
	id := ctx.Params("customerId")
	var body struct {
		Notes string `json:"notes"`
	}
	_ = ctx.BodyParser(&body)

	auth := middleware.GetUser(ctx)
	response, err := c.UseCase.Terminate(ctx.UserContext(), id, body.Notes, auth.ID)
	if err != nil {
		return err
	}
	return ctx.JSON(model.WebResponse[model.CustomerResponse]{Data: response})
}

func (c *CustomerController) GetCurrentCustomer(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	response, err := c.UseCase.FindByUserId(ctx.UserContext(), auth.ID)
	if err != nil {
		return err
	}
	return ctx.JSON(model.WebResponse[model.CustomerResponse]{Data: response})
}

func (c *CustomerController) RecreatePPPoE(ctx *fiber.Ctx) error {
	id := ctx.Params("customerId")
	auth := middleware.GetUser(ctx)
	response, err := c.UseCase.RecreatePPPoE(ctx.UserContext(), id, auth.ID)
	if err != nil {
		return err
	}
	return ctx.JSON(model.WebResponse[model.CustomerResponse]{Data: response})
}
