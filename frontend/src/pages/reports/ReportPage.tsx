import React, { useEffect, useState } from 'react';
import { reportsApi } from '@/api/reports';
import { BusinessReportData } from '@/types/api';
import { formatRupiah } from '@/utils/currency';
import { formatJakartaDate } from '@/utils/date';
import { Card } from '@/components/ui/Card';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';
import { Skeleton } from '@/components/ui/EmptyState';
import {
  TrendingUp,
  Wallet,
  Receipt,
  Printer,
  Calendar,
  Layers,
  PieChart,
} from 'lucide-react';

export const ReportPage: React.FC = () => {
  const [report, setReport] = useState<BusinessReportData | null>(null);
  const [startDate, setStartDate] = useState('');
  const [endDate, setEndDate] = useState('');
  const [activePreset, setActivePreset] = useState<'today' | '7d' | 'month' | 'custom'>('month');

  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const setPresetDates = (preset: 'today' | '7d' | 'month') => {
    setActivePreset(preset);
    const now = new Date();
    const todayStr = now.toISOString().split('T')[0];

    if (preset === 'today') {
      setStartDate(todayStr);
      setEndDate(todayStr);
    } else if (preset === '7d') {
      const past = new Date();
      past.setDate(past.getDate() - 7);
      setStartDate(past.toISOString().split('T')[0]);
      setEndDate(todayStr);
    } else if (preset === 'month') {
      const firstDay = new Date(now.getFullYear(), now.getMonth(), 1);
      setStartDate(firstDay.toISOString().split('T')[0]);
      setEndDate(todayStr);
    }
  };

  useEffect(() => {
    // Default to this month
    setPresetDates('month');
  }, []);

  const fetchReport = async () => {
    setIsLoading(true);
    setError(null);
    try {
      const data = await reportsApi.getBusinessReport(startDate || undefined, endDate || undefined);
      setReport(data);
    } catch (err: any) {
      setError(err.message || 'Gagal memuat laporan bisnis.');
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    if (startDate && endDate) {
      fetchReport();
    }
  }, [startDate, endDate]);

  const handleCustomFilterSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setActivePreset('custom');
    fetchReport();
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-slate-900 tracking-tight">Laporan Keuangan &amp; Operasional</h1>
          <p className="text-xs text-slate-500 mt-1">
            Analisis menyeluruh omzet, arus kas, pengeluaran, dan laba operasional bersih.
          </p>
        </div>
        <Button
          variant="outline"
          size="sm"
          onClick={() => window.print()}
          leftIcon={<Printer className="w-4 h-4" />}
        >
          Cetak Laporan
        </Button>
      </div>

      {/* Date Filter & Presets */}
      <Card className="p-4">
        <div className="flex flex-col lg:flex-row lg:items-center justify-between gap-4">
          <div className="flex items-center gap-2 overflow-x-auto pb-1">
            <button
              onClick={() => setPresetDates('today')}
              className={`px-3 py-1.5 rounded-lg text-xs font-semibold transition-colors shrink-0 ${
                activePreset === 'today'
                  ? 'bg-teal-700 text-white'
                  : 'bg-slate-100 text-slate-600 hover:bg-slate-200'
              }`}
            >
              Hari Ini
            </button>
            <button
              onClick={() => setPresetDates('7d')}
              className={`px-3 py-1.5 rounded-lg text-xs font-semibold transition-colors shrink-0 ${
                activePreset === '7d'
                  ? 'bg-teal-700 text-white'
                  : 'bg-slate-100 text-slate-600 hover:bg-slate-200'
              }`}
            >
              7 Hari Terakhir
            </button>
            <button
              onClick={() => setPresetDates('month')}
              className={`px-3 py-1.5 rounded-lg text-xs font-semibold transition-colors shrink-0 ${
                activePreset === 'month'
                  ? 'bg-teal-700 text-white'
                  : 'bg-slate-100 text-slate-600 hover:bg-slate-200'
              }`}
            >
              Bulan Ini
            </button>
          </div>

          <form onSubmit={handleCustomFilterSubmit} className="flex items-center gap-2.5">
            <Input
              type="date"
              value={startDate}
              onChange={e => {
                setStartDate(e.target.value);
                setActivePreset('custom');
              }}
              className="text-xs py-1.5"
            />
            <span className="text-slate-400 text-xs">s/d</span>
            <Input
              type="date"
              value={endDate}
              onChange={e => {
                setEndDate(e.target.value);
                setActivePreset('custom');
              }}
              className="text-xs py-1.5"
            />
            <Button type="submit" size="sm" variant="secondary" className="shrink-0">
              Terapkan
            </Button>
          </form>
        </div>
      </Card>

      {isLoading ? (
        <div className="space-y-6">
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
            {Array.from({ length: 4 }).map((_, i) => (
              <Skeleton key={i} className="h-28" />
            ))}
          </div>
          <Skeleton className="h-80" />
        </div>
      ) : error || !report ? (
        <div className="p-8 text-center text-xs text-rose-700 bg-rose-50 rounded-xl">
          {error || 'Gagal memuat data laporan.'}
        </div>
      ) : (
        <>
          {/* Summary Financial Metric Cards */}
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
            <Card className="p-5">
              <span className="text-xs font-medium text-slate-500 block">Omzet Bersih (Penjualan)</span>
              <p className="text-2xl font-bold text-slate-900 mt-1">
                {formatRupiah(report.net_revenue)}
              </p>
              <p className="text-[11px] text-slate-400 mt-0.5">
                Total {report.total_orders} pesanan
              </p>
            </Card>

            <Card className="p-5">
              <span className="text-xs font-medium text-slate-500 block">Kas Masuk (Uang Riil)</span>
              <p className="text-2xl font-bold text-emerald-700 mt-1">
                {formatRupiah(report.total_received)}
              </p>
              <p className="text-[11px] text-slate-400 mt-0.5">
                Piutang: {formatRupiah(report.unpaid_amount)}
              </p>
            </Card>

            <Card className="p-5">
              <span className="text-xs font-medium text-slate-500 block">Total Pengeluaran</span>
              <p className="text-2xl font-bold text-slate-900 mt-1">
                {formatRupiah(report.total_expenses)}
              </p>
              <p className="text-[11px] text-slate-400 mt-0.5">Biaya operasional</p>
            </Card>

            <Card className="p-5 bg-teal-50/50 border-teal-200">
              <span className="text-xs font-semibold text-teal-800 block">Laba Operasional Bersih</span>
              <p className={`text-2xl font-bold mt-1 ${report.net_profit >= 0 ? 'text-teal-900' : 'text-rose-700'}`}>
                {formatRupiah(report.net_profit)}
              </p>
              <p className="text-[11px] text-teal-700 mt-0.5">Omzet bersih - Biaya</p>
            </Card>
          </div>

          {/* Breakdown Section: Pendapatan per Layanan & Pengeluaran per Kategori */}
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            {/* Revenue By Service */}
            <Card title="Pendapatan per Layanan Laundry" subtitle="Kontribusi omzet dari masing-masing jenis layanan">
              <div className="overflow-x-auto pt-1">
                <table className="w-full text-left text-xs divide-y divide-slate-100">
                  <thead className="text-[11px] font-bold text-slate-400 uppercase">
                    <tr>
                      <th className="pb-2">Layanan</th>
                      <th className="pb-2 text-center">Volume</th>
                      <th className="pb-2 text-right">Pendapatan</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-slate-100 text-slate-700">
                    {report.revenue_by_service && report.revenue_by_service.length > 0 ? (
                      report.revenue_by_service.map(s => (
                        <tr key={s.service_id}>
                          <td className="py-2.5 font-semibold text-slate-900">{s.service_name}</td>
                          <td className="py-2.5 text-center text-slate-500">{s.quantity_total}</td>
                          <td className="py-2.5 text-right font-bold text-slate-900">
                            {formatRupiah(s.revenue_total)}
                          </td>
                        </tr>
                      ))
                    ) : (
                      <tr>
                        <td colSpan={3} className="py-6 text-center text-slate-400">
                          Belum ada data layanan pada periode ini.
                        </td>
                      </tr>
                    )}
                  </tbody>
                </table>
              </div>
            </Card>

            {/* Expenses By Category */}
            <Card title="Pengeluaran per Kategori" subtitle="Distribusi pengeluaran biaya usaha laundry">
              <div className="overflow-x-auto pt-1">
                <table className="w-full text-left text-xs divide-y divide-slate-100">
                  <thead className="text-[11px] font-bold text-slate-400 uppercase">
                    <tr>
                      <th className="pb-2">Kategori</th>
                      <th className="pb-2 text-right">Total Biaya</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-slate-100 text-slate-700">
                    {report.expenses_by_category && report.expenses_by_category.length > 0 ? (
                      report.expenses_by_category.map(c => (
                        <tr key={c.category_id}>
                          <td className="py-2.5 font-semibold text-slate-900">{c.category_name}</td>
                          <td className="py-2.5 text-right font-bold text-slate-900">
                            {formatRupiah(c.amount_total)}
                          </td>
                        </tr>
                      ))
                    ) : (
                      <tr>
                        <td colSpan={2} className="py-6 text-center text-slate-400">
                          Belum ada pengeluaran pada periode ini.
                        </td>
                      </tr>
                    )}
                  </tbody>
                </table>
              </div>
            </Card>
          </div>

          {/* Daily Trends Table */}
          <Card title="Rincian Transaksi Harian" subtitle="Perkembangan pesanan, omzet, kas masuk, dan pengeluaran per hari">
            <div className="overflow-x-auto pt-1">
              <table className="w-full text-left text-xs divide-y divide-slate-100">
                <thead className="text-[11px] font-bold text-slate-400 uppercase bg-slate-50/50">
                  <tr>
                    <th className="p-3">Tanggal</th>
                    <th className="p-3 text-center">Jumlah Pesanan</th>
                    <th className="p-3 text-right">Omzet</th>
                    <th className="p-3 text-right">Kas Masuk</th>
                    <th className="p-3 text-right">Pengeluaran</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-slate-100 text-slate-700">
                  {report.daily_trends && report.daily_trends.length > 0 ? (
                    report.daily_trends.map(d => (
                      <tr key={d.date} className="hover:bg-slate-50/80">
                        <td className="p-3 font-semibold text-slate-900">{d.date}</td>
                        <td className="p-3 text-center text-slate-600">{d.order_count}</td>
                        <td className="p-3 text-right font-bold text-slate-900">{formatRupiah(d.revenue)}</td>
                        <td className="p-3 text-right text-emerald-700 font-semibold">{formatRupiah(d.cash_received)}</td>
                        <td className="p-3 text-right text-slate-600">{formatRupiah(d.expenses)}</td>
                      </tr>
                    ))
                  ) : (
                    <tr>
                      <td colSpan={5} className="py-8 text-center text-slate-400">
                        Belum ada data tren harian pada periode ini.
                      </td>
                    </tr>
                  )}
                </tbody>
              </table>
            </div>
          </Card>
        </>
      )}
    </div>
  );
};
