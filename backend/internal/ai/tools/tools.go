package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	"satset-app/internal/model"
	"satset-app/internal/repository"
)

// ToolDefinition defines a function that Gemini can call.
type ToolDefinition struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

// ToolExecutor executes read-only tool calls against the SatSet database.
type ToolExecutor struct {
	db           *gorm.DB
	orderRepo    *repository.OrderRepository
	paymentRepo  *repository.PaymentRepository
	expenseRepo  *repository.ExpenseRepository
	customerRepo *repository.CustomerRepository
	serviceRepo  *repository.ServiceRepository
}

// NewToolExecutor creates a new ToolExecutor instance.
func NewToolExecutor(
	db *gorm.DB,
	orderRepo *repository.OrderRepository,
	paymentRepo *repository.PaymentRepository,
	expenseRepo *repository.ExpenseRepository,
	customerRepo *repository.CustomerRepository,
	serviceRepo *repository.ServiceRepository,
) *ToolExecutor {
	return &ToolExecutor{
		db:           db,
		orderRepo:    orderRepo,
		paymentRepo:  paymentRepo,
		expenseRepo:  expenseRepo,
		customerRepo: customerRepo,
		serviceRepo:  serviceRepo,
	}
}

// GetToolDefinitions returns the 13 read-only tool definitions for Gemini function calling.
func GetToolDefinitions() []ToolDefinition {
	return []ToolDefinition{
		{
			Name:        "get_business_summary",
			Description: "Mendapatkan ringkasan performa bisnis (total omzet, uang masuk/kas, pengeluaran, estimasi laba kotor, dan total pesanan) untuk rentang tanggal tertentu (format YYYY-MM-DD).",
			Parameters: map[string]interface{}{
				"type": "OBJECT",
				"properties": map[string]interface{}{
					"start_date": map[string]interface{}{"type": "STRING", "description": "Tanggal awal YYYY-MM-DD"},
					"end_date":   map[string]interface{}{"type": "STRING", "description": "Tanggal akhir YYYY-MM-DD"},
				},
				"required": []string{"start_date", "end_date"},
			},
		},
		{
			Name:        "get_revenue_report",
			Description: "Mendapatkan rincian laporan omzet dan pendapatan berdasarkan rentang tanggal (format YYYY-MM-DD).",
			Parameters: map[string]interface{}{
				"type": "OBJECT",
				"properties": map[string]interface{}{
					"start_date": map[string]interface{}{"type": "STRING", "description": "Tanggal awal YYYY-MM-DD"},
					"end_date":   map[string]interface{}{"type": "STRING", "description": "Tanggal akhir YYYY-MM-DD"},
				},
				"required": []string{"start_date", "end_date"},
			},
		},
		{
			Name:        "get_expense_report",
			Description: "Mendapatkan ringkasan laporan pengeluaran bisnis berdasarkan rentang tanggal dan kategori opsional.",
			Parameters: map[string]interface{}{
				"type": "OBJECT",
				"properties": map[string]interface{}{
					"start_date":  map[string]interface{}{"type": "STRING", "description": "Tanggal awal YYYY-MM-DD"},
					"end_date":    map[string]interface{}{"type": "STRING", "description": "Tanggal akhir YYYY-MM-DD"},
					"category_id": map[string]interface{}{"type": "INTEGER", "description": "ID kategori opsional"},
				},
				"required": []string{"start_date", "end_date"},
			},
		},
		{
			Name:        "get_orders_by_status",
			Description: "Mendapatkan daftar pesanan berdasarkan status pengerjaan (DITERIMA, DICUCI, DIKERINGKAN, DISETRIKA, DIKEMAS, SIAP_DIAMBIL, SELESAI, DIBATALKAN). Maksimal 50 data.",
			Parameters: map[string]interface{}{
				"type": "OBJECT",
				"properties": map[string]interface{}{
					"status": map[string]interface{}{"type": "STRING", "description": "Status pengerjaan"},
					"limit":  map[string]interface{}{"type": "INTEGER", "description": "Maksimal data (default 20, max 50)"},
				},
				"required": []string{"status"},
			},
		},
		{
			Name:        "get_popular_services",
			Description: "Mendapatkan daftar layanan laundry yang paling populer/banyak dipesan beserta frekuensi dan total nilainya.",
			Parameters: map[string]interface{}{
				"type": "OBJECT",
				"properties": map[string]interface{}{
					"start_date": map[string]interface{}{"type": "STRING", "description": "Tanggal awal YYYY-MM-DD"},
					"end_date":   map[string]interface{}{"type": "STRING", "description": "Tanggal akhir YYYY-MM-DD"},
					"limit":      map[string]interface{}{"type": "INTEGER", "description": "Maksimal data (default 5)"},
				},
			},
		},
		{
			Name:        "get_unpaid_orders",
			Description: "Mendapatkan daftar pesanan yang status pembayarannya BELUM_LUNAS. Maksimal 50 data.",
			Parameters: map[string]interface{}{
				"type": "OBJECT",
				"properties": map[string]interface{}{
					"limit": map[string]interface{}{"type": "INTEGER", "description": "Maksimal data (default 20, max 50)"},
				},
			},
		},
		{
			Name:        "get_overdue_orders",
			Description: "Mendapatkan daftar pesanan yang melewati estimasi selesai namun belum berstatus SIAP_DIAMBIL, SELESAI, atau DIBATALKAN. Maksimal 50 data.",
			Parameters: map[string]interface{}{
				"type": "OBJECT",
				"properties": map[string]interface{}{
					"limit": map[string]interface{}{"type": "INTEGER", "description": "Maksimal data (default 20, max 50)"},
				},
			},
		},
		{
			Name:        "get_customer_info",
			Description: "Mendapatkan informasi riwayat dan profil pelanggan berdasarkan ID pelanggan atau nomor telepon.",
			Parameters: map[string]interface{}{
				"type": "OBJECT",
				"properties": map[string]interface{}{
					"customer_id": map[string]interface{}{"type": "INTEGER", "description": "ID pelanggan"},
					"phone":       map[string]interface{}{"type": "STRING", "description": "Nomor telepon pelanggan"},
				},
			},
		},
		{
			Name:        "get_recent_transactions",
			Description: "Mendapatkan daftar transaksi pembayaran kasir terbaru.",
			Parameters: map[string]interface{}{
				"type": "OBJECT",
				"properties": map[string]interface{}{
					"limit": map[string]interface{}{"type": "INTEGER", "description": "Jumlah transaksi (default 10, max 30)"},
				},
			},
		},
		{
			Name:        "get_cash_flow_summary",
			Description: "Mendapatkan perbandingan ringkas uang masuk (pembayaran lunas) vs uang keluar (pengeluaran) untuk periode tertentu.",
			Parameters: map[string]interface{}{
				"type": "OBJECT",
				"properties": map[string]interface{}{
					"start_date": map[string]interface{}{"type": "STRING", "description": "Tanggal awal YYYY-MM-DD"},
					"end_date":   map[string]interface{}{"type": "STRING", "description": "Tanggal akhir YYYY-MM-DD"},
				},
				"required": []string{"start_date", "end_date"},
			},
		},
		{
			Name:        "get_service_list",
			Description: "Mendapatkan daftar semua layanan laundry yang tersedia, tarif harga, satuan, dan durasi pengerjaan.",
			Parameters: map[string]interface{}{
				"type": "OBJECT",
			},
		},
		{
			Name:        "get_expense_categories",
			Description: "Mendapatkan daftar kategori pengeluaran operasional laundry.",
			Parameters: map[string]interface{}{
				"type": "OBJECT",
			},
		},
		{
			Name:        "draft_customer_message",
			Description: "Membuat draf pesan WhatsApp/SMS untuk pelanggan terkait status pesanan (contoh: pesanan selesai, pengingat pengambilan, konfirmasi). TIDAK MENGIRIM PESAN.",
			Parameters: map[string]interface{}{
				"type": "OBJECT",
				"properties": map[string]interface{}{
					"order_id":     map[string]interface{}{"type": "INTEGER", "description": "ID Pesanan"},
					"message_type": map[string]interface{}{"type": "STRING", "description": "Tipe pesan: 'selesai', 'pengingat', 'konfirmasi'"},
				},
				"required": []string{"order_id", "message_type"},
			},
		},
	}
}

