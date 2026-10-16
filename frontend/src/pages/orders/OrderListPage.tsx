import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { ordersApi, OrderFilterParams } from '@/api/orders';
import { Order, WorkStatus, PaymentStatus } from '@/types/api';
import { formatRupiah } from '@/utils/currency';
import { formatJakartaDate, formatJakartaTimeOnly, isOverdue } from '@/utils/date';
import { WORK_STATUS_META, PAYMENT_STATUS_META } from '@/utils/status';
import { Card } from '@/components/ui/Card';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';
import { Select } from '@/components/ui/Select';
import { Skeleton, EmptyState } from '@/components/ui/EmptyState';
import {
  Search,
  Filter,
  Plus,
  ShoppingBag,
  Clock,
  Eye,
  AlertTriangle,
  ChevronLeft,
  ChevronRight,
} from 'lucide-react';

const STATUS_TABS: Array<{ value: string; label: string }> = [
  { value: '', label: 'Semua Status' },
  { value: 'DITERIMA', label: 'Diterima' },
  { value: 'DICUCI', label: 'Dicuci' },
  { value: 'DIKERINGKAN', label: 'Dikeringkan' },
  { value: 'DISETRIKA', label: 'Disetrika' },
  { value: 'DIKEMAS', label: 'Dikemas' },
  { value: 'SIAP_DIAMBIL', label: 'Siap Diambil' },
  { value: 'SELESAI', label: 'Selesai' },
  { value: 'DIBATALKAN', label: 'Dibatalkan' },
];

