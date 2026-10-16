import React, { useEffect, useState } from 'react';
import { NavLink } from 'react-router-dom';
import { auditApi, AuditLogFilterParams } from '@/api/audit';
import { AuditLog } from '@/types/api';
import { formatJakartaDateTime } from '@/utils/date';
import { Card } from '@/components/ui/Card';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';
import { Select } from '@/components/ui/Select';
import { Modal } from '@/components/ui/Modal';
import { Skeleton, EmptyState } from '@/components/ui/EmptyState';
import {
  FileText,
  Store,
  Users,
  Search,
  Filter,
  Eye,
  ChevronLeft,
  ChevronRight,
} from 'lucide-react';

export const AuditLogsPage: React.FC = () => {
  const [logs, setLogs] = useState<AuditLog[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [perPage] = useState(20);
  const [totalPages, setTotalPages] = useState(1);

  // Filters
  const [selectedAction, setSelectedAction] = useState('');
  const [selectedEntity, setSelectedEntity] = useState('');
  const [startDate, setStartDate] = useState('');
  const [endDate, setEndDate] = useState('');

  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Detail Modal
  const [detailLog, setDetailLog] = useState<AuditLog | null>(null);

  const fetchLogs = async () => {
    setIsLoading(true);
    setError(null);
    try {
      const params: AuditLogFilterParams = {
        page,
        per_page: perPage,
        action: selectedAction || undefined,
        entity_type: selectedEntity || undefined,
        start_date: startDate || undefined,
        end_date: endDate || undefined,
      };

      const res = await auditApi.list(params);
      setLogs(res.logs);
      if (res.meta) {
        setTotal(res.meta.total);
        setTotalPages(res.meta.total_pages);
      }
    } catch (err: any) {
      setError(err.message || 'Gagal memuat log audit.');
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchLogs();
  }, [page, selectedAction, selectedEntity, startDate, endDate]);

  const handleResetFilters = () => {
    setSelectedAction('');
    setSelectedEntity('');
    setStartDate('');
    setEndDate('');
    setPage(1);
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-2xl font-bold text-slate-900 tracking-tight">Log Audit Trail</h1>
        <p className="text-xs text-slate-500 mt-1">
          Rekaman jejak aktivitas sistem, perubahan data, transaksi, dan riwayat keamanan.
        </p>
      </div>

      {/* Sub-navigation tabs */}
      <div className="flex items-center gap-2 border-b border-slate-200 pb-1">
        <NavLink
          to="/settings"
          end
          className={({ isActive }) =>
            `px-3.5 py-2 text-xs font-semibold rounded-lg flex items-center gap-2 transition-colors ${
              isActive
                ? 'bg-teal-50 text-teal-800'
                : 'text-slate-600 hover:text-slate-900 hover:bg-slate-50'
            }`
          }
        >
          <Store className="w-4 h-4" />
          Profil Toko &amp; Nota
        </NavLink>
        <NavLink
          to="/settings/users"
          className={({ isActive }) =>
            `px-3.5 py-2 text-xs font-semibold rounded-lg flex items-center gap-2 transition-colors ${
              isActive
                ? 'bg-teal-50 text-teal-800'
                : 'text-slate-600 hover:text-slate-900 hover:bg-slate-50'
            }`
          }
        >
          <Users className="w-4 h-4" />
          Manajemen Pengguna
        </NavLink>
        <NavLink
          to="/settings/audit-logs"
          className={({ isActive }) =>
            `px-3.5 py-2 text-xs font-semibold rounded-lg flex items-center gap-2 transition-colors ${
              isActive
                ? 'bg-teal-50 text-teal-800'
                : 'text-slate-600 hover:text-slate-900 hover:bg-slate-50'
            }`
          }
        >
          <FileText className="w-4 h-4" />
          Log Audit Trail
        </NavLink>
      </div>

      {/* Filter Bar */}
      <Card className="p-4">
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-3">
          <Select
            label="Aksi"
            value={selectedAction}
            onChange={e => {
              setSelectedAction(e.target.value);
              setPage(1);
            }}
            options={[
              { value: '', label: 'Semua Aksi' },
              { value: 'CREATE', label: 'CREATE (Buat)' },
              { value: 'UPDATE', label: 'UPDATE (Ubah)' },
              { value: 'DELETE', label: 'DELETE (Hapus)' },
              { value: 'STATUS_CHANGE', label: 'STATUS_CHANGE (Transisi Status)' },
              { value: 'PAYMENT', label: 'PAYMENT (Pembayaran)' },
              { value: 'PAYMENT_CANCEL', label: 'PAYMENT_CANCEL (Batal Bayar)' },
              { value: 'ORDER_CANCEL', label: 'ORDER_CANCEL (Batal Pesanan)' },
            ]}
          />

          <Select
            label="Entitas"
            value={selectedEntity}
            onChange={e => {
              setSelectedEntity(e.target.value);
              setPage(1);
            }}
            options={[
              { value: '', label: 'Semua Entitas' },
              { value: 'order', label: 'Order (Pesanan)' },
              { value: 'payment', label: 'Payment (Pembayaran)' },
              { value: 'user', label: 'User (Pengguna)' },
              { value: 'customer', label: 'Customer (Pelanggan)' },
              { value: 'service', label: 'Service (Layanan)' },
              { value: 'expense', label: 'Expense (Pengeluaran)' },
              { value: 'setting', label: 'Setting (Pengaturan)' },
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

        <div className="flex justify-end pt-2 border-t border-slate-100 mt-2">
          <Button size="sm" variant="ghost" onClick={handleResetFilters}>
            Reset Filter
          </Button>
        </div>
      </Card>

      {/* Logs Table */}
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
        ) : logs.length === 0 ? (
          <div className="p-8">
            <EmptyState
              icon={<FileText className="w-6 h-6 text-slate-400" />}
              title="Belum Ada Log Audit"
              description="Aktivitas transaksi, pembaruan status, dan perubahan pengaturan akan terekam di sini."
            />
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-left text-xs divide-y divide-slate-200">
              <thead className="bg-slate-50 text-[11px] font-bold text-slate-700 uppercase tracking-wider">
                <tr>
                  <th className="px-5 py-3.5">Waktu</th>
                  <th className="px-5 py-3.5">Pengguna</th>
                  <th className="px-5 py-3.5">Aksi</th>
                  <th className="px-5 py-3.5">Entitas</th>
                  <th className="px-5 py-3.5">IP Address</th>
                  <th className="px-5 py-3.5 text-right">Detail</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-100 bg-white">
                {logs.map(log => (
                  <tr key={log.id} className="hover:bg-slate-50/80 transition-colors">
                    <td className="px-5 py-3.5 font-medium text-slate-900 whitespace-nowrap">
                      {formatJakartaDateTime(log.created_at)}
                    </td>
                    <td className="px-5 py-3.5">
                      <span className="font-semibold text-slate-800 block">
                        {log.user?.full_name || 'System / Guest'}
                      </span>
                      <span className="text-[10px] text-slate-400 uppercase font-semibold">
                        {log.user?.role || '-'}
                      </span>
                    </td>
                    <td className="px-5 py-3.5">
                      <span className="bg-slate-100 text-slate-800 font-mono text-[11px] px-2 py-0.5 rounded font-bold">
                        {log.action}
                      </span>
                    </td>
                    <td className="px-5 py-3.5">
                      <span className="text-slate-700 capitalize font-medium">
                        {log.entity_type} #{log.entity_id}
                      </span>
                    </td>
                    <td className="px-5 py-3.5 text-slate-500 font-mono text-[11px]">
                      {log.ip_address || '-'}
                    </td>
                    <td className="px-5 py-3.5 text-right whitespace-nowrap">
                      <Button
                        size="sm"
                        variant="ghost"
                        onClick={() => setDetailLog(log)}
                        leftIcon={<Eye className="w-3.5 h-3.5" />}
                      >
                        Lihat Data
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
              Menampilkan {logs.length} dari total <span className="font-semibold">{total}</span> baris audit
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

      {/* Modal: JSON Values Diff */}
      <Modal
        isOpen={detailLog !== null}
        onClose={() => setDetailLog(null)}
        title={`Detail Audit: ${detailLog?.action} (${detailLog?.entity_type} #${detailLog?.entity_id})`}
        maxWidth="lg"
      >
        {detailLog && (
          <div className="space-y-4 text-xs">
            <div className="p-3 bg-slate-50 rounded-lg border border-slate-200 grid grid-cols-2 gap-2 text-[11px]">
              <div>
                <span className="text-slate-400 block">Pengguna:</span>
                <span className="font-semibold text-slate-800">
                  {detailLog.user?.full_name} (@{detailLog.user?.username})
                </span>
              </div>
              <div>
                <span className="text-slate-400 block">Waktu &amp; IP:</span>
                <span className="font-semibold text-slate-800">
                  {formatJakartaDateTime(detailLog.created_at)} ({detailLog.ip_address})
                </span>
              </div>
            </div>

            {detailLog.old_values && Object.keys(detailLog.old_values).length > 0 && (
              <div>
                <span className="font-bold text-slate-700 block mb-1 text-[11px]">Nilai Sebelumnya:</span>
                <pre className="p-3 rounded-lg bg-slate-900 text-slate-100 overflow-x-auto text-[11px] font-mono">
                  {JSON.stringify(detailLog.old_values, null, 2)}
                </pre>
              </div>
            )}

            {detailLog.new_values && Object.keys(detailLog.new_values).length > 0 && (
              <div>
                <span className="font-bold text-slate-700 block mb-1 text-[11px]">Nilai Baru:</span>
                <pre className="p-3 rounded-lg bg-slate-900 text-slate-100 overflow-x-auto text-[11px] font-mono">
                  {JSON.stringify(detailLog.new_values, null, 2)}
                </pre>
              </div>
            )}
          </div>
        )}
      </Modal>
    </div>
  );
};
