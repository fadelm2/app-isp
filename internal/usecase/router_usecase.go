package usecase

import (
	"context"
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

type RouterUseCase struct {
	DB               *gorm.DB
	Log              *logrus.Logger
	Validate         *validator.Validate
	RouterRepository *repository.RouterRepository
	MikrotikClient   *mikrotik.MikrotikClient
}

func NewRouterUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, repo *repository.RouterRepository, mkClient *mikrotik.MikrotikClient) *RouterUseCase {
	return &RouterUseCase{
		DB:               db,
		Log:              log,
		Validate:         validate,
		RouterRepository: repo,
		MikrotikClient:   mkClient,
	}
}

func (c *RouterUseCase) Create(ctx context.Context, request *model.CreateRouterRequest) (*model.RouterResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		return nil, fiber.ErrBadRequest
	}

	router := &entity.Router{
		ID:       uuid.NewString(),
		Name:     request.Name,
		Host:     request.Host,
		Port:     request.Port,
		Username: request.Username,
		Password: request.Password,
		Status:   "offline",
	}

	// Test connection
	status := "offline"
	ok, err := c.MikrotikClient.PingRouter(router.Host, router.Port, router.Username, router.Password)
	if err == nil && ok {
		status = "online"
	}
	router.Status = status

	if err := c.RouterRepository.Create(tx, router); err != nil {
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fiber.ErrInternalServerError
	}

	return converter.RouterToResponse(router), nil
}

func (c *RouterUseCase) List(ctx context.Context) ([]model.RouterResponse, error) {
	tx := c.DB.WithContext(ctx)
	routers, err := c.RouterRepository.FindAll(tx)
	if err != nil {
		return nil, fiber.ErrInternalServerError
	}

	var responses []model.RouterResponse
	for _, r := range routers {
		responses = append(responses, *converter.RouterToResponse(&r))
	}
	return responses, nil
}

func (c *RouterUseCase) Get(ctx context.Context, id string) (*model.RouterResponse, error) {
	tx := c.DB.WithContext(ctx)
	router := new(entity.Router)
	if err := c.RouterRepository.FindById(tx, router, id); err != nil {
		return nil, fiber.ErrNotFound
	}
	return converter.RouterToResponse(router), nil
}

func (c *RouterUseCase) Delete(ctx context.Context, id string) error {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	router := new(entity.Router)
	if err := c.RouterRepository.FindById(tx, router, id); err != nil {
		return fiber.ErrNotFound
	}

	if err := c.RouterRepository.Delete(tx, router); err != nil {
		return fiber.ErrInternalServerError
	}

	return tx.Commit().Error
}
