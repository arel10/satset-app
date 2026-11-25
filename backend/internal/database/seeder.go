package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"time"

	"gorm.io/gorm"

	"satset-app/internal/model"
	"satset-app/internal/repository"
)

type sampleOrderItemSpec struct {
	serviceName string
	quantity    float64
}

type sampleOrderSpec struct {
	customerPhone string
	daysAgo       int
	hourOffset    int
	items         []sampleOrderItemSpec
	discount      int64
	notes         string
	workStatus    string
	paid          bool
	paymentMethod string
}

type sampleExpenseSpec struct {
	categoryName string
	amount       int64
	description  string
	daysAgo      int
}

// SeedSampleData populates a compact, realistic set of sample laundry data
// (services, customers, orders with items/history/payments, and expenses).
// It is idempotent: existing services/customers are reused, and sample orders/expenses
// are skipped if records already exist.
func SeedSampleData(ctx context.Context, db *gorm.DB) error {
	var owner, kasir, operator model.User
	if err := db.WithContext(ctx).Where("role = ?", model.RoleOwner).First(&owner).Error; err != nil {
		return fmt.Errorf("akun owner belum tersedia: %w", err)
	}
	if err := db.WithContext(ctx).Where("role = ?", model.RoleKasir).First(&kasir).Error; err != nil {
		kasir = owner
	}
	if err := db.WithContext(ctx).Where("role = ?", model.RoleOperator).First(&operator).Error; err != nil {
		operator = owner
	}

	serviceMap, err := seedServices(ctx, db)
	if err != nil {
		return err
	}

	customerMap, err := seedCustomers(ctx, db)
	if err != nil {
		return err
	}

	if err := seedOrders(ctx, db, serviceMap, customerMap, kasir.ID, operator.ID); err != nil {
		return err
	}

	if err := seedExpenses(ctx, db, owner.ID); err != nil {
		return err
	}

	return nil
}

func seedServices(ctx context.Context, db *gorm.DB) (map[string]model.Service, error) {
	samples := []model.Service{
		{
			Name:              "Cuci Komplit Reguler",
			Category:          "Kiloan",
			Unit:              model.UnitKg,
			PricePerUnit:      7000,
			EstimatedDuration: 48,
			IsActive:          true,
			Description:       "Cuci, kering, setrika, dan lipat rapi (estimasi 2 hari)",
		},
		{
			Name:              "Cuci Komplit Express",
			Category:          "Kiloan",
			Unit:              model.UnitKg,
			PricePerUnit:      12000,
			EstimatedDuration: 12,
			IsActive:          true,
			Description:       "Layanan kilat selesai dalam 12 jam",
		},
		{
			Name:              "Setrika Saja",
			Category:          "Kiloan",
			Unit:              model.UnitKg,
			PricePerUnit:      5000,
			EstimatedDuration: 24,
			IsActive:          true,
			Description:       "Setrika uap rapi dan wangi",
		},
		{
			Name:              "Cuci Bedcover Besar",
			Category:          "Satuan",
			Unit:              model.UnitPcs,
			PricePerUnit:      35000,
			EstimatedDuration: 72,
			IsActive:          true,
			Description:       "Cuci bersih bedcover ukuran Queen/King",
		},
		{
			Name:              "Cuci Jaket / Jas",
			Category:          "Satuan",
			Unit:              model.UnitPcs,
			PricePerUnit:      25000,
			EstimatedDuration: 48,
			IsActive:          true,
			Description:       "Perawatan khusus jaket tebal atau jas",
		},
		{
			Name:              "Cuci Sepatu",
			Category:          "Satuan",
			Unit:              model.UnitPasang,
			PricePerUnit:      30000,
			EstimatedDuration: 48,
			IsActive:          true,
			Description:       "Deep clean sepatu sneakers atau kanvas",
		},
	}

	result := make(map[string]model.Service, len(samples))
	created := 0

	for _, svc := range samples {
		var existing model.Service
		err := db.WithContext(ctx).Where("name = ?", svc.Name).First(&existing).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := db.WithContext(ctx).Create(&svc).Error; err != nil {
				return nil, fmt.Errorf("gagal membuat layanan %q: %w", svc.Name, err)
			}
			result[svc.Name] = svc
			created++
		} else if err != nil {
			return nil, fmt.Errorf("gagal memeriksa layanan %q: %w", svc.Name, err)
		} else {
			result[existing.Name] = existing
		}
	}

	log.Printf("Services: %d baru ditambahkan (%d total siap dipakai)", created, len(result))
	return result, nil
}

