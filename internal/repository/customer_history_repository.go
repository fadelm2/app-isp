package repository

import (
	"github.com/sirupsen/logrus"
	"golang-clean-architecture/internal/entity"
	"gorm.io/gorm"
)

type CustomerHistoryRepository struct {
	Repository[entity.CustomerHistory]
	Log *logrus.Logger
}

func NewCustomerHistoryRepository(log *logrus.Logger) *CustomerHistoryRepository {
	return &CustomerHistoryRepository{
		Log: log,
	}
}

func (r *CustomerHistoryRepository) FindAllByCustomerId(db *gorm.DB, customerId string) ([]entity.CustomerHistory, error) {
	var histories []entity.CustomerHistory
	if err := db.Preload("User").Where("customer_id = ?", customerId).Order("created_at desc").Find(&histories).Error; err != nil {
		return nil, err
	}
	return histories, nil
}