export const OrderListPage: React.FC = () => {
  const navigate = useNavigate();

  const [orders, setOrders] = useState<Order[]>([]);
  const [total, setTotal] = useState<number>(0);
  const [page, setPage] = useState<number>(1);
  const [perPage] = useState<number>(15);
  const [totalPages, setTotalPages] = useState<number>(1);

  // Filters
  const [selectedStatus, setSelectedStatus] = useState<string>('');
  const [selectedPaymentStatus, setSelectedPaymentStatus] = useState<string>('');
  const [searchQuery, setSearchQuery] = useState<string>('');
  const [isOverdueOnly, setIsOverdueOnly] = useState<boolean>(false);
  const [startDate, setStartDate] = useState<string>('');
  const [endDate, setEndDate] = useState<string>('');

  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  const fetchOrders = async () => {
    setIsLoading(true);
    setError(null);
    try {
      const params: OrderFilterParams = {
        page,
        per_page: perPage,
        status: selectedStatus || undefined,
        payment_status: selectedPaymentStatus || undefined,
        q: searchQuery.trim() || undefined,
        overdue: isOverdueOnly || undefined,
        start_date: startDate || undefined,
        end_date: endDate || undefined,
      };

      const res = await ordersApi.list(params);
      setOrders(res.orders);
      if (res.meta) {
        setTotal(res.meta.total);
        setTotalPages(res.meta.total_pages);
      }
    } catch (err: any) {
      setError(err.message || 'Gagal memuat daftar pesanan.');
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchOrders();
  }, [page, selectedStatus, selectedPaymentStatus, isOverdueOnly, startDate, endDate]);

  const handleSearchSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setPage(1);
    fetchOrders();
  };

  const handleResetFilters = () => {
    setSelectedStatus('');
    setSelectedPaymentStatus('');
    setSearchQuery('');
    setIsOverdueOnly(false);
    setStartDate('');
    setEndDate('');
    setPage(1);
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-slate-900 tracking-tight">Daftar Pesanan Laundry</h1>
          <p className="text-xs text-slate-500 mt-1">
            Kelola tahapan pengerjaan cucian, pembayaran, dan cetak nota.
          </p>
        </div>
        <Button
          variant="primary"
          onClick={() => navigate('/orders/new')}
          leftIcon={<Plus className="w-4 h-4" />}
        >
          Buat Pesanan Baru
        </Button>
      </div>

      {/* Status Tabs Navigation */}
      <div className="flex items-center gap-1.5 overflow-x-auto pb-1 border-b border-slate-200">
        {STATUS_TABS.map(tab => (
          <button
            key={tab.value}
            onClick={() => {
              setSelectedStatus(tab.value);
              setPage(1);
            }}
            className={`px-3 py-2 text-xs font-medium whitespace-nowrap border-b-2 transition-all -mb-[1px] ${
              selectedStatus === tab.value
                ? 'border-teal-700 text-teal-800 font-semibold'
                : 'border-transparent text-slate-500 hover:text-slate-800 hover:border-slate-300'
            }`}
          >
            {tab.label}
          </button>
        ))}
      </div>

      {/* Filter Bar */}
      <Card className="p-4">
        <form onSubmit={handleSearchSubmit} className="space-y-3">
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-3">
            <Input
              placeholder="Cari no. nota / nama / HP..."
              value={searchQuery}
              onChange={e => setSearchQuery(e.target.value)}
              leftIcon={<Search className="w-4 h-4" />}
            />

            <Select
              value={selectedPaymentStatus}
              onChange={e => {
                setSelectedPaymentStatus(e.target.value);
                setPage(1);
              }}
              options={[
                { value: '', label: 'Semua Status Pembayaran' },
                { value: 'LUNAS', label: 'Lunas' },
                { value: 'BELUM_LUNAS', label: 'Belum Lunas' },
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

          <div className="flex items-center justify-between pt-1 border-t border-slate-100 text-xs">
            <label className="flex items-center gap-2 cursor-pointer select-none text-slate-700 font-medium">
              <input
                type="checkbox"
                checked={isOverdueOnly}
                onChange={e => {
                  setIsOverdueOnly(e.target.checked);
                  setPage(1);
                }}
                className="w-4 h-4 rounded text-teal-700 focus:ring-teal-600 border-slate-300"
              />
              <span className="flex items-center gap-1 text-rose-700 font-semibold">
                <AlertTriangle className="w-3.5 h-3.5" />
                Hanya Pesanan Terlambat (Overdue)
              </span>
            </label>

            <div className="flex items-center gap-2">
              <Button type="button" size="sm" variant="ghost" onClick={handleResetFilters}>
                Reset Filter
              </Button>
              <Button type="submit" size="sm" variant="secondary">
                Terapkan
              </Button>
            </div>
          </div>
        </form>
      </Card>

      {/* Orders Table Card */}
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
        ) : orders.length === 0 ? (
          <div className="p-8">
            <EmptyState
              icon={<ShoppingBag className="w-6 h-6 text-slate-400" />}
              title="Tidak Ada Pesanan Ditemukan"
              description="Belum ada transaksi pesanan yang sesuai dengan filter pencarian ini."
              action={
                <Button variant="primary" size="sm" onClick={() => navigate('/orders/new')}>
                  Buat Pesanan Sekarang
                </Button>
              }
            />
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-left text-xs text-slate-600 divide-y divide-slate-200">
              <thead className="bg-slate-50 text-[11px] font-bold text-slate-700 uppercase tracking-wider">
                <tr>
                  <th className="px-5 py-3.5">No. Nota</th>
                  <th className="px-5 py-3.5">Pelanggan</th>
                  <th className="px-5 py-3.5">Estimasi Selesai</th>
                  <th className="px-5 py-3.5">Status Pengerjaan</th>
                  <th className="px-5 py-3.5">Pembayaran</th>
                  <th className="px-5 py-3.5 text-right">Total Tagihan</th>
                  <th className="px-5 py-3.5 text-center">Aksi</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-100 bg-white">
                {orders.map(order => {
                  const workMeta = WORK_STATUS_META[order.work_status] || {
                    label: order.work_status,
                    badgeClass: 'bg-slate-100 text-slate-700',
                  };
                  const payMeta = PAYMENT_STATUS_META[order.payment_status] || {
                    label: order.payment_status,
                    badgeClass: 'bg-slate-100 text-slate-700',
                  };
                  const overdue =
                    order.work_status !== 'SELESAI' &&
                    order.work_status !== 'DIBATALKAN' &&
                    isOverdue(order.estimated_completion);

                  return (
                    <tr
                      key={order.id}
                      onClick={() => navigate(`/orders/${order.id}`)}
                      className="hover:bg-slate-50/80 cursor-pointer transition-colors"
                    >
                      <td className="px-5 py-3.5 font-bold text-slate-900 whitespace-nowrap">
                        {order.order_number}
                        <span className="block text-[10px] font-normal text-slate-400">
                          {formatJakartaDate(order.order_date || order.created_at)}
                        </span>
                      </td>

                      <td className="px-5 py-3.5">
                        <span className="font-semibold text-slate-800 block">
                          {order.customer?.name || 'Umum'}
                        </span>
                        <span className="text-[11px] text-slate-400">
                          {order.customer?.phone || '-'}
                        </span>
                      </td>

                      <td className="px-5 py-3.5 whitespace-nowrap">
                        <span
                          className={`font-medium block ${
                            overdue ? 'text-rose-600 font-semibold' : 'text-slate-700'
                          }`}
                        >
                          {formatJakartaDate(order.estimated_completion)}
                        </span>
                        <span className="text-[10px] text-slate-400">
                          {formatJakartaTimeOnly(order.estimated_completion)}
                        </span>
                        {overdue && (
                          <span className="inline-block text-[9px] bg-rose-50 text-rose-700 px-1.5 py-0.2 rounded border border-rose-200 mt-0.5">
                            Terlambat
                          </span>
                        )}
                      </td>

                      <td className="px-5 py-3.5 whitespace-nowrap">
                        <span
                          className={`inline-flex items-center px-2.5 py-1 rounded-full text-[11px] border ${workMeta.badgeClass}`}
                        >
                          {workMeta.label}
                        </span>
                      </td>

                      <td className="px-5 py-3.5 whitespace-nowrap">
                        <span
                          className={`inline-flex items-center px-2.5 py-1 rounded-full text-[11px] border ${payMeta.badgeClass}`}
                        >
                          {payMeta.label}
                        </span>
                      </td>

                      <td className="px-5 py-3.5 text-right font-bold text-slate-900 whitespace-nowrap">
                        {formatRupiah(order.total)}
                      </td>

                      <td className="px-5 py-3.5 text-center whitespace-nowrap" onClick={e => e.stopPropagation()}>
                        <Button
                          size="sm"
                          variant="ghost"
                          onClick={() => navigate(`/orders/${order.id}`)}
                          leftIcon={<Eye className="w-3.5 h-3.5" />}
                        >
                          Detail
                        </Button>
                      </td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </div>
        )}

        {/* Pagination Bar */}
        {!isLoading && totalPages > 1 && (
          <div className="px-5 py-3.5 border-t border-slate-100 bg-slate-50/50 flex items-center justify-between text-xs text-slate-500">
            <span>
              Menampilkan {orders.length} dari total <span className="font-semibold">{total}</span> pesanan
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
    </div>
  );
};