func seedCustomers(ctx context.Context, db *gorm.DB) (map[string]model.Customer, error) {
	samples := []model.Customer{
		{
			Name:    "Budi Santoso",
			Phone:   "081234567801",
			Address: "Jl. Melati No. 12, Kebayoran Baru",
			Notes:   "Pelanggan tetap, suka pewangi lavender",
		},
		{
			Name:    "Siti Rahmawati",
			Phone:   "081234567802",
			Address: "Jl. Anggrek Blok B3 No. 5",
			Notes:   "Pisahkan baju putih saat mencuci",
		},
		{
			Name:    "Andi Wijaya",
			Phone:   "081234567803",
			Address: "Apartemen Taman Rasuna Tower 2",
			Notes:   "Biasa ambil sore hari",
		},
		{
			Name:    "Dewi Lestari",
			Phone:   "081234567804",
			Address: "Jl. Kemang Selatan No. 45",
			Notes:   "",
		},
		{
			Name:    "Rizky Pratama",
			Phone:   "081234567805",
			Address: "Jl. Tebet Barat Dalam No. 18",
			Notes:   "Sering pakai layanan express",
		},
		{
			Name:    "Maya Kusuma",
			Phone:   "081234567806",
			Address: "Jl. Pondok Indah Raya No. 88",
			Notes:   "Minta nota dicetak",
		},
	}

	result := make(map[string]model.Customer, len(samples))
	created := 0

	for _, cust := range samples {
		var existing model.Customer
		err := db.WithContext(ctx).Where("phone = ?", cust.Phone).First(&existing).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := db.WithContext(ctx).Create(&cust).Error; err != nil {
				return nil, fmt.Errorf("gagal membuat pelanggan %q: %w", cust.Name, err)
			}
			result[cust.Phone] = cust
			created++
		} else if err != nil {
			return nil, fmt.Errorf("gagal memeriksa pelanggan %q: %w", cust.Name, err)
		} else {
			result[existing.Phone] = existing
		}
	}

	log.Printf("Customers: %d baru ditambahkan (%d total siap dipakai)", created, len(result))
	return result, nil
}

