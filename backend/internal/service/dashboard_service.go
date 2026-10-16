package service

import (
	"context"
	"time"

	"satset-app/internal/config"
	"satset-app/internal/repository"
)

// DashboardService provides dashboard data for Owner and Kasir.
type DashboardService struct {
	orderRepo    *repository.OrderRepository
	paymentRepo  *repository.PaymentRepository
	expenseRepo  *repository.ExpenseRepository
	customerRepo *repository.CustomerRepository
	serviceRepo  *repository.ServiceRepository
	cfg          *config.Config
}

// NewDashboardService creates a new DashboardService.
func NewDashboardService(
	orderRepo *repository.OrderRepository,
	paymentRepo *repository.PaymentRepository,
	expenseRepo *repository.ExpenseRepository,
	customerRepo *repository.CustomerRepository,
	serviceRepo *repository.ServiceRepository,
	cfg *config.Config,
) *DashboardService {
	return &DashboardService{
		orderRepo:    orderRepo,
		paymentRepo:  paymentRepo,
		expenseRepo:  expenseRepo,
		customerRepo: customerRepo,
		serviceRepo:  serviceRepo,
		cfg:          cfg,
	}
}

// OwnerDashboardData contains the full dashboard data for Owner.
type OwnerDashboardData struct {
	Period          string                        `json:"period"`
	Revenue         *repository.RevenueSummary    `json:"revenue"`
	CashReceived    int64                         `json:"cash_received"`
	TotalExpenses   int64                         `json:"total_expenses"`
	EstimatedProfit int64                         `json:"estimated_profit"`
	ActiveOrders    int64                         `json:"active_orders"`
	OverdueOrders   int64                         `json:"overdue_orders"`
	ReadyForPickup  int64                         `json:"ready_for_pickup"`
	UnpaidOrders    int64                         `json:"unpaid_orders"`
	TopServices     []repository.TopServiceResult `json:"top_services"`
	DailyRevenue    []repository.DailyRevenue     `json:"daily_revenue"`
	OrdersByStatus  []repository.StatusCount      `json:"orders_by_status"`
}

// GetOwnerDashboard returns full dashboard data for Owner.
func (s *DashboardService) GetOwnerDashboard(ctx context.Context, startDate, endDate string) (*OwnerDashboardData, error) {
	// If no dates provided, default to current month
	if startDate == "" || endDate == "" {
		now := time.Now().In(s.cfg.App.Timezone)
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, s.cfg.App.Timezone).Format("2006-01-02")
		endDate = now.Format("2006-01-02")
	}

	revenue, err := s.orderRepo.GetRevenueSummary(ctx, startDate, endDate)
	if err != nil {
		return nil, err
	}

	paymentSummary, err := s.paymentRepo.GetPaymentSummary(ctx, startDate, endDate)
	if err != nil {
		return nil, err
	}

	totalExpenses, err := s.expenseRepo.GetExpenseSummary(ctx, startDate, endDate)
	if err != nil {
		return nil, err
	}

	activeOrders, err := s.orderRepo.CountActive(ctx)
	if err != nil {
		return nil, err
	}

	overdueOrders, err := s.orderRepo.CountOverdue(ctx)
	if err != nil {
		return nil, err
	}

	readyForPickup, err := s.orderRepo.CountReadyForPickup(ctx)
	if err != nil {
		return nil, err
	}

	unpaidOrders, err := s.orderRepo.CountUnpaid(ctx)
	if err != nil {
		return nil, err
	}

	topServices, err := s.serviceRepo.FindTopServices(ctx, startDate, endDate, 5)
	if err != nil {
		return nil, err
	}

	dailyRevenue, err := s.orderRepo.GetDailyRevenue(ctx, startDate, endDate)
	if err != nil {
		return nil, err
	}

	ordersByStatus, err := s.orderRepo.CountByStatus(ctx, startDate, endDate)
	if err != nil {
		return nil, err
	}

	return &OwnerDashboardData{
		Period:          startDate + " s/d " + endDate,
		Revenue:         revenue,
		CashReceived:    paymentSummary.TotalReceived,
		TotalExpenses:   totalExpenses,
		EstimatedProfit: paymentSummary.TotalReceived - totalExpenses,
		ActiveOrders:    activeOrders,
		OverdueOrders:   overdueOrders,
		ReadyForPickup:  readyForPickup,
		UnpaidOrders:    unpaidOrders,
		TopServices:     topServices,
		DailyRevenue:    dailyRevenue,
		OrdersByStatus:  ordersByStatus,
	}, nil
}

// KasirDashboardData contains the dashboard data for Kasir.
type KasirDashboardData struct {
	TodayOrders    int64 `json:"today_orders"`
	ActiveOrders   int64 `json:"active_orders"`
	ReadyForPickup int64 `json:"ready_for_pickup"`
	UnpaidOrders   int64 `json:"unpaid_orders"`
}

// GetKasirDashboard returns dashboard data for Kasir.
func (s *DashboardService) GetKasirDashboard(ctx context.Context) (*KasirDashboardData, error) {
	today := time.Now().In(s.cfg.App.Timezone).Format("2006-01-02")

	todayOrders, err := s.orderRepo.GetTodayOrderCount(ctx, today)
	if err != nil {
		return nil, err
	}

	activeOrders, err := s.orderRepo.CountActive(ctx)
	if err != nil {
		return nil, err
	}

	readyForPickup, err := s.orderRepo.CountReadyForPickup(ctx)
	if err != nil {
		return nil, err
	}

	unpaidOrders, err := s.orderRepo.CountUnpaid(ctx)
	if err != nil {
		return nil, err
	}

	return &KasirDashboardData{
		TodayOrders:    todayOrders,
		ActiveOrders:   activeOrders,
		ReadyForPickup: readyForPickup,
		UnpaidOrders:   unpaidOrders,
	}, nil
}
