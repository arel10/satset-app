import React, { useEffect, useState } from 'react';
import { servicesApi } from '@/api/services';
import { LaundryService } from '@/types/api';
import { formatRupiah } from '@/utils/currency';
import { Card } from '@/components/ui/Card';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';
import { Select } from '@/components/ui/Select';
import { Modal } from '@/components/ui/Modal';
import { Skeleton, EmptyState } from '@/components/ui/EmptyState';
import { Shirt, Plus, Edit2, CheckCircle2, XCircle } from 'lucide-react';

export const ServiceListPage: React.FC = () => {
  const [services, setServices] = useState<LaundryService[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Form Modal State
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingService, setEditingService] = useState<LaundryService | null>(null);
  const [formData, setFormData] = useState({
    name: '',
    category: 'Kiloan',
    unit: 'kg',
    price_per_unit: 8000,
    estimated_duration: 24,
    is_active: true,
  });
  const [isSubmitting, setIsSubmitting] = useState(false);

  const fetchServices = async () => {
    setIsLoading(true);
    setError(null);
    try {
      const res = await servicesApi.list(1, 100);
      setServices(res.services);
    } catch (err: any) {
      setError(err.message || 'Gagal memuat katalog layanan.');
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchServices();
  }, []);

  const handleOpenAddModal = () => {
    setEditingService(null);
    setFormData({
      name: '',
      category: 'Kiloan',
      unit: 'kg',
      price_per_unit: 8000,
      estimated_duration: 24,
      is_active: true,
    });
    setIsModalOpen(true);
  };

  const handleOpenEditModal = (s: LaundryService) => {
    setEditingService(s);
    setFormData({
      name: s.name,
      category: s.category || 'Kiloan',
      unit: s.unit || 'kg',
      price_per_unit: s.price_per_unit,
      estimated_duration: s.estimated_duration || 24,
      is_active: s.is_active,
    });
    setIsModalOpen(true);
  };

  const handleSaveService = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!formData.name.trim()) return;

    setIsSubmitting(true);
    try {
      if (editingService) {
        await servicesApi.update(editingService.id, formData);
      } else {
        await servicesApi.create(formData);
      }
      setIsModalOpen(false);
      fetchServices();
    } catch (err: any) {
      alert(err.message || 'Gagal menyimpan layanan.');
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleToggleActive = async (s: LaundryService) => {
    try {
      await servicesApi.update(s.id, { is_active: !s.is_active });
      fetchServices();
    } catch (err: any) {
      alert(err.message || 'Gagal mengubah status layanan.');
    }
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-slate-900 tracking-tight">Katalog Layanan &amp; Tarif</h1>
          <p className="text-xs text-slate-500 mt-1">
            Atur jenis cucian, satuan, tarif per unit, dan estimasi durasi pengerjaan.
          </p>
        </div>
        <Button variant="primary" onClick={handleOpenAddModal} leftIcon={<Plus className="w-4 h-4" />}>
          Tambah Layanan Baru
        </Button>
      </div>

      {/* Services Table Card */}
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
        ) : services.length === 0 ? (
          <div className="p-8">
            <EmptyState
              icon={<Shirt className="w-6 h-6 text-slate-400" />}
              title="Belum Ada Layanan Laundry"
              description="Tambahkan katalog layanan kiloan, satuan, atau express untuk memulai transaksi POS."
              action={
                <Button variant="primary" size="sm" onClick={handleOpenAddModal}>
                  Tambah Layanan
                </Button>
              }
            />
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-left text-xs divide-y divide-slate-200">
              <thead className="bg-slate-50 text-[11px] font-bold text-slate-700 uppercase tracking-wider">
                <tr>
                  <th className="px-5 py-3.5">Nama Layanan</th>
                  <th className="px-5 py-3.5">Kategori</th>
                  <th className="px-5 py-3.5">Satuan</th>
                  <th className="px-5 py-3.5 text-right">Tarif / Satuan</th>
                  <th className="px-5 py-3.5 text-center">Durasi (Jam)</th>
                  <th className="px-5 py-3.5 text-center">Status</th>
                  <th className="px-5 py-3.5 text-right">Aksi</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-100 bg-white">
                {services.map(s => (
                  <tr key={s.id} className="hover:bg-slate-50/80 transition-colors">
                    <td className="px-5 py-3.5 font-bold text-slate-900">{s.name}</td>
                    <td className="px-5 py-3.5 text-slate-600">
                      <span className="bg-slate-100 text-slate-700 px-2 py-0.5 rounded text-[11px] font-medium">
                        {s.category || 'Umum'}
                      </span>
                    </td>
                    <td className="px-5 py-3.5 font-medium text-slate-700">{s.unit}</td>
                    <td className="px-5 py-3.5 text-right font-bold text-teal-800">
                      {formatRupiah(s.price_per_unit)}
                    </td>
                    <td className="px-5 py-3.5 text-center text-slate-600">
                      {s.estimated_duration} jam
                    </td>
                    <td className="px-5 py-3.5 text-center">
                      <button
                        type="button"
                        onClick={() => handleToggleActive(s)}
                        className={`inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-[11px] font-medium border transition-colors ${
                          s.is_active
                            ? 'bg-emerald-50 text-emerald-700 border-emerald-200 hover:bg-emerald-100'
                            : 'bg-slate-100 text-slate-500 border-slate-200 hover:bg-slate-200'
                        }`}
                      >
                        {s.is_active ? (
                          <>
                            <CheckCircle2 className="w-3 h-3 text-emerald-600" /> Aktif
                          </>
                        ) : (
                          <>
                            <XCircle className="w-3 h-3 text-slate-400" /> Nonaktif
                          </>
                        )}
                      </button>
                    </td>
                    <td className="px-5 py-3.5 text-right whitespace-nowrap">
                      <Button
                        size="sm"
                        variant="outline"
                        onClick={() => handleOpenEditModal(s)}
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
      </Card>

      {/* Modal: Tambah/Edit Layanan */}
      <Modal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        title={editingService ? 'Edit Layanan Laundry' : 'Tambah Layanan Baru'}
      >
        <form onSubmit={handleSaveService} className="space-y-3.5">
          <Input
            label="Nama Layanan"
            required
            autoFocus
            value={formData.name}
            onChange={e => setFormData({ ...formData, name: e.target.value })}
            placeholder="Contoh: Cuci Kering Setrika Reguler"
          />

          <div className="grid grid-cols-2 gap-3">
            <Input
              label="Kategori"
              required
              value={formData.category}
              onChange={e => setFormData({ ...formData, category: e.target.value })}
              placeholder="Contoh: Kiloan / Satuan"
            />
            <Input
              label="Satuan Unit"
              required
              value={formData.unit}
              onChange={e => setFormData({ ...formData, unit: e.target.value })}
              placeholder="Contoh: kg / pcs / meter"
            />
          </div>

          <div className="grid grid-cols-2 gap-3">
            <Input
              label="Tarif per Satuan (Rp)"
              required
              type="number"
              min="0"
              value={formData.price_per_unit}
              onChange={e =>
                setFormData({ ...formData, price_per_unit: parseInt(e.target.value, 10) || 0 })
              }
            />
            <Input
              label="Estimasi Durasi (Jam)"
              required
              type="number"
              min="1"
              value={formData.estimated_duration}
              onChange={e =>
                setFormData({ ...formData, estimated_duration: parseInt(e.target.value, 10) || 24 })
              }
            />
          </div>

          <div className="pt-1">
            <label className="flex items-center gap-2 cursor-pointer text-xs font-semibold text-slate-700">
              <input
                type="checkbox"
                checked={formData.is_active}
                onChange={e => setFormData({ ...formData, is_active: e.target.checked })}
                className="w-4 h-4 rounded text-teal-700 focus:ring-teal-600 border-slate-300"
              />
              <span>Tampilkan layanan ini di kasir POS (Aktif)</span>
            </label>
          </div>

          <div className="flex justify-end gap-2 pt-2">
            <Button type="button" variant="outline" onClick={() => setIsModalOpen(false)}>
              Batal
            </Button>
            <Button type="submit" variant="primary" isLoading={isSubmitting}>
              Simpan Layanan
            </Button>
          </div>
        </form>
      </Modal>
    </div>
  );
};