func seedOrders(
	ctx context.Context,
	db *gorm.DB,
	services map[string]model.Service,
	customers map[string]model.Customer,
	kasirID uint,
	operatorID uint,
) error {
	var existingCount int64
	if err := db.WithContext(ctx).Model(&model.Order{}).Count(&existingCount).Error; err != nil {
		return fmt.Errorf("gagal memeriksa jumlah pesanan: %w", err)
	}
	if existingCount > 0 {
		log.Printf("Orders: dilewati (sudah ada %d pesanan di database)", existingCount)
		return nil
	}

	specs := []sampleOrderSpec{
		{
			customerPhone: "081234567801",
			daysAgo:       6,
			hourOffset:    2,
			items: []sampleOrderItemSpec{
				{serviceName: "Cuci Komplit Reguler", quantity: 4.5},
			},
			notes:         "Pewangi lavender",
			workStatus:    model.WorkStatusSelesai,
			paid:          true,
			paymentMethod: model.PaymentMethodTunai,
		},
		{
			customerPhone: "081234567802",
			daysAgo:       5,
			hourOffset:    3,
			items: []sampleOrderItemSpec{
				{serviceName: "Cuci Komplit Reguler", quantity: 3.0},
				{serviceName: "Cuci Bedcover Besar", quantity: 1},
			},
			discount:      5000,
			notes:         "Diskon pelanggan tetap",
			workStatus:    model.WorkStatusSelesai,
			paid:          true,
			paymentMethod: model.PaymentMethodQRIS,
		},
		{
			customerPhone: "081234567803",
			daysAgo:       4,
			hourOffset:    5,
			items: []sampleOrderItemSpec{
				{serviceName: "Cuci Jaket / Jas", quantity: 2},
				{serviceName: "Setrika Saja", quantity: 2.5},
			},
			workStatus:    model.WorkStatusSelesai,
			paid:          true,
			paymentMethod: model.PaymentMethodTransfer,
		},
		{
			customerPhone: "081234567804",
			daysAgo:       2,
			hourOffset:    2,
			items: []sampleOrderItemSpec{
				{serviceName: "Cuci Komplit Reguler", quantity: 5.0},
				{serviceName: "Cuci Sepatu", quantity: 1},
			},
			workStatus:    model.WorkStatusSiapDiambil,
			paid:          true,
			paymentMethod: model.PaymentMethodQRIS,
		},
		{
			customerPhone: "081234567805",
			daysAgo:       2,
			hourOffset:    6,
			items: []sampleOrderItemSpec{
				{serviceName: "Cuci Komplit Express", quantity: 3.5},
			},
			notes:      "Bayar saat ambil",
			workStatus: model.WorkStatusSiapDiambil,
			paid:       false,
		},
		{
			customerPhone: "081234567806",
			daysAgo:       1,
			hourOffset:    3,
			items: []sampleOrderItemSpec{
				{serviceName: "Cuci Komplit Reguler", quantity: 6.0},
			},
			workStatus:    model.WorkStatusDisetrika,
			paid:          true,
			paymentMethod: model.PaymentMethodTunai,
		},
		{
			customerPhone: "081234567801",
			daysAgo:       1,
			hourOffset:    7,
			items: []sampleOrderItemSpec{
				{serviceName: "Cuci Bedcover Besar", quantity: 1},
				{serviceName: "Setrika Saja", quantity: 2.0},
			},
			workStatus: model.WorkStatusDicuci,
			paid:       false,
		},
		{
			customerPhone: "081234567802",
			daysAgo:       0,
			hourOffset:    1,
			items: []sampleOrderItemSpec{
				{serviceName: "Cuci Komplit Express", quantity: 2.5},
				{serviceName: "Cuci Sepatu", quantity: 1},
			},
			notes:      "Butuh cepat untuk besok pagi",
			workStatus: model.WorkStatusDiterima,
			paid:       false,
		},
	}

	orderRepo := repository.NewOrderRepository(db)
	nowUTC := time.Now().UTC()
	baseDayUTC := time.Date(nowUTC.Year(), nowUTC.Month(), nowUTC.Day(), 2, 0, 0, 0, time.UTC)

	for _, spec := range specs {
		cust, ok := customers[spec.customerPhone]
		if !ok {
			return fmt.Errorf("pelanggan dengan nomor %s tidak ditemukan", spec.customerPhone)
		}

		orderTime := baseDayUTC.AddDate(0, 0, -spec.daysAgo).Add(time.Duration(spec.hourOffset) * time.Hour)
		if orderTime.After(nowUTC) {
			orderTime = nowUTC.Add(-30 * time.Minute)
		}

		var items []model.OrderItem
		var subtotal int64
		var maxDuration int

		for _, itemSpec := range spec.items {
			svc, ok := services[itemSpec.serviceName]
			if !ok {
				return fmt.Errorf("layanan %q tidak ditemukan", itemSpec.serviceName)
			}
			itemSubtotal := int64(math.Round(itemSpec.quantity * float64(svc.PricePerUnit)))
			items = append(items, model.OrderItem{
				ServiceID:    svc.ID,
				ServiceName:  svc.Name,
				ServiceUnit:  svc.Unit,
				PricePerUnit: svc.PricePerUnit,
				Quantity:     itemSpec.quantity,
				Subtotal:     itemSubtotal,
				CreatedAt:    orderTime,
			})
			subtotal += itemSubtotal
			if svc.EstimatedDuration > maxDuration {
				maxDuration = svc.EstimatedDuration
			}
		}

		if maxDuration == 0 {
			maxDuration = 24
		}
		total := subtotal - spec.discount
		paymentStatus := model.PaymentStatusBelumLunas
		if spec.paid {
			paymentStatus = model.PaymentStatusLunas
		}

		err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			orderNumber, err := orderRepo.GenerateOrderNumber(ctx, tx, orderTime)
			if err != nil {
				return fmt.Errorf("gagal generate nomor pesanan: %w", err)
			}

			order := &model.Order{
				OrderNumber:         orderNumber,
				CustomerID:          cust.ID,
				OrderDate:           orderTime,
				EstimatedCompletion: orderTime.Add(time.Duration(maxDuration) * time.Hour),
				Subtotal:            subtotal,
				Discount:            spec.discount,
				Total:               total,
				Notes:               spec.notes,
				WorkStatus:          spec.workStatus,
				PaymentStatus:       paymentStatus,
				CreatedBy:           kasirID,
				Items:               items,
				CreatedAt:           orderTime,
				UpdatedAt:           orderTime,
			}

			if err := tx.Create(order).Error; err != nil {
				return fmt.Errorf("gagal menyimpan pesanan %s: %w", orderNumber, err)
			}

			histories := buildStatusHistories(order.ID, spec.workStatus, orderTime, kasirID, operatorID)
			if err := tx.Create(&histories).Error; err != nil {
				return fmt.Errorf("gagal menyimpan riwayat status pesanan %s: %w", orderNumber, err)
			}

			if spec.paid {
				paidAt := orderTime.Add(10 * time.Minute)
				payment := &model.Payment{
					OrderID:       order.ID,
					Amount:        total,
					PaymentMethod: spec.paymentMethod,
					PaidAt:        paidAt,
					ReceivedBy:    kasirID,
					Notes:         "Pembayaran lunas",
					CreatedAt:     paidAt,
					UpdatedAt:     paidAt,
				}
				if err := tx.Create(payment).Error; err != nil {
					return fmt.Errorf("gagal menyimpan pembayaran pesanan %s: %w", orderNumber, err)
				}
			}

			return nil
		})
		if err != nil {
			return err
		}
	}

	log.Printf("Orders: %d pesanan sampel berhasil ditambahkan", len(specs))
	return nil
}

