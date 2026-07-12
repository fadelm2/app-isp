package repository

import (
	"github.com/sirupsen/logrus"
	"golang-clean-architecture/internal/entity"
	"gorm.io/gorm"
)

type InvoiceRepository struct {
	Repository[entity.Invoice]
	Log *logrus.Logger
}

func NewInvoiceRepository(log *logrus.Logger) *InvoiceRepository {
	return &InvoiceRepository{
		Log: log,
	}
}

func (r *InvoiceRepository) FindAll(db *gorm.DB) ([]entity.Invoice, error) {
	var invoices []entity.Invoice
	if err := db.Preload("Customer").Preload("Customer.User").Preload("Customer.Package").Find(&invoices).Error; err != nil {
		return nil, err
	}
	return invoices, nil
}

func (r *InvoiceRepository) FindByCustomerId(db *gorm.DB, customerId string) ([]entity.Invoice, error) {
	var invoices []entity.Invoice
	if err := db.Preload("Customer").Preload("Customer.User").Preload("Customer.Package").Where("customer_id = ?", customerId).Find(&invoices).Error; err != nil {
		return nil, err
	}
	return invoices, nil
}

func (r *InvoiceRepository) FindByIdWithDetails(db *gorm.DB, invoice *entity.Invoice, id string) error {
	return db.Preload("Customer").Preload("Customer.User").Preload("Customer.Package").Where("id = ?", id).First(invoice).Error
}
