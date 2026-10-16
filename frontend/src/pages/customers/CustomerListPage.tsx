import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { customersApi } from '@/api/customers';
import { Customer, Order } from '@/types/api';
import { formatRupiah } from '@/utils/currency';
import { formatJakartaDate } from '@/utils/date';
import { WORK_STATUS_META, PAYMENT_STATUS_META } from '@/utils/status';
import { Card } from '@/components/ui/Card';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';
import { Modal } from '@/components/ui/Modal';
import { Skeleton, EmptyState } from '@/components/ui/EmptyState';
import {
  Users,
  Search,
  Plus,
  Edit2,
  History,
  Phone,
  MapPin,
  ChevronLeft,
  ChevronRight,
} from 'lucide-react';

export const CustomerListPage: React.FC = () => {
  const navigate = useNavigate();

  const [customers, setCustomers] = useState<Customer[]>([]);
  const [total, setTotal] = useState<number>(0);
  const [page, setPage] = useState<number>(1);
  const [perPage] = useState<number>(15);
  const [totalPages, setTotalPages] = useState<number>(1);
  const [searchQuery, setSearchQuery] = useState<string>('');
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  // Add/Edit Modal
  const [isFormModalOpen, setIsFormModalOpen] = useState(false);
  const [editingCustomer, setEditingCustomer] = useState<Customer | null>(null);
  const [formData, setFormData] = useState({ name: '', phone: '', address: '', notes: '' });
  const [isSubmitting, setIsSubmitting] = useState(false);

  // History Drawer/Modal
  const [historyCustomer, setHistoryCustomer] = useState<Customer | null>(null);
  const [orderHistory, setOrderHistory] = useState<Order[]>([]);
  const [isLoadingHistory, setIsLoadingHistory] = useState(false);

  const fetchCustomers = async () => {
    setIsLoading(true);
    setError(null);
    try {
      const res = await customersApi.search(searchQuery.trim(), page, perPage);
      setCustomers(res.customers);
      if (res.meta) {
        setTotal(res.meta.total);
        setTotalPages(res.meta.total_pages);
      }
    } catch (err: any) {
      setError(err.message || 'Gagal memuat data pelanggan.');
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchCustomers();
  }, [page]);

  const handleSearchSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setPage(1);
    fetchCustomers();
  };

  const handleOpenAddModal = () => {
    setEditingCustomer(null);
    setFormData({ name: '', phone: '', address: '', notes: '' });
    setIsFormModalOpen(true);
  };

  const handleOpenEditModal = (c: Customer) => {
    setEditingCustomer(c);
    setFormData({
      name: c.name,
      phone: c.phone,
      address: c.address || '',
      notes: c.notes || '',
    });
    setIsFormModalOpen(true);
  };

  const handleSaveCustomer = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!formData.name.trim() || !formData.phone.trim()) return;

    setIsSubmitting(true);
    try {
      if (editingCustomer) {
        await customersApi.update(editingCustomer.id, formData);
      } else {
        await customersApi.create(formData);
      }
      setIsFormModalOpen(false);
      fetchCustomers();
    } catch (err: any) {
      alert(err.message || 'Gagal menyimpan data pelanggan.');
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleViewHistory = async (c: Customer) => {
    setHistoryCustomer(c);
    setIsLoadingHistory(true);
    try {
      const orders = await customersApi.getOrderHistory(c.id, 10);
      setOrderHistory(orders);
    } catch {
      setOrderHistory([]);
    } finally {
      setIsLoadingHistory(false);
    }
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-slate-900 tracking-tight">Data Pelanggan</h1>
          <p className="text-xs text-slate-500 mt-1">
            Kelola data kontak pelanggan dan riwayat transaksi cucian.
          </p>
        </div>
        <Button variant="primary" onClick={handleOpenAddModal} leftIcon={<Plus className="w-4 h-4" />}>
          Tambah Pelanggan
        </Button>
      </div>

      {/* Search Bar */}
      <Card className="p-4">
        <form onSubmit={handleSearchSubmit} className="flex gap-2">
          <Input
            placeholder="Cari berdasarkan nama atau no. telepon WhatsApp..."
            value={searchQuery}
            onChange={e => setSearchQuery(e.target.value)}
            leftIcon={<Search className="w-4 h-4" />}
          />
          <Button type="submit" variant="secondary" className="shrink-0">
            Cari
          </Button>
        </form>
      </Card>

      {/* Table Card */}
      <Card className="overflow-hidden p-0">
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
        ) : customers.length === 0 ? (
          <div className="p-8">
            <EmptyState
              icon={<Users className="w-6 h-6 text-slate-400" />}
              title="Belum Ada Pelanggan"
              description="Pelanggan akan otomatis tersimpan saat kasir membuat pesanan, atau tambahkan secara manual."
              action={
                <Button variant="primary" size="sm" onClick={handleOpenAddModal}>
                  Tambah Pelanggan Baru
                </Button>
              }
            />
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-left text-xs divide-y divide-slate-200">
              <thead className="bg-slate-50 text-[11px] font-bold text-slate-700 uppercase tracking-wider">
                <tr>
                  <th className="px-5 py-3.5">Nama Pelanggan</th>
                  <th className="px-5 py-3.5">No. Telepon / WA</th>
                  <th className="px-5 py-3.5">Alamat</th>
                  <th className="px-5 py-3.5">Catatan</th>
                  <th className="px-5 py-3.5 text-right">Aksi</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-100 bg-white">
                {customers.map(c => (
                  <tr key={c.id} className="hover:bg-slate-50/80 transition-colors">
                    <td className="px-5 py-3.5 font-bold text-slate-900">
                      {c.name}
                      <span className="block text-[10px] font-normal text-slate-400">
                        Terdaftar: {formatJakartaDate(c.created_at)}
                      </span>
                    </td>
                    <td className="px-5 py-3.5 font-medium text-slate-700">
                      <a
                        href={`https://wa.me/${c.phone.replace(/[^0-9]/g, '')}`}
                        target="_blank"
                        rel="noreferrer"
                        className="hover:text-teal-700 hover:underline flex items-center gap-1"
                      >
                        <Phone className="w-3 h-3 text-slate-400" />
                        {c.phone}
                      </a>
                    </td>
                    <td className="px-5 py-3.5 text-slate-600 max-w-xs truncate">
                      {c.address || '-'}
                    </td>
                    <td className="px-5 py-3.5 text-slate-500 italic max-w-xs truncate">
                      {c.notes ? `"${c.notes}"` : '-'}
                    </td>
                    <td className="px-5 py-3.5 text-right whitespace-nowrap space-x-1.5">
                      <Button
                        size="sm"
                        variant="ghost"
                        onClick={() => handleViewHistory(c)}
                        leftIcon={<History className="w-3.5 h-3.5 text-slate-500" />}
                      >
                        Riwayat
                      </Button>
                      <Button
                        size="sm"
                        variant="outline"
                        onClick={() => handleOpenEditModal(c)}
                        leftIcon={<Edit2 className="w-3.5 h-3.5" />}
                      >
                        Edit
                      </Button>
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
              Menampilkan {customers.length} dari total <span className="font-semibold">{total}</span> pelanggan
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

      {/* Modal: Tambah/Edit Pelanggan */}
      <Modal
        isOpen={isFormModalOpen}
        onClose={() => setIsFormModalOpen(false)}
        title={editingCustomer ? 'Edit Data Pelanggan' : 'Tambah Pelanggan Baru'}
      >
        <form onSubmit={handleSaveCustomer} className="space-y-3.5">
          <Input
            label="Nama Lengkap"
            required
            autoFocus
            value={formData.name}
            onChange={e => setFormData({ ...formData, name: e.target.value })}
            placeholder="Contoh: Pak Bambang"
          />
          <Input
            label="Nomor WhatsApp / HP"
            required
            type="tel"
            value={formData.phone}
            onChange={e => setFormData({ ...formData, phone: e.target.value })}
            placeholder="Contoh: 081234567890"
          />
          <Input
            label="Alamat Domisili"
            value={formData.address}
            onChange={e => setFormData({ ...formData, address: e.target.value })}
            placeholder="Contoh: Perumahan Indah Blok C No. 5"
          />
          <Input
            label="Catatan Khusus Pelanggan"
            value={formData.notes}
            onChange={e => setFormData({ ...formData, notes: e.target.value })}
            placeholder="Contoh: Selalu minta pewangi lavender / langganan kantor"
          />
          <div className="flex justify-end gap-2 pt-2">
            <Button type="button" variant="outline" onClick={() => setIsFormModalOpen(false)}>
              Batal
            </Button>
            <Button type="submit" variant="primary" isLoading={isSubmitting}>
              {editingCustomer ? 'Simpan Perubahan' : 'Tambah Pelanggan'}
            </Button>
          </div>
        </form>
      </Modal>

      {/* Modal: Riwayat Pesanan Pelanggan */}
      <Modal
        isOpen={historyCustomer !== null}
        onClose={() => setHistoryCustomer(null)}
        title={`Riwayat Transaksi: ${historyCustomer?.name || ''}`}
        maxWidth="lg"
      >
        {isLoadingHistory ? (
          <div className="space-y-3 p-4">
            <Skeleton className="h-10 w-full" />
            <Skeleton className="h-10 w-full" />
          </div>
        ) : orderHistory.length > 0 ? (
          <div className="divide-y divide-slate-100 max-h-96 overflow-y-auto">
            {orderHistory.map(ord => (
              <div
                key={ord.id}
                onClick={() => navigate(`/orders/${ord.id}`)}
                className="py-3 flex items-center justify-between text-xs hover:bg-slate-50 p-2 rounded-lg cursor-pointer"
              >
                <div>
                  <span className="font-bold text-slate-900 block">{ord.order_number}</span>
                  <span className="text-[11px] text-slate-400">
                    {formatJakartaDate(ord.order_date || ord.created_at)}
                  </span>
                </div>
                <div className="text-right">
                  <span className="font-bold text-slate-900 block">{formatRupiah(ord.total)}</span>
                  <span
                    className={`inline-block text-[10px] px-2 py-0.2 rounded border ${
                      PAYMENT_STATUS_META[ord.payment_status]?.badgeClass
                    }`}
                  >
                    {ord.payment_status}
                  </span>
                </div>
              </div>
            ))}
          </div>
        ) : (
          <p className="text-xs text-slate-400 text-center py-8">
            Belum ada riwayat pesanan untuk pelanggan ini.
          </p>
        )}
      </Modal>
    </div>
  );
};
