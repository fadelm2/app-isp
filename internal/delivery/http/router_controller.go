package http

import (
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type RouterController struct {
	Log     *logrus.Logger
	UseCase *usecase.RouterUseCase
}

func NewRouterController(log *logrus.Logger, useCase *usecase.RouterUseCase) *RouterController {
	return &RouterController{
		Log:     log,
		UseCase: useCase,
	}
}

func (c *RouterController) Create(ctx *fiber.Ctx) error {
	request := new(model.CreateRouterRequest)
	if err := ctx.BodyParser(request); err != nil {
		return fiber.ErrBadRequest
	}

	response, err := c.UseCase.Create(ctx.UserContext(), request)
	if err != nil {
		return err
	}

	return ctx.JSON(model.WebResponse[*model.RouterResponse]{Data: response})
}

func (c *RouterController) List(ctx *fiber.Ctx) error {
	response, err := c.UseCase.List(ctx.UserContext())
	if err != nil {
		return err
	}

	return ctx.JSON(model.WebResponse[[]model.RouterResponse]{Data: response})
}

func (c *RouterController) Get(ctx *fiber.Ctx) error {
	id := ctx.Params("routerId")
	response, err := c.UseCase.Get(ctx.UserContext(), id)
	if err != nil {
		return err
	}

	return ctx.JSON(model.WebResponse[*model.RouterResponse]{Data: response})
}

func (c *RouterController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("routerId")
	err := c.UseCase.Delete(ctx.UserContext(), id)
	if err != nil {
		return err
	}

	return ctx.JSON(model.WebResponse[bool]{Data: true})
}