// Execute parses tool arguments and runs the appropriate read-only query.
func (e *ToolExecutor) Execute(ctx context.Context, name string, argsJSON string) (interface{}, error) {
	var args map[string]interface{}
	if argsJSON != "" && argsJSON != "{}" {
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			return nil, fmt.Errorf("invalid tool arguments json: %w", err)
		}
	} else {
		args = make(map[string]interface{})
	}

	switch name {
	case "get_business_summary":
		return e.executeBusinessSummary(ctx, args)
	case "get_revenue_report":
		return e.executeRevenueReport(ctx, args)
	case "get_expense_report":
		return e.executeExpenseReport(ctx, args)
	case "get_orders_by_status":
		return e.executeOrdersByStatus(ctx, args)
	case "get_popular_services":
		return e.executePopularServices(ctx, args)
	case "get_unpaid_orders":
		return e.executeUnpaidOrders(ctx, args)
	case "get_overdue_orders":
		return e.executeOverdueOrders(ctx, args)
	case "get_customer_info":
		return e.executeCustomerInfo(ctx, args)
	case "get_recent_transactions":
		return e.executeRecentTransactions(ctx, args)
	case "get_cash_flow_summary":
		return e.executeCashFlowSummary(ctx, args)
	case "get_service_list":
		return e.executeServiceList(ctx)
	case "get_expense_categories":
		return e.executeExpenseCategories(ctx)
	case "draft_customer_message":
		return e.executeDraftCustomerMessage(ctx, args)
	default:
		return nil, fmt.Errorf("tool tidak dikenal: %s", name)
	}
}