// buildStatusHistories builds sequential status transitions from DITERIMA up to targetStatus.
func buildStatusHistories(orderID uint, targetStatus string, startTime time.Time, kasirID, operatorID uint) []model.OrderStatusHistory {
	histories := []model.OrderStatusHistory{
		{
			OrderID:    orderID,
			FromStatus: "",
			ToStatus:   model.WorkStatusDiterima,
			ChangedBy:  kasirID,
			Notes:      "Pesanan diterima di kasir",
			CreatedAt:  startTime,
		},
	}

	if targetStatus == model.WorkStatusDiterima {
		return histories
	}

	prev := model.WorkStatusDiterima
	stepTime := startTime
	for _, st := range model.ValidWorkStatuses[1:] {
		stepTime = stepTime.Add(2 * time.Hour)
		changer := operatorID
		note := "Status diperbarui oleh operator"
		if st == model.WorkStatusSelesai {
			changer = kasirID
			note = "Pesanan diambil pelanggan"
		}
		histories = append(histories, model.OrderStatusHistory{
			OrderID:    orderID,
			FromStatus: prev,
			ToStatus:   st,
			ChangedBy:  changer,
			Notes:      note,
			CreatedAt:  stepTime,
		})
		prev = st
		if st == targetStatus {
			break
		}
	}

	return histories
}

func seedExpenses(ctx context.Context, db *gorm.DB, ownerID uint) error {
	var existingCount int64
	if err := db.WithContext(ctx).Model(&model.Expense{}).Count(&existingCount).Error; err != nil {
		return fmt.Errorf("gagal memeriksa jumlah pengeluaran: %w", err)
	}
	if existingCount > 0 {
		log.Printf("Expenses: dilewati (sudah ada %d pengeluaran di database)", existingCount)
		return nil
	}

	var categories []model.ExpenseCategory
	if err := db.WithContext(ctx).Find(&categories).Error; err != nil {
		return fmt.Errorf("gagal mengambil kategori pengeluaran: %w", err)
	}
	catMap := make(map[string]uint, len(categories))
	for _, c := range categories {
		catMap[c.Name] = c.ID
	}

	specs := []sampleExpenseSpec{
		{
			categoryName: "Deterjen",
			amount:       85000,
			description:  "Pembelian deterjen cair konsentrat 5L",
			daysAgo:      6,
		},
		{
			categoryName: "Pewangi",
			amount:       65000,
			description:  "Pewangi laundry aroma Sakura & Lavender 5L",
			daysAgo:      5,
		},
		{
			categoryName: "Listrik",
			amount:       350000,
			description:  "Token listrik operasional outlet",
			daysAgo:      3,
		},
		{
			categoryName: "Perlengkapan",
			amount:       45000,
			description:  "Plastik packing berbagai ukuran dan lakban",
			daysAgo:      2,
		},
		{
			categoryName: "Perawatan Mesin",
			amount:       120000,
			description:  "Pembersihan filter dan cek rutin mesin pengering",
			daysAgo:      1,
		},
	}

	nowUTC := time.Now().UTC()
	baseDayUTC := time.Date(nowUTC.Year(), nowUTC.Month(), nowUTC.Day(), 4, 0, 0, 0, time.UTC)

	created := 0
	for _, spec := range specs {
		catID, ok := catMap[spec.categoryName]
		if !ok {
			return fmt.Errorf("kategori pengeluaran %q tidak ditemukan", spec.categoryName)
		}
		expDate := baseDayUTC.AddDate(0, 0, -spec.daysAgo)
		expense := &model.Expense{
			CategoryID:  catID,
			Amount:      spec.amount,
			Description: spec.description,
			ExpenseDate: expDate,
			RecordedBy:  ownerID,
			CreatedAt:   expDate,
			UpdatedAt:   expDate,
		}
		if err := db.WithContext(ctx).Create(expense).Error; err != nil {
			return fmt.Errorf("gagal membuat pengeluaran %q: %w", spec.description, err)
		}
		created++
	}

	log.Printf("Expenses: %d pengeluaran sampel berhasil ditambahkan", created)
	return nil
}
