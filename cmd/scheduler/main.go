package main

import (
	"context"
	"fmt"
	"time"

	"golang-clean-architecture/internal/config"
	"golang-clean-architecture/internal/entity"
	"golang-clean-architecture/internal/gateway/mikrotik"
	"golang-clean-architecture/internal/gateway/notification"
	"golang-clean-architecture/internal/repository"
	"golang-clean-architecture/internal/usecase"

	"github.com/sirupsen/logrus"
	"github.com/subosito/gotenv"
	"gorm.io/gorm"
)

func main() {
	gotenv.Load()

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
	notifClient := notification.NewNotificationClient(log)

	custUseCase := usecase.NewCustomerUseCase(db, log, validate, custRepo, packageRepo, routerRepo, radRepo, histRepo, mkClient)
	invUseCase := usecase.NewInvoiceUseCase(db, log, validate, invRepo, custRepo, payRepo, nil, custUseCase, notifClient)

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
	if err := invUseCase.AutoGenerateInvoices(ctx); err != nil {
		log.Errorf("Scheduler failed to auto generate invoices: %v", err)
	}

	// 2. Suspend Owed Customers
	suspendOwedCustomers(ctx, db, log, custUseCase)
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
