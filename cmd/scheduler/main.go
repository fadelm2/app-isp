package main

import (
	"context"
	"fmt"
	"time"

	"golang-clean-architecture/internal/config"
	"golang-clean-architecture/internal/entity"
	"golang-clean-architecture/internal/gateway/mikrotik"
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/internal/repository"
	"golang-clean-architecture/internal/usecase"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func main() {
	viperConfig := config.NewViper()
	log := config.NewLogger(viperConfig)
	db := config.NewDatabase(viperConfig, log)
	validate := config.NewValidator(viperConfig)

	log.Info("Starting ISP Management System Scheduler Daemon...")

	packageRepo := repository.NewPackageRepository(log)
	custRepo := repository.NewCustomerRepository(log)
	routerRepo := repository.NewRouterRepository(log)
	radRepo := repository.NewRadiusRepository(log)
	histRepo := repository.NewCustomerHistoryRepository(log)
	invRepo := repository.NewInvoiceRepository(log)
	payRepo := repository.NewPaymentRepository(log)

	mkClient := mikrotik.NewMikrotikClient(log)

	custUseCase := usecase.NewCustomerUseCase(db, log, validate, custRepo, packageRepo, routerRepo, radRepo, histRepo, mkClient)
	invUseCase := usecase.NewInvoiceUseCase(db, log, validate, invRepo, custRepo, payRepo, nil, custUseCase)

	ctx := context.Background()

	// Run initially on startup
	runSchedulerTasks(ctx, db, log, invUseCase, custUseCase)

	// Run every 10 minutes
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		runSchedulerTasks(ctx, db, log, invUseCase, custUseCase)
	}
}

func runSchedulerTasks(ctx context.Context, db *gorm.DB, log *logrus.Logger, invUseCase *usecase.InvoiceUseCase, custUseCase *usecase.CustomerUseCase) {
	log.Info("[Scheduler] Running recurring tasks...")

	// 1. Generate Monthly Invoices
	generateMonthlyInvoices(ctx, db, log, invUseCase)

	// 2. Suspend Owed Customers
	suspendOwedCustomers(ctx, db, log, custUseCase)
}

func generateMonthlyInvoices(ctx context.Context, db *gorm.DB, log *logrus.Logger, invUseCase *usecase.InvoiceUseCase) {
	var customers []entity.Customer
	if err := db.Where("status = ?", "active").Find(&customers).Error; err != nil {
		log.Errorf("Scheduler failed to load customers: %v", err)
		return
	}

	now := time.Now()
	for _, cust := range customers {
		// Attempt to create invoice for the current month and year
		_, err := invUseCase.Create(ctx, &model.CreateInvoiceRequest{
			CustomerID:  cust.ID,
			PeriodMonth: int(now.Month()),
			PeriodYear:  now.Year(),
		})
		if err == nil {
			log.Infof("[Scheduler] Auto-generated monthly invoice for customer: %s", cust.ID)
		}
	}
}

func suspendOwedCustomers(ctx context.Context, db *gorm.DB, log *logrus.Logger, custUseCase *usecase.CustomerUseCase) {
	var unpaidInvoices []entity.Invoice
	now := time.Now().UnixMilli()

	// Find invoices whose due date has passed, and status is pending
	err := db.Where("due_date < ? AND status = ?", now, "pending").Find(&unpaidInvoices).Error
	if err != nil {
		log.Errorf("Scheduler failed to load overdue invoices: %v", err)
		return
	}

	for _, inv := range unpaidInvoices {
		// Update invoice status to owed
		inv.Status = "owed"
		db.Save(&inv)

		log.Warnf("[Scheduler] Invoice %s is overdue. Suspending customer: %s", inv.ID, inv.CustomerID)
		_, err = custUseCase.Suspend(ctx, inv.CustomerID, fmt.Sprintf("Automatic suspension: invoice %s overdue", inv.ID), "SYSTEM")
		if err != nil {
			log.Errorf("Scheduler failed to suspend customer %s: %v", inv.CustomerID, err)
		}
	}
}
