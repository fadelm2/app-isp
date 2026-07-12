package repository

import (
	"github.com/sirupsen/logrus"
	"golang-clean-architecture/internal/entity"
	"gorm.io/gorm"
)

type CustomerRepository struct {
	Repository[entity.Customer]
	Log *logrus.Logger
}

func NewCustomerRepository(log *logrus.Logger) *CustomerRepository {
	return &CustomerRepository{
		Log: log,
	}
}

func (r *CustomerRepository) FindAll(db *gorm.DB) ([]entity.Customer, error) {
	var customers []entity.Customer
	if err := db.Preload("User").Preload("Package").Preload("Router").Find(&customers).Error; err != nil {
		return nil, err
	}
	return customers, nil
}

func (r *CustomerRepository) FindByIdWithDetails(db *gorm.DB, customer *entity.Customer, id string) error {
	return db.Preload("User").Preload("Package").Preload("Router").Where("id = ?", id).First(customer).Error
}

func (r *CustomerRepository) FindByUserId(db *gorm.DB, customer *entity.Customer, userId string) error {
	return db.Preload("User").Preload("Package").Preload("Router").Where("user_id = ?", userId).First(customer).Error
}
