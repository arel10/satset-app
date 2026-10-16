import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '@/context/AuthContext';
import { dashboardApi } from '@/api/dashboard';
import { OwnerDashboardData, KasirDashboardData } from '@/types/api';
import { formatRupiah } from '@/utils/currency';
import { Card } from '@/components/ui/Card';
import { Button } from '@/components/ui/Button';
import { Skeleton } from '@/components/ui/EmptyState';
import {
  TrendingUp,
  ShoppingBag,
  Clock,
  CheckCircle2,
  AlertTriangle,
  Receipt,
  Plus,
  Sparkles,
  ArrowRight,
  Wallet,
} from 'lucide-react';

export const DashboardPage: React.FC = () => {
  const { user, role } = useAuth();
  const navigate = useNavigate();

  const [ownerData, setOwnerData] = useState<OwnerDashboardData | null>(null);
  const [kasirData, setKasirData] = useState<KasirDashboardData | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchDashboard = async () => {
      setIsLoading(true);
      setError(null);
      try {
        if (role === 'owner') {
          const data = await dashboardApi.getOwnerDashboard();
          setOwnerData(data);
        } else if (role === 'kasir') {
          const data = await dashboardApi.getKasirDashboard();
          setKasirData(data);
        } else if (role === 'operator') {
          navigate('/orders', { replace: true });
        }
      } catch (err: any) {
        setError(err.message || 'Gagal memuat data dashboard.');
      } finally {
        setIsLoading(false);
      }
    };

    fetchDashboard();
  }, [role, navigate]);

  if (isLoading) {
    return (
      <div className="space-y-6">
        <div className="flex justify-between items-center">
          <Skeleton className="h-8 w-48" />
          <Skeleton className="h-10 w-36" />
        </div>
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
          <Skeleton className="h-28" />
          <Skeleton className="h-28" />
          <Skeleton className="h-28" />
          <Skeleton className="h-28" />
        </div>
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <Skeleton className="h-64" />
          <Skeleton className="h-64" />
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="p-6 rounded-xl bg-rose-50 border border-rose-200 text-rose-800 text-center">
        <p className="font-semibold text-sm mb-2">{error}</p>
        <Button size="sm" variant="outline" onClick={() => window.location.reload()}>
          Muat Ulang
        </Button>
      </div>
    );
  }

  // --- KASIR DASHBOARD VIEW ---
  if (role === 'kasir') {
    return (
      <div className="space-y-6">
        <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
          <div>
            <h1 className="text-2xl font-bold text-slate-900 tracking-tight">Dashboard Kasir</h1>
            <p className="text-xs text-slate-500 mt-1">
              Selamat bertugas, <span className="font-semibold text-slate-700">{user?.full_name}</span>.
            </p>
          </div>
          <Button
            variant="primary"
            onClick={() => navigate('/orders/new')}
            leftIcon={<Plus className="w-4 h-4" />}
          >
            Buat Pesanan Baru (POS)
          </Button>
        </div>

        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
          <Card className="p-5">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-xs font-medium text-slate-500">Pesanan Hari Ini</p>
                <p className="text-2xl font-bold text-slate-900 mt-1">
                  {kasirData?.today_orders || 0}
                </p>
              </div>
              <div className="w-10 h-10 rounded-xl bg-teal-50 text-teal-700 flex items-center justify-center">
                <ShoppingBag className="w-5 h-5" />
              </div>
            </div>
          </Card>

          <Card className="p-5">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-xs font-medium text-slate-500">Siap Diambil</p>
                <p className="text-2xl font-bold text-teal-700 mt-1">
                  {kasirData?.ready_for_pickup || 0}
                </p>
              </div>
              <div className="w-10 h-10 rounded-xl bg-teal-50 text-teal-700 flex items-center justify-center">
                <CheckCircle2 className="w-5 h-5" />
              </div>
            </div>
          </Card>

          <Card className="p-5">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-xs font-medium text-slate-500">Belum Lunas</p>
                <p className="text-2xl font-bold text-amber-700 mt-1">
                  {kasirData?.unpaid_orders || 0}
                </p>
              </div>
              <div className="w-10 h-10 rounded-xl bg-amber-50 text-amber-700 flex items-center justify-center">
                <Clock className="w-5 h-5" />
              </div>
            </div>
          </Card>

          <Card className="p-5">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-xs font-medium text-slate-500">Kas Diterima Hari Ini</p>
                <p className="text-xl font-bold text-slate-900 mt-1">
                  {formatRupiah(kasirData?.today_cash_received || 0)}
                </p>
              </div>
              <div className="w-10 h-10 rounded-xl bg-blue-50 text-blue-700 flex items-center justify-center">
                <Wallet className="w-5 h-5" />
              </div>
            </div>
          </Card>
        </div>

        <div className="flex justify-end">
          <Button
            variant="outline"
            onClick={() => navigate('/orders')}
            rightIcon={<ArrowRight className="w-4 h-4" />}
          >
            Buka Daftar Pesanan
          </Button>
        </div>
      </div>
    );
  }

  // --- OWNER DASHBOARD VIEW ---
  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-slate-900 tracking-tight">Ringkasan Bisnis</h1>
          <p className="text-xs text-slate-500 mt-1">
            Pantau performa harian dan metrik operasional laundry secara real-time.
          </p>
        </div>
        <div className="flex items-center gap-2.5">
          <Button
            variant="outline"
            onClick={() => navigate('/ai')}
            leftIcon={<Sparkles className="w-4 h-4 text-teal-600" />}
          >
            Tanya AI
          </Button>
          <Button
            variant="primary"
            onClick={() => navigate('/orders/new')}
            leftIcon={<Plus className="w-4 h-4" />}
          >
            Buat Pesanan
          </Button>
        </div>
      </div>

      {/* Periode Info */}
      <div className="flex items-center gap-2 text-xs text-slate-500 bg-teal-50/50 border border-teal-100/80 px-3.5 py-2 rounded-xl">
        <Clock className="w-4 h-4 text-teal-600" />
        <span>Periode Data: <strong className="text-teal-900">{ownerData?.period || 'Bulan Berjalan'}</strong></span>
      </div>

      {/* KPI Cards Row 1: Finansial */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <Card className="p-5">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-xs font-medium text-slate-500">Total Omzet</p>
              <p className="text-2xl font-bold text-slate-900 mt-1">
                {formatRupiah(ownerData?.revenue?.total_revenue || 0)}
              </p>
              <p className="text-[11px] text-slate-400 mt-0.5">
                {ownerData?.revenue?.total_orders || 0} pesanan selesai/berjalan
              </p>
            </div>
            <div className="w-10 h-10 rounded-xl bg-teal-50 text-teal-700 flex items-center justify-center">
              <TrendingUp className="w-5 h-5" />
            </div>
          </div>
        </Card>

        <Card className="p-5">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-xs font-medium text-slate-500">Kas Masuk Riil</p>
              <p className="text-2xl font-bold text-emerald-700 mt-1">
                {formatRupiah(ownerData?.cash_received || 0)}
              </p>
              <p className="text-[11px] text-slate-400 mt-0.5">Pembayaran lunas diterima</p>
            </div>
            <div className="w-10 h-10 rounded-xl bg-emerald-50 text-emerald-700 flex items-center justify-center">
              <Wallet className="w-5 h-5" />
            </div>
          </div>
        </Card>

        <Card className="p-5">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-xs font-medium text-slate-500">Total Pengeluaran</p>
              <p className="text-2xl font-bold text-slate-900 mt-1">
                {formatRupiah(ownerData?.total_expenses || 0)}
              </p>
              <p className="text-[11px] text-slate-400 mt-0.5">Biaya operasional</p>
            </div>
            <div className="w-10 h-10 rounded-xl bg-slate-100 text-slate-700 flex items-center justify-center">
              <Receipt className="w-5 h-5" />
            </div>
          </div>
        </Card>

        <Card className="p-5">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-xs font-medium text-slate-500">Estimasi Laba Bersih</p>
              <p className={`text-2xl font-bold mt-1 ${(ownerData?.estimated_profit || 0) >= 0 ? 'text-teal-700' : 'text-rose-700'}`}>
                {formatRupiah(ownerData?.estimated_profit || 0)}
              </p>
              <p className="text-[11px] text-slate-400 mt-0.5">Kas Masuk - Pengeluaran</p>
            </div>
            <div className="w-10 h-10 rounded-xl bg-blue-50 text-blue-700 flex items-center justify-center">
              <TrendingUp className="w-5 h-5" />
            </div>
          </div>
        </Card>
      </div>

      {/* KPI Cards Row 2: Operasional */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <Card className="p-5">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-xs font-medium text-slate-500">Pesanan Aktif</p>
              <p className="text-2xl font-bold text-slate-900 mt-1">
                {ownerData?.active_orders || 0}
              </p>
              <p className="text-[11px] text-slate-400 mt-0.5">Sedang dalam proses</p>
            </div>
            <div className="w-10 h-10 rounded-xl bg-slate-100 text-slate-700 flex items-center justify-center">
              <ShoppingBag className="w-5 h-5" />
            </div>
          </div>
        </Card>

        <Card className="p-5">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-xs font-medium text-slate-500">Siap Diambil</p>
              <p className="text-2xl font-bold text-teal-700 mt-1">
                {ownerData?.ready_for_pickup || 0}
              </p>
              <p className="text-[11px] text-slate-400 mt-0.5">Menunggu diambil pelanggan</p>
            </div>
            <div className="w-10 h-10 rounded-xl bg-teal-50 text-teal-700 flex items-center justify-center">
              <CheckCircle2 className="w-5 h-5" />
            </div>
          </div>
        </Card>

        <Card className="p-5">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-xs font-medium text-slate-500">Pesanan Terlambat</p>
              <p className={`text-2xl font-bold mt-1 ${ownerData?.overdue_orders ? 'text-rose-700' : 'text-slate-900'}`}>
                {ownerData?.overdue_orders || 0}
              </p>
              <p className="text-[11px] text-slate-400 mt-0.5">Melewati estimasi selesai</p>
            </div>
            <div className="w-10 h-10 rounded-xl bg-rose-50 text-rose-700 flex items-center justify-center">
              <AlertTriangle className="w-5 h-5" />
            </div>
          </div>
        </Card>

        <Card className="p-5">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-xs font-medium text-slate-500">Belum Lunas</p>
              <p className="text-2xl font-bold text-amber-700 mt-1">
                {ownerData?.unpaid_orders || 0}
              </p>
              <p className="text-[11px] text-slate-400 mt-0.5">Tagihan pesanan aktif</p>
            </div>
            <div className="w-10 h-10 rounded-xl bg-amber-50 text-amber-700 flex items-center justify-center">
              <Clock className="w-5 h-5" />
            </div>
          </div>
        </Card>
      </div>

      {/* Breakdown Section: Tren Pendapatan & Status Pengerjaan */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <Card
          title="Tren Pendapatan Harian"
          subtitle="Grafik aktivitas order dan omzet harian periode ini"
          action={
            <Button size="sm" variant="ghost" onClick={() => navigate('/reports')}>
              Laporan Lengkap
            </Button>
          }
        >
          <div className="space-y-3 pt-1">
            {ownerData?.daily_revenue && ownerData.daily_revenue.length > 0 ? (
              ownerData.daily_revenue.map(trend => (
                <div key={trend.date} className="flex items-center justify-between text-xs py-1.5 border-b border-slate-100 last:border-0">
                  <span className="font-medium text-slate-700">{trend.date}</span>
                  <div className="flex items-center gap-4">
                    <span className="text-slate-500">{trend.order_count} pesanan</span>
                    <span className="font-semibold text-slate-900 min-w-[100px] text-right">
                      {formatRupiah(trend.revenue)}
                    </span>
                  </div>
                </div>
              ))
            ) : (
              <p className="text-xs text-slate-400 text-center py-6">Belum ada transaksi pada periode ini.</p>
            )}
          </div>
        </Card>

        <Card
          title="Distribusi Status Pengerjaan"
          subtitle="Jumlah cucian berdasarkan tahapan pengerjaan saat ini"
          action={
            <Button size="sm" variant="ghost" onClick={() => navigate('/orders')}>
              Kelola Pesanan
            </Button>
          }
        >
          <div className="grid grid-cols-2 gap-3 pt-1">
            {ownerData?.orders_by_status && ownerData.orders_by_status.length > 0 ? (
              ownerData.orders_by_status.map(item => {
                const statusName = item.work_status || item.status || 'Lainnya';
                return (
                  <div
                    key={statusName}
                    className="p-3 rounded-lg border border-slate-100 bg-slate-50 flex items-center justify-between"
                  >
                    <span className="text-xs text-slate-600 font-medium capitalize">
                      {statusName.toLowerCase().replace('_', ' ')}
                    </span>
                    <span className="text-sm font-bold text-slate-900">{item.count}</span>
                  </div>
                );
              })
            ) : (
              <p className="text-xs text-slate-400 text-center py-6 col-span-2">Belum ada cucian aktif.</p>
            )}
          </div>
        </Card>
      </div>
    </div>
  );
};
