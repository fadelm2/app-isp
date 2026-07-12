package repository

import (
	"github.com/sirupsen/logrus"
	"golang-clean-architecture/internal/entity"
	"gorm.io/gorm"
)

type RegistrationRepository struct {
	Repository[entity.Registration]
	Log *logrus.Logger
}

func NewRegistrationRepository(log *logrus.Logger) *RegistrationRepository {
	return &RegistrationRepository{
		Log: log,
	}
}

func (r *RegistrationRepository) FindAll(db *gorm.DB) ([]entity.Registration, error) {
	var registrations []entity.Registration
	if err := db.Preload("Package").Find(&registrations).Error; err != nil {
		return nil, err
	}
	return registrations, nil
}
