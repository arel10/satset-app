import React, { useEffect, useState, useRef } from 'react';
import { useNavigate } from 'react-router-dom';
import { servicesApi } from '@/api/services';
import { customersApi } from '@/api/customers';
import { ordersApi } from '@/api/orders';
import { paymentsApi } from '@/api/payments';
import { settingsApi } from '@/api/settings';
import { Customer, LaundryService, Order, CreateOrderItemInput } from '@/types/api';
import { formatRupiah, formatQuantity } from '@/utils/currency';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';
import { Card } from '@/components/ui/Card';
import { Modal } from '@/components/ui/Modal';
import { ThermalReceipt } from '@/components/receipt/ThermalReceipt';
import {
  Search,
  UserPlus,
  Plus,
  Minus,
  Trash2,
  CheckCircle,
  Printer,
  Sparkles,
  ArrowRight,
  AlertCircle,
  X,
} from 'lucide-react';

interface CartItem extends CreateOrderItemInput {
  service: LaundryService;
  subtotal: number;
}

export const POSPage: React.FC = () => {
  const navigate = useNavigate();

  // Master Data
  const [services, setServices] = useState<LaundryService[]>([]);
  const [categories, setCategories] = useState<string[]>([]);
  const [selectedCategory, setSelectedCategory] = useState<string>('all');
  const [isLoadingServices, setIsLoadingServices] = useState(true);

  // Customer State
  const [customerSearch, setCustomerSearch] = useState('');
  const [customerResults, setCustomerResults] = useState<Customer[]>([]);
  const [selectedCustomer, setSelectedCustomer] = useState<Customer | null>(null);
  const [isCustomerDropdownOpen, setIsCustomerDropdownOpen] = useState(false);
  const [isSearchingCustomer, setIsSearchingCustomer] = useState(false);
  const [isNewCustomerModalOpen, setIsNewCustomerModalOpen] = useState(false);
  const [newCustomerForm, setNewCustomerForm] = useState({ name: '', phone: '', address: '', notes: '' });
  const [isSavingCustomer, setIsSavingCustomer] = useState(false);
  const customerDropdownRef = useRef<HTMLDivElement>(null);

  // Cart / Items State
  const [cart, setCart] = useState<CartItem[]>([]);
  const [discount, setDiscount] = useState<number>(0);
  const [notes, setNotes] = useState('');
  const [isSubmittingOrder, setIsSubmittingOrder] = useState(false);
  const [orderError, setOrderError] = useState<string | null>(null);

  // Success & Payment State
  const [createdOrder, setCreatedOrder] = useState<Order | null>(null);
  const [isSuccessModalOpen, setIsSuccessModalOpen] = useState(false);
  const [paymentMethod, setPaymentMethod] = useState<'TUNAI' | 'TRANSFER' | 'QRIS'>('TUNAI');
  const [isProcessingPayment, setIsProcessingPayment] = useState(false);
  const [businessSettings, setBusinessSettings] = useState<any>(null);

  // Load Services & Settings
  useEffect(() => {
    const loadInitData = async () => {
      try {
        const [activeServices, settings] = await Promise.all([
          servicesApi.getActive(),
          settingsApi.getAll(),
        ]);
        setServices(activeServices);
        setBusinessSettings(settings);

        const cats = Array.from(new Set(activeServices.map(s => s.category).filter(Boolean)));
        setCategories(cats);
      } catch (err: any) {
        setOrderError(err.message || 'Gagal memuat katalog layanan.');
      } finally {
        setIsLoadingServices(false);
      }
    };
    loadInitData();
  }, []);

  // Close customer dropdown when clicking outside
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (customerDropdownRef.current && !customerDropdownRef.current.contains(event.target as Node)) {
        setIsCustomerDropdownOpen(false);
      }
    };
    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  // Search Customer debounce
  useEffect(() => {
    if (!customerSearch.trim() || selectedCustomer?.name === customerSearch) {
      setCustomerResults([]);
      setIsSearchingCustomer(false);
      return;
    }
    setIsSearchingCustomer(true);
    const timer = setTimeout(async () => {
      try {
        const res = await customersApi.search(customerSearch.trim(), 1, 8);
        setCustomerResults(res.customers);
        setIsCustomerDropdownOpen(true);
      } catch {
        setCustomerResults([]);
      } finally {
        setIsSearchingCustomer(false);
      }
    }, 250);
    return () => clearTimeout(timer);
  }, [customerSearch, selectedCustomer]);

  // Cart calculations
  const subtotal = cart.reduce((sum, item) => sum + item.subtotal, 0);
  const safeDiscount = Math.min(Math.max(0, discount), subtotal);
  const total = subtotal - safeDiscount;

  // Add Item to Cart
  const handleAddToCart = (service: LaundryService) => {
    setCart(prev => {
      const existing = prev.find(item => item.service_id === service.id);
      if (existing) {
        const updatedQty = existing.quantity + 1;
        const updatedSubtotal = Math.round(updatedQty * service.price_per_unit);
        return prev.map(item =>
          item.service_id === service.id
            ? { ...item, quantity: updatedQty, subtotal: updatedSubtotal }
            : item
        );
      }
      return [
        ...prev,
        {
          service_id: service.id,
          service,
          quantity: 1,
          subtotal: service.price_per_unit,
        },
      ];
    });
  };

  // Update Item Quantity
  const handleUpdateQty = (serviceID: number, newQty: number) => {
    if (newQty <= 0) {
      handleRemoveItem(serviceID);
      return;
    }
    setCart(prev =>
      prev.map(item => {
        if (item.service_id === serviceID) {
          const sub = Math.round(newQty * item.service.price_per_unit);
          return { ...item, quantity: newQty, subtotal: sub };
        }
        return item;
      })
    );
  };

  // Remove Item
  const handleRemoveItem = (serviceID: number) => {
    setCart(prev => prev.filter(item => item.service_id !== serviceID));
  };

  // Save New Customer Modal
  const handleSaveNewCustomer = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newCustomerForm.name.trim() || !newCustomerForm.phone.trim()) return;

    setIsSavingCustomer(true);
    try {
      const created = await customersApi.create(newCustomerForm);
      setSelectedCustomer(created);
      setCustomerSearch(created.name);
      setIsNewCustomerModalOpen(false);
      setNewCustomerForm({ name: '', phone: '', address: '', notes: '' });
    } catch (err: any) {
      alert(err.message || 'Gagal menambahkan pelanggan.');
    } finally {
      setIsSavingCustomer(false);
    }
  };

  // Submit Order
  const handleSubmitOrder = async () => {
    if (!selectedCustomer) {
      setOrderError('Silakan pilih atau tambahkan pelanggan terlebih dahulu.');
      return;
    }
    if (cart.length === 0) {
      setOrderError('Tambahkan minimal 1 item layanan ke dalam pesanan.');
      return;
    }

    setOrderError(null);
    setIsSubmittingOrder(true);

    try {
      const idempotencyKey = 'pos-' + Date.now() + '-' + Math.random().toString(36).substring(2, 9);
      const payload = {
        customer_id: selectedCustomer.id,
        items: cart.map(item => ({
          service_id: item.service_id,
          quantity: item.quantity,
        })),
        discount: safeDiscount,
        notes: notes.trim(),
      };

      const order = await ordersApi.create(payload, idempotencyKey);
      setCreatedOrder(order);
      setIsSuccessModalOpen(true);
    } catch (err: any) {
      setOrderError(err.message || 'Gagal membuat pesanan.');
    } finally {
      setIsSubmittingOrder(false);
    }
  };

  // Immediate Full Payment Process
  const handlePayNow = async () => {
    if (!createdOrder) return;
    setIsProcessingPayment(true);
    try {
      await paymentsApi.record(createdOrder.id, createdOrder.total, paymentMethod);
      const updated = await ordersApi.getByID(createdOrder.id);
      setCreatedOrder(updated);
    } catch (err: any) {
      alert(err.message || 'Gagal mencatat pembayaran.');
    } finally {
      setIsProcessingPayment(false);
    }
  };

  // Reset form for next order
  const handleResetForNewOrder = () => {
    setCart([]);
    setSelectedCustomer(null);
    setCustomerSearch('');
    setDiscount(0);
    setNotes('');
    setCreatedOrder(null);
    setIsSuccessModalOpen(false);
  };

  const filteredServices =
    selectedCategory === 'all'
      ? services
      : services.filter(s => s.category === selectedCategory);

  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-slate-900 tracking-tight">Kasir &amp; Buat Pesanan (POS)</h1>
          <p className="text-xs text-slate-500 mt-1">
            Pencatatan pesanan kilat dengan snapshot harga otomatis dan cetak nota.
          </p>
        </div>
      </div>

      {orderError && (
        <div className="rounded-lg bg-rose-50 border border-rose-200 p-3.5 flex items-start gap-2.5 text-xs text-rose-800">
          <AlertCircle className="w-4 h-4 text-rose-600 shrink-0 mt-0.5" />
          <span>{orderError}</span>
        </div>
      )}

      {/* POS Two-Column Layout */}
      <div className="grid grid-cols-1 lg:grid-cols-12 gap-6 items-start">
        {/* Left Column: Services Catalog (7 Cols) */}
        <div className="lg:col-span-7 space-y-4">
          {/* Customer Selection Card */}
          <div className="bg-white rounded-xl border border-slate-200/80 shadow-xs p-4 relative z-30">
            <div className="flex items-center justify-between mb-2">
              <label className="text-xs font-semibold text-slate-700 flex items-center gap-1.5">
                <span>Pelanggan</span>
                {selectedCustomer && (
                  <span className="text-[10px] font-bold text-teal-700 bg-teal-50 px-1.5 py-0.5 rounded border border-teal-200">
                    Terpilih
                  </span>
                )}
              </label>
              <button
                type="button"
                onClick={() => {
                  setNewCustomerForm({ name: customerSearch.trim(), phone: '', address: '', notes: '' });
                  setIsNewCustomerModalOpen(true);
                }}
                className="text-xs font-semibold text-teal-700 hover:text-teal-800 flex items-center gap-1"
              >
                <UserPlus className="w-3.5 h-3.5" />
                Tambah Baru
              </button>
            </div>

            <div className="relative" ref={customerDropdownRef}>
              <Input
                placeholder="Cari nama atau no. telepon pelanggan..."
                value={customerSearch}
                onChange={e => {
                  setCustomerSearch(e.target.value);
                  if (selectedCustomer) setSelectedCustomer(null);
                  setIsCustomerDropdownOpen(true);
                }}
                onFocus={() => {
                  if (customerResults.length > 0 && !selectedCustomer) {
                    setIsCustomerDropdownOpen(true);
                  }
                }}
                leftIcon={<Search className="w-4 h-4 text-slate-400" />}
                rightIcon={
                  customerSearch ? (
                    <button
                      type="button"
                      onClick={() => {
                        setCustomerSearch('');
                        setSelectedCustomer(null);
                        setCustomerResults([]);
                        setIsCustomerDropdownOpen(false);
                      }}
                      className="p-1 text-slate-400 hover:text-slate-600 pointer-events-auto"
                    >
                      <X className="w-3.5 h-3.5" />
                    </button>
                  ) : null
                }
              />

              {/* Autocomplete Dropdown with robust z-50 and no clipping */}
              {isCustomerDropdownOpen && customerSearch.trim() && !selectedCustomer && (
                <div className="absolute top-full left-0 right-0 mt-1.5 bg-white border border-slate-200/90 rounded-xl shadow-2xl z-50 overflow-hidden divide-y divide-slate-100 max-h-64 overflow-y-auto">
                  {isSearchingCustomer ? (
                    <div className="p-4 text-center text-xs text-slate-500 flex items-center justify-center gap-2">
                      <div className="w-3.5 h-3.5 border-2 border-teal-600 border-t-transparent rounded-full animate-spin" />
                      <span>Mencari pelanggan...</span>
                    </div>
                  ) : customerResults.length > 0 ? (
                    <>
                      <div className="px-3.5 py-1.5 bg-slate-50 text-[10px] font-bold text-slate-400 uppercase tracking-wider flex items-center justify-between">
                        <span>Hasil Pencarian ({customerResults.length})</span>
                        <span className="text-[10px] text-teal-700 font-normal">Klik untuk memilih</span>
                      </div>
                      {customerResults.map(c => (
                        <button
                          key={c.id}
                          type="button"
                          onClick={() => {
                            setSelectedCustomer(c);
                            setCustomerSearch(c.name);
                            setIsCustomerDropdownOpen(false);
                          }}
                          className="w-full text-left px-3.5 py-2.5 hover:bg-teal-50/70 transition-colors flex items-center justify-between gap-3 group"
                        >
                          <div className="min-w-0 pr-2">
                            <div className="flex items-center gap-2">
                              <span className="font-bold text-slate-900 text-xs truncate group-hover:text-teal-900">
                                {c.name}
                              </span>
                              <span className="text-[11px] font-mono text-slate-600 bg-slate-100 px-1.5 py-0.5 rounded border border-slate-200/60 shrink-0">
                                {c.phone}
                              </span>
                            </div>
                            {c.address && (
                              <p className="text-slate-500 text-[11px] truncate mt-0.5 max-w-sm">
                                {c.address}
                              </p>
                            )}
                          </div>
                          <span className="text-[10px] text-teal-700 bg-teal-50 border border-teal-200 px-2.5 py-1 rounded-md font-semibold shrink-0 group-hover:bg-teal-700 group-hover:text-white transition-colors">
                            Pilih
                          </span>
                        </button>
                      ))}
                    </>
                  ) : (
                    <div className="p-4 text-center">
                      <p className="text-xs font-medium text-slate-700">Pelanggan tidak ditemukan</p>
                      <p className="text-[11px] text-slate-400 mt-0.5">
                        Belum ada pelanggan dengan nama/telepon "{customerSearch.trim()}".
                      </p>
                      <button
                        type="button"
                        onClick={() => {
                          setIsCustomerDropdownOpen(false);
                          setNewCustomerForm({
                            name: customerSearch.trim(),
                            phone: '',
                            address: '',
                            notes: '',
                          });
                          setIsNewCustomerModalOpen(true);
                        }}
                        className="mt-2.5 text-xs font-semibold text-teal-700 hover:text-teal-900 bg-teal-50 hover:bg-teal-100 border border-teal-200 px-3 py-1.5 rounded-lg inline-flex items-center gap-1.5 transition-colors"
                      >
                        <UserPlus className="w-3.5 h-3.5" />
                        Daftarkan Pelanggan Ini
                      </button>
                    </div>
                  )}
                </div>
              )}
            </div>

            {selectedCustomer && (
              <div className="mt-3 p-3 rounded-xl bg-teal-50/80 border border-teal-200 flex items-center justify-between gap-3 text-xs">
                <div className="flex items-center gap-3 min-w-0">
                  <div className="w-8 h-8 rounded-full bg-teal-700 text-white font-bold text-xs flex items-center justify-center shrink-0">
                    {selectedCustomer.name.charAt(0).toUpperCase()}
                  </div>
                  <div className="min-w-0">
                    <div className="flex items-center gap-2">
                      <span className="font-bold text-slate-900 text-xs truncate">
                        {selectedCustomer.name}
                      </span>
                      <span className="text-[11px] font-mono text-teal-800 bg-white px-1.5 py-0.5 rounded border border-teal-200 shrink-0">
                        {selectedCustomer.phone}
                      </span>
                    </div>
                    {selectedCustomer.address && (
                      <p className="text-[11px] text-slate-600 truncate max-w-sm mt-0.5">
                        {selectedCustomer.address}
                      </p>
                    )}
                  </div>
                </div>
                <button
                  type="button"
                  onClick={() => {
                    setSelectedCustomer(null);
                    setCustomerSearch('');
                    setIsCustomerDropdownOpen(false);
                  }}
                  className="text-xs font-semibold text-teal-700 hover:text-teal-900 bg-white px-2.5 py-1 rounded-md border border-teal-200 shadow-2xs shrink-0"
                >
                  Ganti
                </button>
              </div>
            )}
          </div>

          {/* Service Categories Filter */}
          <div className="flex items-center gap-1.5 overflow-x-auto pb-1">
            <button
              onClick={() => setSelectedCategory('all')}
              className={`px-3 py-1.5 rounded-lg text-xs font-medium transition-colors shrink-0 ${
                selectedCategory === 'all'
                  ? 'bg-teal-700 text-white'
                  : 'bg-white border border-slate-200 text-slate-600 hover:bg-slate-50'
              }`}
            >
              Semua Layanan
            </button>
            {categories.map(cat => (
              <button
                key={cat}
                onClick={() => setSelectedCategory(cat)}
                className={`px-3 py-1.5 rounded-lg text-xs font-medium transition-colors shrink-0 ${
                  selectedCategory === cat
                    ? 'bg-teal-700 text-white'
                    : 'bg-white border border-slate-200 text-slate-600 hover:bg-slate-50'
                }`}
              >
                {cat}
              </button>
            ))}
          </div>

          {/* Service Cards Grid */}
          <div className="grid grid-cols-2 sm:grid-cols-3 gap-3">
            {isLoadingServices ? (
              Array.from({ length: 6 }).map((_, i) => (
                <div key={i} className="h-24 bg-slate-100 rounded-xl animate-pulse" />
              ))
            ) : filteredServices.length > 0 ? (
              filteredServices.map(service => (
                <button
                  key={service.id}
                  onClick={() => handleAddToCart(service)}
                  type="button"
                  className="text-left p-3.5 rounded-xl border border-slate-200/90 bg-white hover:border-teal-600 hover:shadow-xs transition-all flex flex-col justify-between h-28 group"
                >
                  <div>
                    <span className="text-[10px] font-semibold text-slate-400 uppercase tracking-wider block">
                      {service.category || 'Reguler'}
                    </span>
                    <h4 className="text-xs font-semibold text-slate-900 group-hover:text-teal-700 line-clamp-2 mt-0.5">
                      {service.name}
                    </h4>
                  </div>
                  <div className="flex items-baseline justify-between border-t border-slate-100 pt-2 mt-2">
                    <span className="text-xs font-bold text-teal-800">
                      {formatRupiah(service.price_per_unit)}
                    </span>
                    <span className="text-[10px] text-slate-400">/{service.unit}</span>
                  </div>
                </button>
              ))
            ) : (
              <div className="col-span-full py-8 text-center text-xs text-slate-400 bg-white rounded-xl border border-dashed border-slate-200">
                Tidak ada layanan pada kategori ini.
              </div>
            )}
          </div>
        </div>

        {/* Right Column: Order Cart & Summary (5 Cols) */}
        <div className="lg:col-span-5">
          <Card className="sticky top-20">
            <h3 className="font-bold text-slate-900 text-sm border-b border-slate-100 pb-3 flex items-center justify-between">
              <span>Keranjang Cucian</span>
              <span className="text-xs font-normal text-slate-500">{cart.length} item</span>
            </h3>

            {/* Cart Items List */}
            <div className="py-3 divide-y divide-slate-100 max-h-[320px] overflow-y-auto">
              {cart.length > 0 ? (
                cart.map(item => (
                  <div key={item.service_id} className="py-2.5 first:pt-0 last:pb-0 space-y-1.5">
                    <div className="flex items-start justify-between">
                      <span className="text-xs font-semibold text-slate-800">
                        {item.service.name}
                      </span>
                      <span className="text-xs font-bold text-slate-900">
                        {formatRupiah(item.subtotal)}
                      </span>
                    </div>

                    <div className="flex items-center justify-between text-xs text-slate-500">
                      <span className="text-[11px]">
                        {formatRupiah(item.service.price_per_unit)} / {item.service.unit}
                      </span>

                      {/* Quantity Stepper */}
                      <div className="flex items-center gap-1.5 bg-slate-50 border border-slate-200 rounded-lg p-0.5">
                        <button
                          type="button"
                          onClick={() => handleUpdateQty(item.service_id, item.quantity - 1)}
                          className="w-6 h-6 rounded bg-white text-slate-600 hover:bg-slate-200 flex items-center justify-center"
                        >
                          <Minus className="w-3 h-3" />
                        </button>
                        <input
                          type="number"
                          step="0.1"
                          min="0.1"
                          value={item.quantity}
                          onChange={e =>
                            handleUpdateQty(item.service_id, parseFloat(e.target.value) || 0)
                          }
                          className="w-12 text-center text-xs font-bold bg-transparent focus:outline-none"
                        />
                        <button
                          type="button"
                          onClick={() => handleUpdateQty(item.service_id, item.quantity + 1)}
                          className="w-6 h-6 rounded bg-white text-slate-600 hover:bg-slate-200 flex items-center justify-center"
                        >
                          <Plus className="w-3 h-3" />
                        </button>
                        <button
                          type="button"
                          onClick={() => handleRemoveItem(item.service_id)}
                          className="w-6 h-6 text-slate-400 hover:text-rose-600 flex items-center justify-center ml-1"
                        >
                          <Trash2 className="w-3.5 h-3.5" />
                        </button>
                      </div>
                    </div>
                  </div>
                ))
              ) : (
                <div className="py-10 text-center text-slate-400 text-xs">
                  Keranjang masih kosong. Klik layanan di sebelah kiri untuk menambahkan.
                </div>
              )}
            </div>

            {/* Calculations & Discounts */}
            <div className="border-t border-slate-100 pt-3 space-y-2 text-xs">
              <div className="flex justify-between text-slate-600">
                <span>Subtotal</span>
                <span>{formatRupiah(subtotal)}</span>
              </div>

              <div className="flex items-center justify-between gap-4">
                <span className="text-slate-600 shrink-0">Diskon (Rp)</span>
                <input
                  type="number"
                  min="0"
                  max={subtotal}
                  value={discount || ''}
                  placeholder="0"
                  onChange={e => setDiscount(parseInt(e.target.value, 10) || 0)}
                  className="w-28 text-right text-xs rounded-md border border-slate-200 px-2 py-1 focus:outline-none focus:ring-1 focus:ring-teal-600"
                />
              </div>

              <div className="flex justify-between items-baseline pt-2 border-t border-slate-100">
                <span className="text-sm font-bold text-slate-900">Total Tagihan</span>
                <span className="text-lg font-bold text-teal-800">{formatRupiah(total)}</span>
              </div>

              <div className="pt-2">
                <Input
                  label="Catatan Khusus Pesanan"
                  placeholder="Contoh: Pisahkan pakaian putih / pewangi ekstra"
                  value={notes}
                  onChange={e => setNotes(e.target.value)}
                />
              </div>

              <Button
                variant="primary"
                size="lg"
                className="w-full mt-4"
                disabled={cart.length === 0 || !selectedCustomer}
                isLoading={isSubmittingOrder}
                onClick={handleSubmitOrder}
              >
                Simpan &amp; Lanjutkan Pesanan
              </Button>
            </div>
          </Card>
        </div>
      </div>

      {/* Modal: Tambah Pelanggan Baru */}
      <Modal
        isOpen={isNewCustomerModalOpen}
        onClose={() => setIsNewCustomerModalOpen(false)}
        title="Tambah Pelanggan Baru"
        description="Data pelanggan akan langsung dipilih untuk pesanan saat ini."
      >
        <form onSubmit={handleSaveNewCustomer} className="space-y-3.5">
          <Input
            label="Nama Lengkap"
            required
            autoFocus
            value={newCustomerForm.name}
            onChange={e => setNewCustomerForm({ ...newCustomerForm, name: e.target.value })}
            placeholder="Contoh: Ibu Rina"
          />
          <Input
            label="Nomor WhatsApp / HP"
            required
            type="tel"
            value={newCustomerForm.phone}
            onChange={e => setNewCustomerForm({ ...newCustomerForm, phone: e.target.value })}
            placeholder="Contoh: 081234567890"
          />
          <Input
            label="Alamat (Opsional)"
            value={newCustomerForm.address}
            onChange={e => setNewCustomerForm({ ...newCustomerForm, address: e.target.value })}
            placeholder="Contoh: Jl. Mawar No. 12"
          />
          <Input
            label="Catatan Tambahan (Opsional)"
            value={newCustomerForm.notes}
            onChange={e => setNewCustomerForm({ ...newCustomerForm, notes: e.target.value })}
            placeholder="Contoh: Langganan bulanan"
          />
          <div className="flex justify-end gap-2 pt-2">
            <Button
              type="button"
              variant="outline"
              onClick={() => setIsNewCustomerModalOpen(false)}
            >
              Batal
            </Button>
            <Button type="submit" variant="primary" isLoading={isSavingCustomer}>
              Simpan Pelanggan
            </Button>
          </div>
        </form>
      </Modal>

      {/* Modal: Sukses Pesanan & Opsi Pembayaran / Cetak Nota */}
      <Modal
        isOpen={isSuccessModalOpen}
        onClose={() => setIsSuccessModalOpen(false)}
        title="Pesanan Berhasil Dibuat!"
        maxWidth="lg"
      >
        {createdOrder && (
          <div className="space-y-4">
            <div className="p-4 rounded-xl bg-teal-50 border border-teal-200 text-center space-y-1">
              <p className="text-xs text-teal-700 font-medium">Nomor Nota Pesanan</p>
              <p className="text-2xl font-bold text-teal-900 tracking-wider">
                {createdOrder.order_number}
              </p>
              <p className="text-xs text-teal-800">
                Total: <span className="font-bold">{formatRupiah(createdOrder.total)}</span> •{' '}
                <span className="font-semibold">{createdOrder.payment_status}</span>
              </p>
            </div>

            {/* If Unpaid, offer immediate full payment */}
            {createdOrder.payment_status === 'BELUM_LUNAS' && (
              <div className="p-4 rounded-xl border border-slate-200 bg-white space-y-3">
                <h4 className="text-xs font-bold text-slate-800 uppercase tracking-wider">
                  Catat Pembayaran Lunas Sekarang
                </h4>
                <div className="grid grid-cols-3 gap-2">
                  {(['TUNAI', 'TRANSFER', 'QRIS'] as const).map(m => (
                    <button
                      key={m}
                      type="button"
                      onClick={() => setPaymentMethod(m)}
                      className={`py-2 px-3 text-xs font-semibold rounded-lg border transition-all ${
                        paymentMethod === m
                          ? 'border-teal-700 bg-teal-50 text-teal-800 ring-1 ring-teal-700'
                          : 'border-slate-200 text-slate-600 hover:bg-slate-50'
                      }`}
                    >
                      {m}
                    </button>
                  ))}
                </div>
                <Button
                  variant="primary"
                  className="w-full"
                  isLoading={isProcessingPayment}
                  onClick={handlePayNow}
                >
                  Bayar Lunas Sekarang ({formatRupiah(createdOrder.total)})
                </Button>
              </div>
            )}

            {/* Thermal Receipt Preview */}
            <div className="border border-slate-200 rounded-xl p-3 bg-slate-50 max-h-60 overflow-y-auto">
              <ThermalReceipt order={createdOrder} businessSettings={businessSettings} />
            </div>

            {/* Bottom Actions */}
            <div className="flex flex-col sm:flex-row items-center justify-between gap-2.5 pt-2">
              <Button
                variant="outline"
                className="w-full sm:w-auto"
                leftIcon={<Printer className="w-4 h-4" />}
                onClick={() => window.print()}
              >
                Cetak Struk Thermal
              </Button>
              <div className="flex items-center gap-2 w-full sm:w-auto">
                <Button
                  variant="secondary"
                  className="w-full sm:w-auto"
                  onClick={() => navigate(`/orders/${createdOrder.id}`)}
                >
                  Lihat Detail
                </Button>
                <Button
                  variant="primary"
                  className="w-full sm:w-auto"
                  onClick={handleResetForNewOrder}
                  rightIcon={<ArrowRight className="w-4 h-4" />}
                >
                  Transaksi Baru
                </Button>
              </div>
            </div>
          </div>
        )}
      </Modal>
    </div>
  );
};
