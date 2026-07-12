package usecase

import (
	"context"
	"golang-clean-architecture/internal/entity"
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/internal/model/converter"
	"golang-clean-architecture/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PackageUseCase struct {
	DB                *gorm.DB
	Log               *logrus.Logger
	Validate          *validator.Validate
	PackageRepository *repository.PackageRepository
}

func NewPackageUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, repo *repository.PackageRepository) *PackageUseCase {
	return &PackageUseCase{
		DB:                db,
		Log:               log,
		Validate:          validate,
		PackageRepository: repo,
	}
}

func (c *PackageUseCase) Create(ctx context.Context, request *model.CreatePackageRequest) (model.PackageResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body for package create: %+v", err)
		return model.PackageResponse{}, fiber.ErrBadRequest
	}

	pkg := &entity.InternetPackage{
		ID:              uuid.NewString(),
		Name:            request.Name,
		SpeedMbps:       request.SpeedMbps,
		Price:           request.Price,
		InstallationFee: request.InstallationFee,
		TaxRate:         request.TaxRate,
		IsActive:        true,
	}

	if err := c.PackageRepository.Create(tx, pkg); err != nil {
		c.Log.Errorf("Failed to create package: %+v", err)
		return model.PackageResponse{}, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		return model.PackageResponse{}, fiber.ErrInternalServerError
	}

	return converter.PackageToResponse(pkg), nil
}

func (c *PackageUseCase) Update(ctx context.Context, request *model.UpdatePackageRequest) (model.PackageResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		return model.PackageResponse{}, fiber.ErrBadRequest
	}

	pkg := new(entity.InternetPackage)
	if err := c.PackageRepository.FindById(tx, pkg, request.ID); err != nil {
		return model.PackageResponse{}, fiber.ErrNotFound
	}

	if request.Name != "" {
		pkg.Name = request.Name
	}
	if request.SpeedMbps != 0 {
		pkg.SpeedMbps = request.SpeedMbps
	}
	if request.Price != 0 {
		pkg.Price = request.Price
	}
	if request.InstallationFee != 0 {
		pkg.InstallationFee = request.InstallationFee
	}
	if request.TaxRate != 0 {
		pkg.TaxRate = request.TaxRate
	}
	if request.IsActive != nil {
		pkg.IsActive = *request.IsActive
	}

	if err := c.PackageRepository.Update(tx, pkg); err != nil {
		return model.PackageResponse{}, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		return model.PackageResponse{}, fiber.ErrInternalServerError
	}

	return converter.PackageToResponse(pkg), nil
}

func (c *PackageUseCase) List(ctx context.Context) ([]model.PackageResponse, error) {
	tx := c.DB.WithContext(ctx)
	packages, err := c.PackageRepository.FindAll(tx)
	if err != nil {
		return nil, fiber.ErrInternalServerError
	}

	var responses []model.PackageResponse
	for _, p := range packages {
		responses = append(responses, converter.PackageToResponse(&p))
	}
	return responses, nil
}

func (c *PackageUseCase) Get(ctx context.Context, id string) (model.PackageResponse, error) {
	tx := c.DB.WithContext(ctx)
	pkg := new(entity.InternetPackage)
	if err := c.PackageRepository.FindById(tx, pkg, id); err != nil {
		return model.PackageResponse{}, fiber.ErrNotFound
	}
	return converter.PackageToResponse(pkg), nil
}

func (c *PackageUseCase) Delete(ctx context.Context, id string) error {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	pkg := new(entity.InternetPackage)
	if err := c.PackageRepository.FindById(tx, pkg, id); err != nil {
		return fiber.ErrNotFound
	}

	if err := c.PackageRepository.Delete(tx, pkg); err != nil {
		return fiber.ErrInternalServerError
	}

	return tx.Commit().Error
}
