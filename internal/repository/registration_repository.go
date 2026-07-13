package repository

import (
	"github.com/sirupsen/logrus"
	"golang-clean-architecture/internal/entity"
	"golang-clean-architecture/internal/model"
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

func (r *RegistrationRepository) Search(db *gorm.DB, request *model.SearchRegistrationRequest) ([]entity.Registration, int64, error) {
	var registrations []entity.Registration

	query := db.Preload("Package").
		Scopes(r.FilterRegistration(request)).
		Order("created_at DESC").
		Offset((request.Page - 1) * request.Size).
		Limit(request.Size)

	if err := query.Find(&registrations).Error; err != nil {
		return nil, 0, err
	}

	var total int64
	if err := db.Model(&entity.Registration{}).Scopes(r.FilterRegistration(request)).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return registrations, total, nil
}

func (r *RegistrationRepository) FilterRegistration(request *model.SearchRegistrationRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		if status := request.Status; status != "" {
			tx = tx.Where("registrations.status = ?", status)
		}

		if search := request.Search; search != "" {
			keyword := "%" + search + "%"
			tx = tx.Where("registrations.id LIKE ? OR registrations.full_name LIKE ? OR registrations.nik LIKE ? OR registrations.phone LIKE ? OR registrations.email LIKE ? OR registrations.odp_number LIKE ?",
				keyword, keyword, keyword, keyword, keyword, keyword)
		}

		return tx
	}
}