func parseDateStrings(args map[string]interface{}) (string, string) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	today := time.Now().In(loc).Format("2006-01-02")

	startStr, _ := args["start_date"].(string)
	endStr, _ := args["end_date"].(string)

	if startStr == "" {
		startStr = today
	}
	if endStr == "" {
		endStr = today
	}

	return startStr, endStr
}

func (e *ToolExecutor) executeBusinessSummary(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	startStr, endStr := parseDateStrings(args)

	revenueSummary, err := e.orderRepo.GetRevenueSummary(ctx, startStr, endStr)
	if err != nil {
		return nil, err
	}

	paymentSummary, err := e.paymentRepo.GetPaymentSummary(ctx, startStr, endStr)
	if err != nil {
		return nil, err
	}

	expense, err := e.expenseRepo.GetExpenseSummary(ctx, startStr, endStr)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"periode_mulai":       startStr,
		"periode_selesai":     endStr,
		"total_omzet_pesanan": revenueSummary.TotalRevenue,
		"total_kas_masuk":     paymentSummary.TotalReceived,
		"total_pengeluaran":   expense,
		"estimasi_laba_kotor": paymentSummary.TotalReceived - expense,
		"total_pesanan":       revenueSummary.TotalOrders,
	}, nil
}

func (e *ToolExecutor) executeRevenueReport(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	startStr, endStr := parseDateStrings(args)

	summary, err := e.orderRepo.GetRevenueSummary(ctx, startStr, endStr)
	if err != nil {
		return nil, err
	}

	daily, err := e.orderRepo.GetDailyRevenue(ctx, startStr, endStr)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"periode_mulai":   startStr,
		"periode_selesai": endStr,
		"total_omzet":     summary.TotalRevenue,
		"total_pesanan":   summary.TotalOrders,
		"total_diskon":    summary.TotalDiscount,
		"rincian_harian":  daily,
	}, nil
}

func (e *ToolExecutor) executeExpenseReport(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	startStr, endStr := parseDateStrings(args)

	var catID uint
	if cid, ok := args["category_id"].(float64); ok && cid > 0 {
		catID = uint(cid)
	}

	total, err := e.expenseRepo.GetExpenseSummary(ctx, startStr, endStr)
	if err != nil {
		return nil, err
	}

	expenses, _, err := e.expenseRepo.FindAll(ctx, repository.ExpenseFilter{
		CategoryID: catID,
		StartDate:  startStr,
		EndDate:    endStr,
	}, 0, 30)
	if err != nil {
		return nil, err
	}

	type expenseItem struct {
		Tanggal    string `json:"tanggal"`
		Kategori   string `json:"kategori"`
		Jumlah     int64  `json:"jumlah"`
		Keterangan string `json:"keterangan"`
	}

	items := make([]expenseItem, 0, len(expenses))
	for _, ex := range expenses {
		catName := ex.Category.Name
		if catName == "" {
			catName = "Lain-lain"
		}
		items = append(items, expenseItem{
			Tanggal:    ex.ExpenseDate.Format("2006-01-02"),
			Kategori:   catName,
			Jumlah:     ex.Amount,
			Keterangan: ex.Description,
		})
	}

	return map[string]interface{}{
		"total_pengeluaran":  total,
		"daftar_pengeluaran": items,
	}, nil
}

