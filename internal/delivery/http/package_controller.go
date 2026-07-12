package http

import (
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type PackageController struct {
	Log     *logrus.Logger
	UseCase *usecase.PackageUseCase
}

func NewPackageController(log *logrus.Logger, useCase *usecase.PackageUseCase) *PackageController {
	return &PackageController{
		Log:     log,
		UseCase: useCase,
	}
}

func (c *PackageController) Create(ctx *fiber.Ctx) error {
	request := new(model.CreatePackageRequest)
	if err := ctx.BodyParser(request); err != nil {
		return fiber.ErrBadRequest
	}

	response, err := c.UseCase.Create(ctx.UserContext(), request)
	if err != nil {
		return err
	}

	return ctx.JSON(model.WebResponse[model.PackageResponse]{Data: response})
}

func (c *PackageController) Update(ctx *fiber.Ctx) error {
	request := new(model.UpdatePackageRequest)
	if err := ctx.BodyParser(request); err != nil {
		return fiber.ErrBadRequest
	}
	request.ID = ctx.Params("packageId")

	response, err := c.UseCase.Update(ctx.UserContext(), request)
	if err != nil {
		return err
	}

	return ctx.JSON(model.WebResponse[model.PackageResponse]{Data: response})
}

func (c *PackageController) List(ctx *fiber.Ctx) error {
	response, err := c.UseCase.List(ctx.UserContext())
	if err != nil {
		return err
	}

	return ctx.JSON(model.WebResponse[[]model.PackageResponse]{Data: response})
}

func (c *PackageController) Get(ctx *fiber.Ctx) error {
	id := ctx.Params("packageId")
	response, err := c.UseCase.Get(ctx.UserContext(), id)
	if err != nil {
		return err
	}

	return ctx.JSON(model.WebResponse[model.PackageResponse]{Data: response})
}

func (c *PackageController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("packageId")
	err := c.UseCase.Delete(ctx.UserContext(), id)
	if err != nil {
		return err
	}

	return ctx.JSON(model.WebResponse[bool]{Data: true})
}
