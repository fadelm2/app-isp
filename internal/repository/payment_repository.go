package repository

import (
	"github.com/sirupsen/logrus"
	"golang-clean-architecture/internal/entity"
	"gorm.io/gorm"
)

type PaymentRepository struct {
	Repository[entity.Payment]
	Log *logrus.Logger
}

func NewPaymentRepository(log *logrus.Logger) *PaymentRepository {
	return &PaymentRepository{
		Log: log,
	}
}

func (r *PaymentRepository) FindAll(db *gorm.DB) ([]entity.Payment, error) {
	var payments []entity.Payment
	if err := db.Preload("Invoice").Find(&payments).Error; err != nil {
		return nil, err
	}
	return payments, nil
}