func (e *ToolExecutor) executeOrdersByStatus(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	status, _ := args["status"].(string)
	if status == "" {
		return nil, fmt.Errorf("parameter status diperlukan")
	}

	limit := 20
	if l, ok := args["limit"].(float64); ok && l > 0 {
		limit = int(l)
		if limit > 50 {
			limit = 50
		}
	}

	orders, total, err := e.orderRepo.FindAll(ctx, repository.OrderFilter{
		WorkStatus: status,
	}, 0, limit)
	if err != nil {
		return nil, err
	}

	type summaryOrder struct {
		NoPesanan        string `json:"no_pesanan"`
		Pelanggan        string `json:"pelanggan"`
		StatusPengerjaan string `json:"status_pengerjaan"`
		StatusBayar      string `json:"status_bayar"`
		Total            int64  `json:"total"`
		EstimasiSelesai  string `json:"estimasi_selesai"`
	}

	list := make([]summaryOrder, 0, len(orders))
	for _, o := range orders {
		custName := o.Customer.Name
		if custName == "" {
			custName = "-"
		}
		est := "-"
		if !o.EstimatedCompletion.IsZero() {
			est = o.EstimatedCompletion.Format("2006-01-02 15:04")
		}
		list = append(list, summaryOrder{
			NoPesanan:        o.OrderNumber,
			Pelanggan:        custName,
			StatusPengerjaan: o.WorkStatus,
			StatusBayar:      o.PaymentStatus,
			Total:            o.Total,
			EstimasiSelesai:  est,
		})
	}

	return map[string]interface{}{
		"status":          status,
		"total_ditemukan": total,
		"pesanan":         list,
	}, nil
}

func (e *ToolExecutor) executePopularServices(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	startStr, endStr := parseDateStrings(args)

	limit := 5
	if l, ok := args["limit"].(float64); ok && l > 0 {
		limit = int(l)
	}

	topServices, err := e.serviceRepo.FindTopServices(ctx, startStr, endStr, limit)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"periode_mulai":   startStr,
		"periode_selesai": endStr,
		"layanan_populer": topServices,
	}, nil
}

func (e *ToolExecutor) executeUnpaidOrders(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	limit := 20
	if l, ok := args["limit"].(float64); ok && l > 0 {
		limit = int(l)
		if limit > 50 {
			limit = 50
		}
	}

	orders, err := e.orderRepo.FindUnpaidOrders(ctx, limit)
	if err != nil {
		return nil, err
	}

	type unpaidOrder struct {
		NoPesanan string `json:"no_pesanan"`
		Pelanggan string `json:"pelanggan"`
		Telepon   string `json:"telepon"`
		Total     int64  `json:"total"`
		Tanggal   string `json:"tanggal"`
	}

	list := make([]unpaidOrder, 0, len(orders))
	for _, o := range orders {
		cName := o.Customer.Name
		if cName == "" {
			cName = "-"
		}
		cPhone := o.Customer.Phone
		if cPhone == "" {
			cPhone = "-"
		}
		list = append(list, unpaidOrder{
			NoPesanan: o.OrderNumber,
			Pelanggan: cName,
			Telepon:   cPhone,
			Total:     o.Total,
			Tanggal:   o.OrderDate.Format("2006-01-02"),
		})
	}

	return map[string]interface{}{
		"total_belum_lunas": len(list),
		"daftar_pesanan":    list,
	}, nil
}

func (e *ToolExecutor) executeOverdueOrders(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	limit := 20
	if l, ok := args["limit"].(float64); ok && l > 0 {
		limit = int(l)
		if limit > 50 {
			limit = 50
		}
	}

	orders, err := e.orderRepo.FindOverdueOrders(ctx, limit)
	if err != nil {
		return nil, err
	}

	type overdueOrder struct {
		NoPesanan        string `json:"no_pesanan"`
		Pelanggan        string `json:"pelanggan"`
		StatusPengerjaan string `json:"status_pengerjaan"`
		EstimasiSelesai  string `json:"estimasi_selesai"`
	}

	list := make([]overdueOrder, 0, len(orders))
	for _, o := range orders {
		cName := o.Customer.Name
		if cName == "" {
			cName = "-"
		}
		est := "-"
		if !o.EstimatedCompletion.IsZero() {
			est = o.EstimatedCompletion.Format("2006-01-02 15:04")
		}
		list = append(list, overdueOrder{
			NoPesanan:        o.OrderNumber,
			Pelanggan:        cName,
			StatusPengerjaan: o.WorkStatus,
			EstimasiSelesai:  est,
		})
	}

	return map[string]interface{}{
		"total_terlambat": len(list),
		"daftar_pesanan":  list,
	}, nil
}

