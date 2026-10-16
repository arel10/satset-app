import React, { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { ordersApi } from '@/api/orders';
import { paymentsApi } from '@/api/payments';
import { settingsApi } from '@/api/settings';
import { useAuth } from '@/context/AuthContext';
import { Order, WorkStatus, PaymentMethod } from '@/types/api';
import { formatRupiah, formatQuantity } from '@/utils/currency';
import { formatJakartaDate, formatJakartaDateTime, isOverdue } from '@/utils/date';
import { WORK_STATUS_META, PAYMENT_STATUS_META, getNextWorkStatuses } from '@/utils/status';
import { Button } from '@/components/ui/Button';
import { Card } from '@/components/ui/Card';
import { Modal } from '@/components/ui/Modal';
import { Input } from '@/components/ui/Input';
import { Skeleton } from '@/components/ui/EmptyState';
import { ThermalReceipt } from '@/components/receipt/ThermalReceipt';
import {
  ArrowLeft,
  Printer,
  CreditCard,
  Ban,
  CheckCircle2,
  Clock,
  User,
  Phone,
  MapPin,
  Calendar,
  AlertTriangle,
} from 'lucide-react';

export const OrderDetailPage: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { user } = useAuth();
  const isOwner = user?.role === 'owner';

  const [order, setOrder] = useState<Order | null>(null);
  const [businessSettings, setBusinessSettings] = useState<any>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Status Update State
  const [isUpdatingStatus, setIsUpdatingStatus] = useState(false);
  const [statusNotes, setStatusNotes] = useState('');

  // Payment Modal State
  const [isPaymentModalOpen, setIsPaymentModalOpen] = useState(false);
  const [paymentMethod, setPaymentMethod] = useState<PaymentMethod>('TUNAI');
  const [paymentNotes, setPaymentNotes] = useState('');
  const [isRecordingPayment, setIsRecordingPayment] = useState(false);

  // Cancel Order Modal State
  const [isCancelOrderModalOpen, setIsCancelOrderModalOpen] = useState(false);
  const [cancelOrderReason, setCancelOrderReason] = useState('');
  const [isCancellingOrder, setIsCancellingOrder] = useState(false);

  // Cancel Payment Modal State (Owner only)
  const [cancelPaymentID, setCancelPaymentID] = useState<number | null>(null);
  const [cancelPaymentReason, setCancelPaymentReason] = useState('');
  const [isCancellingPayment, setIsCancellingPayment] = useState(false);

  // Print Receipt Modal
  const [isReceiptModalOpen, setIsReceiptModalOpen] = useState(false);

  const fetchOrderDetail = async () => {
    if (!id) return;
    setIsLoading(true);
    setError(null);
    try {
      const [orderData, settings] = await Promise.all([
        ordersApi.getByID(parseInt(id, 10)),
        settingsApi.getAll(),
      ]);
      setOrder(orderData);
      setBusinessSettings(settings);
    } catch (err: any) {
      setError(err.message || 'Gagal memuat rincian pesanan.');
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchOrderDetail();
  }, [id]);

  // Handle status transition
  const handleUpdateStatus = async (nextStatus: WorkStatus) => {
    if (!order) return;
    if (!window.confirm(`Ubah status pesanan ke '${WORK_STATUS_META[nextStatus]?.label}'?`)) {
      return;
    }

    setIsUpdatingStatus(true);
    try {
      const updated = await ordersApi.updateStatus(order.id, nextStatus, statusNotes);
      setOrder(updated);
      setStatusNotes('');
    } catch (err: any) {
      alert(err.message || 'Gagal memperbarui status pesanan.');
    } finally {
      setIsUpdatingStatus(false);
    }
  };

  // Handle full payment
  const handleRecordPayment = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!order) return;

    setIsRecordingPayment(true);
    try {
      await paymentsApi.record(order.id, order.total, paymentMethod, paymentNotes);
      setIsPaymentModalOpen(false);
      setPaymentNotes('');
      await fetchOrderDetail();
    } catch (err: any) {
      alert(err.message || 'Gagal mencatat pembayaran.');
    } finally {
      setIsRecordingPayment(false);
    }
  };

  // Handle cancel order
  const handleCancelOrder = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!order || !cancelOrderReason.trim()) return;

    setIsCancellingOrder(true);
    try {
      const updated = await ordersApi.cancel(order.id, cancelOrderReason.trim());
      setOrder(updated);
      setIsCancelOrderModalOpen(false);
      setCancelOrderReason('');
    } catch (err: any) {
      alert(err.message || 'Gagal membatalkan pesanan.');
    } finally {
      setIsCancellingOrder(false);
    }
  };

  // Handle cancel payment
  const handleCancelPayment = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!cancelPaymentID || !cancelPaymentReason.trim()) return;

    setIsCancellingPayment(true);
    try {
      await paymentsApi.cancel(cancelPaymentID, cancelPaymentReason.trim());
      setCancelPaymentID(null);
      setCancelPaymentReason('');
      await fetchOrderDetail();
    } catch (err: any) {
      alert(err.message || 'Gagal membatalkan pembayaran.');
    } finally {
      setIsCancellingPayment(false);
    }
  };

  if (isLoading) {
    return (
      <div className="space-y-6">
        <Skeleton className="h-10 w-48" />
        <Skeleton className="h-44 w-full" />
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <Skeleton className="h-96 lg:col-span-2" />
          <Skeleton className="h-96" />
        </div>
      </div>
    );
  }

  if (error || !order) {
    return (
      <div className="p-8 text-center text-xs text-rose-700 bg-rose-50 rounded-xl">
        <p className="font-semibold text-sm mb-3">{error || 'Pesanan tidak ditemukan.'}</p>
        <Button variant="outline" size="sm" onClick={() => navigate('/orders')}>
          Kembali ke Daftar Pesanan
        </Button>
      </div>
    );
  }

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

  const nextStatuses = getNextWorkStatuses(order.work_status);
  const isTerminal = order.work_status === 'SELESAI' || order.work_status === 'DIBATALKAN';

  return (
    <div className="space-y-6">
      {/* Top Bar Navigation & Actions */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div className="flex items-center gap-3">
          <Button
            size="sm"
            variant="ghost"
            onClick={() => navigate('/orders')}
            leftIcon={<ArrowLeft className="w-4 h-4" />}
          >
            Daftar Pesanan
          </Button>
          <div className="h-4 w-px bg-slate-300" />
          <h1 className="text-xl font-bold text-slate-900">{order.order_number}</h1>
        </div>

        <div className="flex items-center gap-2">
          <Button
            variant="outline"
            size="sm"
            onClick={() => setIsReceiptModalOpen(true)}
            leftIcon={<Printer className="w-4 h-4" />}
          >
            Cetak Nota Struk
          </Button>

          {order.payment_status === 'BELUM_LUNAS' && order.work_status !== 'DIBATALKAN' && (
            <Button
              variant="primary"
              size="sm"
              onClick={() => setIsPaymentModalOpen(true)}
              leftIcon={<CreditCard className="w-4 h-4" />}
            >
              Catat Pelunasan
            </Button>
          )}

          {!isTerminal && (
            <Button
              variant="destructive"
              size="sm"
              onClick={() => setIsCancelOrderModalOpen(true)}
              leftIcon={<Ban className="w-4 h-4" />}
            >
              Batalkan Pesanan
            </Button>
          )}
        </div>
      </div>

      {/* Main Order Status Overview Card */}
      <Card className="p-5">
        <div className="grid grid-cols-2 sm:grid-cols-4 gap-4 pb-4 border-b border-slate-100 text-xs">
          <div>
            <span className="text-slate-400 block text-[11px]">Tanggal Masuk</span>
            <span className="font-semibold text-slate-800">
              {formatJakartaDateTime(order.order_date || order.created_at)}
            </span>
          </div>

          <div>
            <span className="text-slate-400 block text-[11px]">Estimasi Selesai</span>
            <span className={`font-semibold ${overdue ? 'text-rose-600' : 'text-slate-800'}`}>
              {formatJakartaDateTime(order.estimated_completion)}
            </span>
            {overdue && (
              <span className="inline-block text-[10px] bg-rose-50 text-rose-700 px-1.5 rounded border border-rose-200 ml-1">
                Terlambat
              </span>
            )}
          </div>

          <div>
            <span className="text-slate-400 block text-[11px]">Status Pengerjaan</span>
            <span className={`inline-flex items-center px-2 py-0.5 rounded-full text-xs border mt-0.5 ${workMeta.badgeClass}`}>
              {workMeta.label}
            </span>
          </div>

          <div>
            <span className="text-slate-400 block text-[11px]">Status Pembayaran</span>
            <span className={`inline-flex items-center px-2 py-0.5 rounded-full text-xs border mt-0.5 ${payMeta.badgeClass}`}>
              {payMeta.label}
            </span>
          </div>
        </div>

        {/* Workflow Transition Action Buttons */}
        {!isTerminal && nextStatuses.length > 0 && (
          <div className="pt-4 flex flex-col sm:flex-row sm:items-center justify-between gap-3">
            <span className="text-xs font-semibold text-slate-700">Langkah Pengerjaan Selanjutnya:</span>
            <div className="flex flex-wrap items-center gap-2">
              {nextStatuses.map(st => (
                <Button
                  key={st}
                  size="sm"
                  variant="primary"
                  isLoading={isUpdatingStatus}
                  onClick={() => handleUpdateStatus(st)}
                  leftIcon={<CheckCircle2 className="w-3.5 h-3.5" />}
                >
                  Lanjut ke: {WORK_STATUS_META[st]?.label}
                </Button>
              ))}
            </div>
          </div>
        )}
      </Card>

      {/* Grid: 2 Columns (Items & Customer Details) vs (Timeline & Payments) */}
      <div className="grid grid-cols-1 lg:grid-cols-12 gap-6 items-start">
        {/* Left Column: Items & Customer (7 Cols) */}
        <div className="lg:col-span-7 space-y-6">
          {/* Customer Card */}
          <Card title="Informasi Pelanggan">
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 text-xs">
              <div className="flex items-center gap-2.5">
                <div className="w-8 h-8 rounded-full bg-slate-100 flex items-center justify-center text-slate-500">
                  <User className="w-4 h-4" />
                </div>
                <div>
                  <span className="text-slate-400 block text-[11px]">Nama</span>
                  <span className="font-semibold text-slate-900">{order.customer?.name || 'Umum'}</span>
                </div>
              </div>

              <div className="flex items-center gap-2.5">
                <div className="w-8 h-8 rounded-full bg-slate-100 flex items-center justify-center text-slate-500">
                  <Phone className="w-4 h-4" />
                </div>
                <div>
                  <span className="text-slate-400 block text-[11px]">Telepon / WA</span>
                  <span className="font-semibold text-slate-900">{order.customer?.phone || '-'}</span>
                </div>
              </div>

              {order.customer?.address && (
                <div className="sm:col-span-2 flex items-start gap-2.5">
                  <div className="w-8 h-8 rounded-full bg-slate-100 flex items-center justify-center text-slate-500 shrink-0">
                    <MapPin className="w-4 h-4" />
                  </div>
                  <div>
                    <span className="text-slate-400 block text-[11px]">Alamat</span>
                    <span className="text-slate-700">{order.customer.address}</span>
                  </div>
                </div>
              )}
            </div>
          </Card>

          {/* Items Table */}
          <Card title="Rincian Item Laundry">
            <div className="overflow-x-auto">
              <table className="w-full text-left text-xs divide-y divide-slate-100">
                <thead className="text-[11px] font-bold text-slate-400 uppercase tracking-wider">
                  <tr>
                    <th className="pb-2">Layanan</th>
                    <th className="pb-2 text-right">Harga Satuan</th>
                    <th className="pb-2 text-center">Jumlah</th>
                    <th className="pb-2 text-right">Subtotal</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-slate-100 text-slate-700">
                  {order.items?.map(item => (
                    <tr key={item.id}>
                      <td className="py-2.5 font-semibold text-slate-900">{item.service_name}</td>
                      <td className="py-2.5 text-right text-slate-500">
                        {formatRupiah(item.price_per_unit)} / {item.service_unit}
                      </td>
                      <td className="py-2.5 text-center font-bold">
                        {formatQuantity(item.quantity, item.service_unit)}
                      </td>
                      <td className="py-2.5 text-right font-bold text-slate-900">
                        {formatRupiah(item.subtotal)}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>

            {/* Calculations Breakdown */}
            <div className="border-t border-slate-100 mt-4 pt-3 space-y-1.5 text-xs text-slate-600">
              <div className="flex justify-between">
                <span>Subtotal</span>
                <span>{formatRupiah(order.subtotal)}</span>
              </div>
              {order.discount > 0 && (
                <div className="flex justify-between text-rose-700">
                  <span>Diskon</span>
                  <span>-{formatRupiah(order.discount)}</span>
                </div>
              )}
              <div className="flex justify-between font-bold text-sm text-slate-900 pt-2 border-t border-slate-100">
                <span>Total Tagihan</span>
                <span className="text-teal-800 text-base">{formatRupiah(order.total)}</span>
              </div>
            </div>

            {order.notes && (
              <div className="mt-4 p-3 rounded-lg bg-amber-50/70 border border-amber-200/80 text-xs text-amber-800">
                <span className="font-semibold block mb-0.5">Catatan Pesanan:</span>
                <p>{order.notes}</p>
              </div>
            )}
          </Card>
        </div>

        {/* Right Column: Payments & Timeline History (5 Cols) */}
        <div className="lg:col-span-5 space-y-6">
          {/* Payments History Card */}
          <Card title="Riwayat Pembayaran">
            {order.payments && order.payments.length > 0 ? (
              <div className="space-y-3">
                {order.payments.map(p => (
                  <div
                    key={p.id}
                    className={`p-3 rounded-xl border text-xs ${
                      p.is_cancelled
                        ? 'bg-slate-50 border-slate-200 opacity-60'
                        : 'bg-emerald-50/50 border-emerald-200'
                    }`}
                  >
                    <div className="flex items-center justify-between font-bold">
                      <span className={p.is_cancelled ? 'line-through text-slate-500' : 'text-emerald-800'}>
                        {formatRupiah(p.amount)}
                      </span>
                      <span className="text-[11px] bg-white px-2 py-0.5 rounded border border-slate-200">
                        {p.payment_method}
                      </span>
                    </div>
                    <div className="mt-1 text-[11px] text-slate-500 flex justify-between">
                      <span>Diterima: {p.receiver?.full_name || 'Kasir'}</span>
                      <span>{formatJakartaDateTime(p.paid_at)}</span>
                    </div>

                    {p.is_cancelled && (
                      <div className="mt-2 pt-1 border-t border-slate-200 text-rose-700 text-[10px]">
                        Dibatalkan: {p.cancel_reason}
                      </div>
                    )}

                    {!p.is_cancelled && isOwner && (
                      <div className="mt-2 pt-2 border-t border-emerald-100 flex justify-end">
                        <button
                          type="button"
                          onClick={() => setCancelPaymentID(p.id)}
                          className="text-[10px] text-rose-600 hover:text-rose-800 font-medium underline"
                        >
                          Batalkan Pembayaran Ini
                        </button>
                      </div>
                    )}
                  </div>
                ))}
              </div>
            ) : (
              <p className="text-xs text-slate-400 text-center py-4">Belum ada pembayaran yang tercatat.</p>
            )}
          </Card>

          {/* Timeline Status History Card */}
          <Card title="Riwayat Alur Pengerjaan">
            {order.status_history && order.status_history.length > 0 ? (
              <div className="relative pl-4 space-y-4 before:absolute before:left-1.5 before:top-2 before:bottom-2 before:w-0.5 before:bg-slate-200">
                {order.status_history.map(hist => (
                  <div key={hist.id} className="relative text-xs">
                    <div className="absolute -left-4 top-1 w-2.5 h-2.5 rounded-full bg-teal-600 ring-4 ring-white" />
                    <div>
                      <span className="font-semibold text-slate-800">
                        {WORK_STATUS_META[hist.to_status]?.label || hist.to_status}
                      </span>
                      <span className="text-[10px] text-slate-400 ml-2">
                        {formatJakartaDateTime(hist.created_at)}
                      </span>
                    </div>
                    <p className="text-[11px] text-slate-500 mt-0.5">
                      Oleh: {hist.changer?.full_name || 'Petugas'}
                      {hist.notes && <span className="block italic text-slate-600 mt-0.5">"{hist.notes}"</span>}
                    </p>
                  </div>
                ))}
              </div>
            ) : (
              <p className="text-xs text-slate-400 text-center py-4">Belum ada riwayat status.</p>
            )}
          </Card>
        </div>
      </div>

      {/* Modal: Catat Pelunasan Pembayaran */}
      <Modal
        isOpen={isPaymentModalOpen}
        onClose={() => setIsPaymentModalOpen(false)}
        title="Catat Pembayaran Lunas"
        description="Sesuai kebijakan SatSet, pembayaran pesanan wajib lunas penuh."
      >
        <form onSubmit={handleRecordPayment} className="space-y-4">
          <div className="p-3.5 rounded-xl bg-slate-50 border border-slate-200 text-xs flex justify-between items-center">
            <span className="text-slate-600 font-medium">Tagihan Harus Dibayar:</span>
            <span className="text-base font-bold text-teal-800">{formatRupiah(order.total)}</span>
          </div>

          <div>
            <label className="block text-xs font-semibold text-slate-700 mb-1.5">
              Metode Pembayaran
            </label>
            <div className="grid grid-cols-3 gap-2">
              {(['TUNAI', 'TRANSFER', 'QRIS'] as const).map(m => (
                <button
                  key={m}
                  type="button"
                  onClick={() => setPaymentMethod(m)}
                  className={`py-2 text-xs font-semibold rounded-lg border transition-all ${
                    paymentMethod === m
                      ? 'border-teal-700 bg-teal-50 text-teal-800 ring-1 ring-teal-700'
                      : 'border-slate-200 text-slate-600 hover:bg-slate-50'
                  }`}
                >
                  {m}
                </button>
              ))}
            </div>
          </div>

          <Input
            label="Catatan Pembayaran (Opsional)"
            placeholder="Contoh: Transfer BCA an Budi / QRIS Statis Kasir"
            value={paymentNotes}
            onChange={e => setPaymentNotes(e.target.value)}
          />

          <div className="flex justify-end gap-2 pt-2">
            <Button type="button" variant="outline" onClick={() => setIsPaymentModalOpen(false)}>
              Batal
            </Button>
            <Button type="submit" variant="primary" isLoading={isRecordingPayment}>
              Konfirmasi Lunas ({formatRupiah(order.total)})
            </Button>
          </div>
        </form>
      </Modal>

      {/* Modal: Batalkan Pesanan */}
      <Modal
        isOpen={isCancelOrderModalOpen}
        onClose={() => setIsCancelOrderModalOpen(false)}
        title="Batalkan Pesanan Ini?"
        description="Pesanan yang dibatalkan tidak dapat dikembalikan ke alur pengerjaan normal."
      >
        <form onSubmit={handleCancelOrder} className="space-y-3.5">
          <Input
            label="Alasan Pembatalan"
            required
            autoFocus
            placeholder="Contoh: Pelanggan membatalkan permintaan / barang cacat"
            value={cancelOrderReason}
            onChange={e => setCancelOrderReason(e.target.value)}
          />
          <div className="flex justify-end gap-2 pt-2">
            <Button type="button" variant="outline" onClick={() => setIsCancelOrderModalOpen(false)}>
              Kembali
            </Button>
            <Button type="submit" variant="destructive" isLoading={isCancellingOrder}>
              Ya, Batalkan Pesanan
            </Button>
          </div>
        </form>
      </Modal>

      {/* Modal: Batalkan Pembayaran (Owner only) */}
      <Modal
        isOpen={cancelPaymentID !== null}
        onClose={() => setCancelPaymentID(null)}
        title="Koreksi / Batalkan Pembayaran"
        description="Pembatalan pembayaran akan mengembalikan status pesanan menjadi Belum Lunas."
      >
        <form onSubmit={handleCancelPayment} className="space-y-3.5">
          <Input
            label="Alasan Pembatalan Pembayaran"
            required
            autoFocus
            placeholder="Contoh: Salah input nominal / transaksi transfer tidak masuk"
            value={cancelPaymentReason}
            onChange={e => setCancelPaymentReason(e.target.value)}
          />
          <div className="flex justify-end gap-2 pt-2">
            <Button type="button" variant="outline" onClick={() => setCancelPaymentID(null)}>
              Kembali
            </Button>
            <Button type="submit" variant="destructive" isLoading={isCancellingPayment}>
              Konfirmasi Pembatalan
            </Button>
          </div>
        </form>
      </Modal>

      {/* Modal: Cetak Struk Nota Thermal */}
      <Modal
        isOpen={isReceiptModalOpen}
        onClose={() => setIsReceiptModalOpen(false)}
        title="Pratinjau Nota Kasir"
        maxWidth="md"
        footer={
          <div className="flex justify-end gap-2 w-full">
            <Button variant="outline" onClick={() => setIsReceiptModalOpen(false)}>
              Tutup
            </Button>
            <Button
              variant="primary"
              onClick={() => window.print()}
              leftIcon={<Printer className="w-4 h-4" />}
            >
              Cetak Sekarang
            </Button>
          </div>
        }
      >
        <ThermalReceipt order={order} businessSettings={businessSettings} />
      </Modal>
    </div>
  );
};
