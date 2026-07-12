package usecase

import (
	"context"
	"time"

	"golang-clean-architecture/internal/entity"
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/internal/repository"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type DashboardUseCase struct {
	DB                 *gorm.DB
	Log                *logrus.Logger
	CustomerRepository *repository.CustomerRepository
	InvoiceRepository  *repository.InvoiceRepository
	RouterRepository   *repository.RouterRepository
	RadiusRepository   *repository.RadiusRepository
}

func NewDashboardUseCase(
	db *gorm.DB,
	log *logrus.Logger,
	custRepo *repository.CustomerRepository,
	invRepo *repository.InvoiceRepository,
	routerRepo *repository.RouterRepository,
	radRepo *repository.RadiusRepository,
) *DashboardUseCase {
	return &DashboardUseCase{
		DB:                 db,
		Log:                log,
		CustomerRepository: custRepo,
		InvoiceRepository:  invRepo,
		RouterRepository:   routerRepo,
		RadiusRepository:   radRepo,
	}
}

func (c *DashboardUseCase) GetStats(ctx context.Context) (*model.DashboardStatsResponse, error) {
	tx := c.DB.WithContext(ctx)

	var totalCust, activeCust, suspendedCust int64
	tx.Model(&entity.Customer{}).Count(&totalCust)
	tx.Model(&entity.Customer{}).Where("status = ?", "active").Count(&activeCust)
	tx.Model(&entity.Customer{}).Where("status = ?", "suspended").Count(&suspendedCust)

	var owedCust int64
	tx.Model(&entity.Invoice{}).Where("status = ?", "owed").Distinct("customer_id").Count(&owedCust)

	// Today's Payments
	var todayPayments float64
	todayStart := time.Now().Truncate(24 * time.Hour).UnixMilli()
	tx.Model(&entity.Payment{}).Where("paid_at >= ? AND status = ?", todayStart, "settlement").Select("COALESCE(SUM(paid_amount), 0)").Row().Scan(&todayPayments)

	// Monthly Revenue
	var monthlyRevenue float64
	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local).UnixMilli()
	tx.Model(&entity.Invoice{}).Where("paid_at >= ? AND status = ?", monthStart, "paid").Select("COALESCE(SUM(total_amount), 0)").Row().Scan(&monthlyRevenue)

	// Router Status
	routerStatus := "offline"
	var routers []entity.Router
	tx.Find(&routers)
	for _, r := range routers {
		if r.Status == "online" {
			routerStatus = "online"
			break
		}
	}

	// RADIUS sessions
	var onlineUsers int
	activeSessions, err := c.RadiusRepository.FindActiveSessions(tx)
	if err == nil {
		onlineUsers = len(activeSessions)
	}

	offlineUsers := int(activeCust) - onlineUsers
	if offlineUsers < 0 {
		offlineUsers = 0
	}

	return &model.DashboardStatsResponse{
		TotalCustomers:     totalCust,
		ActiveCustomers:    activeCust,
		SuspendedCustomers: suspendedCust,
		OwedCustomers:      owedCust,
		TodayPayments:      todayPayments,
		MonthlyRevenue:     monthlyRevenue,
		RouterStatus:       routerStatus,
		OnlineUsers:        onlineUsers,
		OfflineUsers:       offlineUsers,
	}, nil
}
