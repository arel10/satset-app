package service

import (
	"context"
	"fmt"
	"time"

	"satset-app/internal/config"
	"satset-app/internal/repository"
)

// ReportService provides business reports for Owner.
type ReportService struct {
	orderRepo    *repository.OrderRepository
	paymentRepo  *repository.PaymentRepository
	expenseRepo  *repository.ExpenseRepository
	serviceRepo  *repository.ServiceRepository
	customerRepo *repository.CustomerRepository
	cfg          *config.Config
}

// NewReportService creates a new ReportService.
func NewReportService(
	orderRepo *repository.OrderRepository,
	paymentRepo *repository.PaymentRepository,
	expenseRepo *repository.ExpenseRepository,
	serviceRepo *repository.ServiceRepository,
	customerRepo *repository.CustomerRepository,
	cfg *config.Config,
) *ReportService {
	return &ReportService{
		orderRepo:    orderRepo,
		paymentRepo:  paymentRepo,
		expenseRepo:  expenseRepo,
		serviceRepo:  serviceRepo,
		customerRepo: customerRepo,
		cfg:          cfg,
	}
}

// BusinessReport contains the full business report for a period.
type BusinessReport struct {
	Period             string                               `json:"period"`
	StartDate          string                               `json:"start_date"`
	EndDate            string                               `json:"end_date"`
	TransactionSummary *repository.RevenueSummary           `json:"transaction_summary"`
	DailyRevenue       []repository.DailyRevenue            `json:"daily_revenue"`
	PaymentSummary     *repository.PaymentSummary           `json:"payment_summary"`
	DailyPayments      []repository.DailyPayment            `json:"daily_payments"`
	TotalExpenses      int64                                `json:"total_expenses"`
	DailyExpenses      []repository.DailyExpense            `json:"daily_expenses"`
	ExpenseByCategory  []repository.ExpenseByCategoryResult `json:"expense_by_category"`
	OrdersByStatus     []repository.StatusCount             `json:"orders_by_status"`
	TopServices        []repository.TopServiceResult        `json:"top_services"`
	UnpaidCount        int64                                `json:"unpaid_count"`
	OverdueCount       int64                                `json:"overdue_count"`
	EstimatedProfit    int64                                `json:"estimated_profit"`
	TopCustomers       []repository.CustomerOrderCount      `json:"top_customers"`
}

// GetBusinessReport generates a full business report for a given date range.
func (s *ReportService) GetBusinessReport(ctx context.Context, startDate, endDate string) (*BusinessReport, error) {
	// Validate dates
	if startDate == "" || endDate == "" {
		now := time.Now().In(s.cfg.App.Timezone)
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, s.cfg.App.Timezone).Format("2006-01-02")
		endDate = now.Format("2006-01-02")
	}

	_, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, fmt.Errorf("format start_date tidak valid")
	}
	_, err = time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, fmt.Errorf("format end_date tidak valid")
	}

	report := &BusinessReport{
		Period:    startDate + " s/d " + endDate,
		StartDate: startDate,
		EndDate:   endDate,
	}

	// Transaction summary (revenue from orders)
	report.TransactionSummary, err = s.orderRepo.GetRevenueSummary(ctx, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil ringkasan transaksi: %w", err)
	}

	// Daily revenue
	report.DailyRevenue, err = s.orderRepo.GetDailyRevenue(ctx, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil omzet harian: %w", err)
	}

	// Payment summary (cash received)
	report.PaymentSummary, err = s.paymentRepo.GetPaymentSummary(ctx, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil ringkasan pembayaran: %w", err)
	}

	// Daily payments
	report.DailyPayments, err = s.paymentRepo.GetDailyPayments(ctx, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil pembayaran harian: %w", err)
	}

	// Total expenses
	report.TotalExpenses, err = s.expenseRepo.GetExpenseSummary(ctx, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil total pengeluaran: %w", err)
	}

	// Daily expenses
	report.DailyExpenses, err = s.expenseRepo.GetDailyExpenses(ctx, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil pengeluaran harian: %w", err)
	}

	// Expense by category
	report.ExpenseByCategory, err = s.expenseRepo.GetExpenseByCategory(ctx, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil pengeluaran per kategori: %w", err)
	}

	// Orders by status
	report.OrdersByStatus, err = s.orderRepo.CountByStatus(ctx, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil pesanan per status: %w", err)
	}

	// Top services
	report.TopServices, err = s.serviceRepo.FindTopServices(ctx, startDate, endDate, 10)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil layanan terlaris: %w", err)
	}

	// Unpaid and overdue counts
	report.UnpaidCount, err = s.orderRepo.CountUnpaid(ctx)
	if err != nil {
		return nil, fmt.Errorf("gagal menghitung pesanan belum lunas: %w", err)
	}

	report.OverdueCount, err = s.orderRepo.CountOverdue(ctx)
	if err != nil {
		return nil, fmt.Errorf("gagal menghitung pesanan terlambat: %w", err)
	}

	// Top customers
	report.TopCustomers, err = s.customerRepo.FindTopCustomers(ctx, 10)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil pelanggan teratas: %w", err)
	}

	// Estimated profit
	report.EstimatedProfit = report.PaymentSummary.TotalReceived - report.TotalExpenses

	return report, nil
}
