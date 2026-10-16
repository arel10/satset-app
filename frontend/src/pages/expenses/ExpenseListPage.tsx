import React, { useEffect, useState } from 'react';
import { expensesApi, ExpenseFilterParams } from '@/api/expenses';
import { Expense, ExpenseCategory } from '@/types/api';
import { formatRupiah } from '@/utils/currency';
import { formatJakartaDate } from '@/utils/date';
import { Card } from '@/components/ui/Card';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';
import { Select } from '@/components/ui/Select';
import { Modal } from '@/components/ui/Modal';
import { Skeleton, EmptyState } from '@/components/ui/EmptyState';
import { Receipt, Plus, Ban, Tag, Calendar, ChevronLeft, ChevronRight } from 'lucide-react';

export const ExpenseListPage: React.FC = () => {
  const [expenses, setExpenses] = useState<Expense[]>([]);
  const [categories, setCategories] = useState<ExpenseCategory[]>([]);
  const [total, setTotal] = useState<number>(0);
  const [page, setPage] = useState<number>(1);
  const [perPage] = useState<number>(15);
  const [totalPages, setTotalPages] = useState<number>(1);

  // Filters
  const [selectedCategory, setSelectedCategory] = useState<string>('');
  const [startDate, setStartDate] = useState<string>('');
  const [endDate, setEndDate] = useState<string>('');

  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Add Expense Modal
  const [isAddModalOpen, setIsAddModalOpen] = useState(false);
  const [formExpense, setFormExpense] = useState({
    category_id: 0,
    amount: 50000,
    expense_date: new Date().toISOString().split('T')[0],
    description: '',
    vendor: '',
  });
  const [isSubmittingExpense, setIsSubmittingExpense] = useState(false);

  // Manage Categories Modal
  const [isCategoryModalOpen, setIsCategoryModalOpen] = useState(false);
  const [newCatName, setNewCatName] = useState('');
  const [newCatDesc, setNewCatDesc] = useState('');
  const [isSubmittingCat, setIsSubmittingCat] = useState(false);

  // Cancel Expense Modal
  const [cancelExpenseID, setCancelExpenseID] = useState<number | null>(null);
  const [cancelReason, setCancelReason] = useState('');
  const [isCancelling, setIsCancelling] = useState(false);

  const fetchInitData = async () => {
    setIsLoading(true);
    setError(null);
    try {
      const [expRes, catList] = await Promise.all([
        expensesApi.list({
          page,
          per_page: perPage,
          category_id: selectedCategory ? parseInt(selectedCategory, 10) : undefined,
          start_date: startDate || undefined,
          end_date: endDate || undefined,
        }),
        expensesApi.listCategories(),
      ]);

      setExpenses(expRes.expenses);
      setCategories(catList);
      if (expRes.meta) {
        setTotal(expRes.meta.total);
        setTotalPages(expRes.meta.total_pages);
      }
      if (catList.length > 0 && formExpense.category_id === 0) {
        setFormExpense(prev => ({ ...prev, category_id: catList[0].id }));
      }
    } catch (err: any) {
      setError(err.message || 'Gagal memuat data pengeluaran.');
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchInitData();
  }, [page, selectedCategory, startDate, endDate]);

  const handleAddExpense = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!formExpense.description.trim() || formExpense.amount <= 0 || !formExpense.category_id) {
      alert('Lengkapi kategori, nominal, dan keterangan pengeluaran.');
      return;
    }

    setIsSubmittingExpense(true);
    try {
      await expensesApi.create(formExpense);
      setIsAddModalOpen(false);
      setFormExpense({
        category_id: categories[0]?.id || 0,
        amount: 50000,
        expense_date: new Date().toISOString().split('T')[0],
        description: '',
        vendor: '',
      });
      fetchInitData();
    } catch (err: any) {
      alert(err.message || 'Gagal mencatat pengeluaran.');
    } finally {
      setIsSubmittingExpense(false);
    }
  };

  const handleCreateCategory = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newCatName.trim()) return;

    setIsSubmittingCat(true);
    try {
      await expensesApi.createCategory({ name: newCatName.trim(), description: newCatDesc.trim() });
      setNewCatName('');
      setNewCatDesc('');
      const updatedCats = await expensesApi.listCategories();
      setCategories(updatedCats);
    } catch (err: any) {
      alert(err.message || 'Gagal membuat kategori.');
    } finally {
      setIsSubmittingCat(false);
    }
  };

  const handleCancelExpense = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!cancelExpenseID || !cancelReason.trim()) return;

    setIsCancelling(true);
    try {
      await expensesApi.cancel(cancelExpenseID, cancelReason.trim());
      setCancelExpenseID(null);
      setCancelReason('');
      fetchInitData();
    } catch (err: any) {
      alert(err.message || 'Gagal membatalkan pengeluaran.');
    } finally {
      setIsCancelling(false);
    }
  };

  const totalFilteredAmount = expenses
    .filter(e => !e.is_cancelled)
    .reduce((sum, e) => sum + e.amount, 0);

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-slate-900 tracking-tight">Pengeluaran Operasional</h1>
          <p className="text-xs text-slate-500 mt-1">
            Catat beban biaya deterjen, listrik, air, sewa, dan operasional toko.
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button
            variant="outline"
            onClick={() => setIsCategoryModalOpen(true)}
            leftIcon={<Tag className="w-4 h-4" />}
          >
            Kelola Kategori
          </Button>
          <Button
            variant="primary"
            onClick={() => setIsAddModalOpen(true)}
            leftIcon={<Plus className="w-4 h-4" />}
          >
            Catat Pengeluaran
          </Button>
        </div>
      </div>

      {/* Filter Bar */}
      <Card className="p-4">
        <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
          <Select
            label="Kategori"
            value={selectedCategory}
            onChange={e => {
              setSelectedCategory(e.target.value);
              setPage(1);
            }}
            options={[
              { value: '', label: 'Semua Kategori' },
              ...categories.map(c => ({ value: c.id, label: c.name })),
            ]}
          />

          <Input
            type="date"
            label="Tanggal Mulai"
            value={startDate}
            onChange={e => {
              setStartDate(e.target.value);
              setPage(1);
            }}
          />

          <Input
            type="date"
            label="Tanggal Akhir"
            value={endDate}
            onChange={e => {
              setEndDate(e.target.value);
              setPage(1);
            }}
          />
        </div>
      </Card>

      {/* Table Card */}
      <Card className="overflow-hidden p-0">
        <div className="px-5 py-3 border-b border-slate-100 bg-slate-50/70 flex items-center justify-between text-xs font-semibold">
          <span>Daftar Pengeluaran</span>
          <span className="text-slate-700">
            Total Halaman Ini: <span className="font-bold text-slate-900">{formatRupiah(totalFilteredAmount)}</span>
          </span>
        </div>

        {isLoading ? (
          <div className="p-6 space-y-3">
            {Array.from({ length: 5 }).map((_, i) => (
              <Skeleton key={i} className="h-12 w-full" />
            ))}
          </div>
        ) : error ? (
          <div className="p-8 text-center text-xs text-rose-700 bg-rose-50">
            {error}
          </div>
        ) : expenses.length === 0 ? (
          <div className="p-8">
            <EmptyState
              icon={<Receipt className="w-6 h-6 text-slate-400" />}
              title="Belum Ada Pengeluaran Tercatat"
              description="Catat biaya operasional harian untuk menghitung laba bersih laundry secara akurat."
              action={
                <Button variant="primary" size="sm" onClick={() => setIsAddModalOpen(true)}>
                  Catat Pengeluaran
                </Button>
              }
            />
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-left text-xs divide-y divide-slate-200">
              <thead className="bg-slate-50 text-[11px] font-bold text-slate-700 uppercase tracking-wider">
                <tr>
                  <th className="px-5 py-3.5">Tanggal</th>
                  <th className="px-5 py-3.5">Kategori</th>
                  <th className="px-5 py-3.5">Keterangan</th>
                  <th className="px-5 py-3.5">Vendor / Toko</th>
                  <th className="px-5 py-3.5 text-right">Jumlah</th>
                  <th className="px-5 py-3.5 text-right">Aksi</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-100 bg-white">
                {expenses.map(exp => (
                  <tr
                    key={exp.id}
                    className={`transition-colors ${
                      exp.is_cancelled ? 'bg-slate-50/80 opacity-60' : 'hover:bg-slate-50/80'
                    }`}
                  >
                    <td className="px-5 py-3.5 font-medium text-slate-900 whitespace-nowrap">
                      {formatJakartaDate(exp.expense_date)}
                    </td>
                    <td className="px-5 py-3.5">
                      <span className="bg-slate-100 text-slate-700 px-2 py-0.5 rounded text-[11px] font-medium">
                        {exp.category?.name || 'Umum'}
                      </span>
                    </td>
                    <td className="px-5 py-3.5 text-slate-700 max-w-sm">
                      <span className={`block font-semibold ${exp.is_cancelled ? 'line-through' : ''}`}>
                        {exp.description}
                      </span>
                      {exp.is_cancelled && (
                        <span className="text-[10px] text-rose-600 block mt-0.5">
                          Dibatalkan: {exp.cancel_reason}
                        </span>
                      )}
                    </td>
                    <td className="px-5 py-3.5 text-slate-500">{exp.vendor || '-'}</td>
                    <td
                      className={`px-5 py-3.5 text-right font-bold whitespace-nowrap ${
                        exp.is_cancelled ? 'line-through text-slate-400' : 'text-slate-900'
                      }`}
                    >
                      {formatRupiah(exp.amount)}
                    </td>
                    <td className="px-5 py-3.5 text-right whitespace-nowrap">
                      {!exp.is_cancelled && (
                        <Button
                          size="sm"
                          variant="ghost"
                          onClick={() => setCancelExpenseID(exp.id)}
                          className="text-rose-600 hover:text-rose-800 hover:bg-rose-50"
                          leftIcon={<Ban className="w-3.5 h-3.5" />}
                        >
                          Batalkan
                        </Button>
                      )}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}

        {/* Pagination Bar */}
        {!isLoading && totalPages > 1 && (
          <div className="px-5 py-3.5 border-t border-slate-100 bg-slate-50/50 flex items-center justify-between text-xs text-slate-500">
            <span>
              Menampilkan {expenses.length} dari total <span className="font-semibold">{total}</span> pengeluaran
            </span>
            <div className="flex items-center gap-1.5">
              <Button
                size="sm"
                variant="outline"
                disabled={page <= 1}
                onClick={() => setPage(p => Math.max(1, p - 1))}
                leftIcon={<ChevronLeft className="w-3.5 h-3.5" />}
              >
                Sebelumnya
              </Button>
              <span className="px-3 py-1 font-semibold text-slate-700">
                Hal {page} / {totalPages}
              </span>
              <Button
                size="sm"
                variant="outline"
                disabled={page >= totalPages}
                onClick={() => setPage(p => Math.min(totalPages, p + 1))}
                rightIcon={<ChevronRight className="w-3.5 h-3.5" />}
              >
                Berikutnya
              </Button>
            </div>
          </div>
        )}
      </Card>

      {/* Modal: Catat Pengeluaran Baru */}
      <Modal
        isOpen={isAddModalOpen}
        onClose={() => setIsAddModalOpen(false)}
        title="Catat Pengeluaran Operasional"
      >
        <form onSubmit={handleAddExpense} className="space-y-3.5">
          <Select
            label="Kategori Pengeluaran"
            required
            value={formExpense.category_id}
            onChange={e =>
              setFormExpense({ ...formExpense, category_id: parseInt(e.target.value, 10) })
            }
            options={categories.map(c => ({ value: c.id, label: c.name }))}
          />

          <div className="grid grid-cols-2 gap-3">
            <Input
              label="Nominal Pengeluaran (Rp)"
              required
              type="number"
              min="1000"
              value={formExpense.amount}
              onChange={e =>
                setFormExpense({ ...formExpense, amount: parseInt(e.target.value, 10) || 0 })
              }
            />
            <Input
              label="Tanggal Transaksi"
              required
              type="date"
              value={formExpense.expense_date}
              onChange={e => setFormExpense({ ...formExpense, expense_date: e.target.value })}
            />
          </div>

          <Input
            label="Keterangan Pengeluaran"
            required
            autoFocus
            placeholder="Contoh: Beli deterjen cair 20 liter"
            value={formExpense.description}
            onChange={e => setFormExpense({ ...formExpense, description: e.target.value })}
          />

          <Input
            label="Vendor / Tempat Pembelian (Opsional)"
            placeholder="Contoh: Toko Kimia Bersih / PLN"
            value={formExpense.vendor}
            onChange={e => setFormExpense({ ...formExpense, vendor: e.target.value })}
          />

          <div className="flex justify-end gap-2 pt-2">
            <Button type="button" variant="outline" onClick={() => setIsAddModalOpen(false)}>
              Batal
            </Button>
            <Button type="submit" variant="primary" isLoading={isSubmittingExpense}>
              Simpan Pengeluaran
            </Button>
          </div>
        </form>
      </Modal>

      {/* Modal: Kelola Kategori Pengeluaran */}
      <Modal
        isOpen={isCategoryModalOpen}
        onClose={() => setIsCategoryModalOpen(false)}
        title="Kategori Pengeluaran"
        maxWidth="lg"
      >
        <div className="space-y-5">
          {/* New Category Form */}
          <form onSubmit={handleCreateCategory} className="p-3.5 rounded-xl border border-slate-200 bg-slate-50 space-y-3">
            <h4 className="text-xs font-bold text-slate-800">Tambah Kategori Baru</h4>
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-2.5">
              <Input
                placeholder="Nama kategori..."
                value={newCatName}
                onChange={e => setNewCatName(e.target.value)}
                required
              />
              <Input
                placeholder="Keterangan opsional..."
                value={newCatDesc}
                onChange={e => setNewCatDesc(e.target.value)}
              />
            </div>
            <div className="flex justify-end">
              <Button type="submit" size="sm" variant="primary" isLoading={isSubmittingCat}>
                Tambah Kategori
              </Button>
            </div>
          </form>

          {/* Existing Categories List */}
          <div className="space-y-2">
            <h4 className="text-xs font-bold text-slate-800">Daftar Kategori Tersedia</h4>
            <div className="divide-y divide-slate-100 max-h-60 overflow-y-auto">
              {categories.map(c => (
                <div key={c.id} className="py-2.5 flex items-center justify-between text-xs">
                  <div>
                    <span className="font-semibold text-slate-900 block">{c.name}</span>
                    {c.description && <span className="text-slate-400 text-[11px]">{c.description}</span>}
                  </div>
                  <span className="text-[10px] text-emerald-700 bg-emerald-50 px-2 py-0.5 rounded font-medium">
                    Aktif
                  </span>
                </div>
              ))}
            </div>
          </div>
        </div>
      </Modal>

      {/* Modal: Batalkan Pengeluaran */}
      <Modal
        isOpen={cancelExpenseID !== null}
        onClose={() => setCancelExpenseID(null)}
        title="Batalkan Catatan Pengeluaran?"
      >
        <form onSubmit={handleCancelExpense} className="space-y-3.5">
          <Input
            label="Alasan Pembatalan"
            required
            autoFocus
            placeholder="Contoh: Salah input nominal / transaksi dibatalkan"
            value={cancelReason}
            onChange={e => setCancelReason(e.target.value)}
          />
          <div className="flex justify-end gap-2 pt-2">
            <Button type="button" variant="outline" onClick={() => setCancelExpenseID(null)}>
              Kembali
            </Button>
            <Button type="submit" variant="destructive" isLoading={isCancelling}>
              Ya, Batalkan Pengeluaran
            </Button>
          </div>
        </form>
      </Modal>
    </div>
  );
};