func (e *ToolExecutor) executeCustomerInfo(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	var cust *model.Customer
	var err error

	if custID, ok := args["customer_id"].(float64); ok && custID > 0 {
		cust, err = e.customerRepo.FindByID(ctx, uint(custID))
	} else if phone, ok := args["phone"].(string); ok && phone != "" {
		cust, err = e.customerRepo.FindByPhone(ctx, phone)
	} else {
		return nil, fmt.Errorf("masukkan customer_id atau phone")
	}

	if err != nil {
		return nil, fmt.Errorf("pelanggan tidak ditemukan")
	}

	// Count customer orders
	orders, _ := e.orderRepo.FindByCustomerID(ctx, cust.ID, 100)

	return map[string]interface{}{
		"id":            cust.ID,
		"nama":          cust.Name,
		"telepon":       cust.Phone,
		"alamat":        cust.Address,
		"total_pesanan": len(orders),
		"catatan":       cust.Notes,
	}, nil
}

func (e *ToolExecutor) executeRecentTransactions(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	limit := 10
	if l, ok := args["limit"].(float64); ok && l > 0 {
		limit = int(l)
		if limit > 30 {
			limit = 30
		}
	}

	payments, err := e.paymentRepo.FindRecent(ctx, limit)
	if err != nil {
		return nil, err
	}

	type paymentItem struct {
		NoPesanan string `json:"no_pesanan"`
		Jumlah    int64  `json:"jumlah"`
		Metode    string `json:"metode"`
		Waktu     string `json:"waktu"`
	}

	list := make([]paymentItem, 0, len(payments))
	for _, p := range payments {
		orderNum := "-"
		if p.Order != nil && p.Order.OrderNumber != "" {
			orderNum = p.Order.OrderNumber
		}
		list = append(list, paymentItem{
			NoPesanan: orderNum,
			Jumlah:    p.Amount,
			Metode:    p.PaymentMethod,
			Waktu:     p.PaidAt.Format("2006-01-02 15:04"),
		})
	}

	return list, nil
}

func (e *ToolExecutor) executeCashFlowSummary(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	return e.executeBusinessSummary(ctx, args)
}

func (e *ToolExecutor) executeServiceList(ctx context.Context) (interface{}, error) {
	services, err := e.serviceRepo.FindActive(ctx)
	if err != nil {
		return nil, err
	}

	type serviceItem struct {
		ID             uint   `json:"id"`
		Nama           string `json:"nama"`
		Tarif          int64  `json:"tarif"`
		Satuan         string `json:"satuan"`
		EstimasiDurasi int    `json:"durasi_jam"`
	}

	list := make([]serviceItem, 0, len(services))
	for _, s := range services {
		list = append(list, serviceItem{
			ID:             s.ID,
			Nama:           s.Name,
			Tarif:          s.PricePerUnit,
			Satuan:         s.Unit,
			EstimasiDurasi: s.EstimatedDuration,
		})
	}

	return list, nil
}

func (e *ToolExecutor) executeExpenseCategories(ctx context.Context) (interface{}, error) {
	categories, err := e.expenseRepo.FindAllCategories(ctx, false)
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (e *ToolExecutor) executeDraftCustomerMessage(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	orderIDVal, ok := args["order_id"].(float64)
	if !ok || orderIDVal <= 0 {
		return nil, fmt.Errorf("parameter order_id wajib diisi")
	}
	messageType, _ := args["message_type"].(string)

	order, err := e.orderRepo.FindByID(ctx, uint(orderIDVal))
	if err != nil {
		return nil, fmt.Errorf("pesanan ID %d tidak ditemukan", uint(orderIDVal))
	}

	custName := order.Customer.Name
	if custName == "" {
		custName = "Pelanggan"
	}

	var draft string
	switch strings.ToLower(messageType) {
	case "selesai":
		draft = fmt.Sprintf("Halo Kak %s, pesanan laundry Anda dengan nomor %s sudah SELESAI dan siap diambil di SatSet Laundry. Total tagihan: Rp %d (%s). Terima kasih!",
			custName, order.OrderNumber, order.Total, order.PaymentStatus)
	case "pengingat":
		draft = fmt.Sprintf("Halo Kak %s, kami mengingatkan bahwa cucian Anda dengan nomor %s sudah selesai dan menunggu untuk diambil di outlet SatSet Laundry. Terima kasih!",
			custName, order.OrderNumber)
	default:
		draft = fmt.Sprintf("Halo Kak %s, pesanan laundry nomor %s saat ini berstatus '%s'. Terima kasih atas kepercayaannya pada SatSet Laundry.",
			custName, order.OrderNumber, order.WorkStatus)
	}

	return map[string]interface{}{
		"pesanan":    order.OrderNumber,
		"pelanggan":  custName,
		"tipe":       messageType,
		"draf_pesan": draft,
		"keterangan": "Draf pesan dibuat (tidak dikirim otomatis ke pelanggan).",
	}, nil
}
