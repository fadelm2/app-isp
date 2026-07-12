package http

import (
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type DashboardController struct {
	Log     *logrus.Logger
	UseCase *usecase.DashboardUseCase
}

func NewDashboardController(log *logrus.Logger, useCase *usecase.DashboardUseCase) *DashboardController {
	return &DashboardController{
		Log:     log,
		UseCase: useCase,
	}
}

func (c *DashboardController) GetStats(ctx *fiber.Ctx) error {
	response, err := c.UseCase.GetStats(ctx.UserContext())
	if err != nil {
		return err
	}

	return ctx.JSON(model.WebResponse[*model.DashboardStatsResponse]{Data: response})
}
