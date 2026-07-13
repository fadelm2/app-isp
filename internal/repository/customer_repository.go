package repository

import (
	"github.com/sirupsen/logrus"
	"golang-clean-architecture/internal/entity"
	"golang-clean-architecture/internal/model"
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
	if err := db.Preload("User").Preload("Package").Preload("Router").Preload("Registration").Find(&customers).Error; err != nil {
		return nil, err
	}
	return customers, nil
}

func (r *CustomerRepository) FindByIdWithDetails(db *gorm.DB, customer *entity.Customer, id string) error {
	return db.Preload("User").Preload("Package").Preload("Router").Preload("Registration").Where("id = ?", id).First(customer).Error
}

func (r *CustomerRepository) FindByUserId(db *gorm.DB, customer *entity.Customer, userId string) error {
	return db.Preload("User").Preload("Package").Preload("Router").Preload("Registration").Where("user_id = ?", userId).First(customer).Error
}

func (r *CustomerRepository) Search(db *gorm.DB, request *model.SearchCustomerRequest) ([]entity.Customer, int64, error) {
	var customers []entity.Customer

	query := db.Preload("User").Preload("Package").Preload("Router").Preload("Registration").
		Scopes(r.FilterCustomer(request)).
		Order("created_at DESC").
		Offset((request.Page - 1) * request.Size).
		Limit(request.Size)

	if err := query.Find(&customers).Error; err != nil {
		return nil, 0, err
	}

	var total int64
	if err := db.Model(&entity.Customer{}).Scopes(r.FilterCustomer(request)).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return customers, total, nil
}

func (r *CustomerRepository) FilterCustomer(request *model.SearchCustomerRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		if status := request.Status; status != "" {
			tx = tx.Where("customers.status = ?", status)
		}

		if search := request.Search; search != "" {
			keyword := "%" + search + "%"
			tx = tx.Joins("LEFT JOIN users ON users.id = customers.user_id").
				Where("customers.id LIKE ? OR customers.ppp_username LIKE ? OR customers.odp_number LIKE ? OR users.username LIKE ? OR users.email LIKE ?",
					keyword, keyword, keyword, keyword, keyword)
		}

		return tx
	}
}
