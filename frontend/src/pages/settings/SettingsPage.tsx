import React, { useEffect, useState } from 'react';
import { NavLink } from 'react-router-dom';
import { settingsApi, BusinessSettings } from '@/api/settings';
import { Card } from '@/components/ui/Card';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';
import { Skeleton } from '@/components/ui/EmptyState';
import { Store, Users, FileText, Upload, CheckCircle2, Image as ImageIcon } from 'lucide-react';

export const SettingsPage: React.FC = () => {
  const [settings, setSettings] = useState<BusinessSettings>({});
  const [isLoading, setIsLoading] = useState(true);
  const [isSaving, setIsSaving] = useState(false);
  const [saveSuccess, setSaveSuccess] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Logo file upload state
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [isUploadingLogo, setIsUploadingLogo] = useState(false);
  const [previewLogoUrl, setPreviewLogoUrl] = useState<string | null>(null);

  useEffect(() => {
    const fetchSettings = async () => {
      setIsLoading(true);
      try {
        const data = await settingsApi.getAll();
        setSettings({
          business_name: data.business_name || '',
          business_address: data.business_address || '',
          business_phone: data.business_phone || '',
          receipt_footer: data.receipt_footer || '',
          business_logo_url: data.business_logo_url || '',
        });
        if (data.business_logo_url) {
          setPreviewLogoUrl(data.business_logo_url);
        }
      } catch (err: any) {
        setError(err.message || 'Gagal memuat pengaturan toko.');
      } finally {
        setIsLoading(false);
      }
    };
    fetchSettings();
  }, []);

  const handleSaveTextSettings = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSaving(true);
    setSaveSuccess(false);
    setError(null);
    try {
      await settingsApi.update(settings);
      setSaveSuccess(true);
      setTimeout(() => setSaveSuccess(false), 3000);
    } catch (err: any) {
      setError(err.message || 'Gagal menyimpan pengaturan.');
    } finally {
      setIsSaving(false);
    }
  };

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      const file = e.target.files[0];
      if (file.size > 2 * 1024 * 1024) {
        alert('Ukuran file maksimal 2MB');
        return;
      }
      setSelectedFile(file);
      setPreviewLogoUrl(URL.createObjectURL(file));
    }
  };

  const handleUploadLogo = async () => {
    if (!selectedFile) return;
    setIsUploadingLogo(true);
    setError(null);
    try {
      const res = await settingsApi.uploadLogo(selectedFile);
      setSettings(prev => ({ ...prev, business_logo_url: res.logo_url }));
      setPreviewLogoUrl(res.logo_url);
      setSelectedFile(null);
      alert('Logo bisnis berhasil diunggah dan disimpan!');
    } catch (err: any) {
      alert(err.message || 'Gagal mengunggah logo.');
    } finally {
      setIsUploadingLogo(false);
    }
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-2xl font-bold text-slate-900 tracking-tight">Pengaturan Sistem</h1>
        <p className="text-xs text-slate-500 mt-1">
          Konfigurasi profil toko, logo nota, pengguna aplikasi, dan log audit trail.
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

      {isLoading ? (
        <div className="space-y-4">
          <Skeleton className="h-64" />
          <Skeleton className="h-44" />
        </div>
      ) : (
        <div className="grid grid-cols-1 lg:grid-cols-12 gap-6 items-start">
          {/* Left Column: Business Form (7 Cols) */}
          <div className="lg:col-span-7">
            <Card title="Identitas Bisnis &amp; Teks Nota">
              <form onSubmit={handleSaveTextSettings} className="space-y-4">
                {saveSuccess && (
                  <div className="p-3 rounded-lg bg-emerald-50 border border-emerald-200 text-xs text-emerald-800 flex items-center gap-2">
                    <CheckCircle2 className="w-4 h-4 text-emerald-600 shrink-0" />
                    <span>Pengaturan toko berhasil disimpan!</span>
                  </div>
                )}
                {error && (
                  <div className="p-3 rounded-lg bg-rose-50 border border-rose-200 text-xs text-rose-800">
                    {error}
                  </div>
                )}

                <Input
                  label="Nama Bisnis / Toko Laundry"
                  required
                  value={settings.business_name || ''}
                  onChange={e => setSettings({ ...settings, business_name: e.target.value })}
                  placeholder="Contoh: SatSet Laundry Express"
                />

                <Input
                  label="Alamat Toko (Muncul pada Struk)"
                  value={settings.business_address || ''}
                  onChange={e => setSettings({ ...settings, business_address: e.target.value })}
                  placeholder="Contoh: Jl. Diponegoro No. 45, Surabaya"
                />

                <Input
                  label="Nomor Telepon / WhatsApp Toko"
                  value={settings.business_phone || ''}
                  onChange={e => setSettings({ ...settings, business_phone: e.target.value })}
                  placeholder="Contoh: 0812-3456-7890"
                />

                <Input
                  label="Catatan Kaki Nota (Receipt Footer)"
                  value={settings.receipt_footer || ''}
                  onChange={e => setSettings({ ...settings, receipt_footer: e.target.value })}
                  placeholder="Contoh: Terima kasih atas kepercayaan Anda!"
                  helperText="Teks ini akan tercetak di bagian paling bawah struk thermal pelanggan."
                />

                <div className="flex justify-end pt-2">
                  <Button type="submit" variant="primary" isLoading={isSaving}>
                    Simpan Pengaturan
                  </Button>
                </div>
              </form>
            </Card>
          </div>

          {/* Right Column: Logo Upload (5 Cols) */}
          <div className="lg:col-span-5">
            <Card title="Logo Nota Thermal">
              <div className="space-y-4 text-center">
                <div className="w-28 h-28 mx-auto rounded-2xl border-2 border-dashed border-slate-300 flex items-center justify-center p-2 bg-slate-50 overflow-hidden">
                  {previewLogoUrl ? (
                    <img
                      src={previewLogoUrl}
                      alt="Logo Toko"
                      className="w-full h-full object-contain"
                    />
                  ) : (
                    <ImageIcon className="w-10 h-10 text-slate-300" />
                  )}
                </div>

                <div className="text-xs text-slate-500">
                  <p className="font-semibold text-slate-700">Unggah Logo Toko</p>
                  <p className="text-[11px] text-slate-400 mt-0.5">
                    Format: JPG, PNG, atau WebP (Maksimal 2MB). Disarankan berlatar transparan/putih.
                  </p>
                </div>

                <input
                  type="file"
                  id="logo-upload-input"
                  accept=".jpg,.jpeg,.png,.webp"
                  onChange={handleFileChange}
                  className="hidden"
                />

                <div className="flex flex-col gap-2">
                  <Button
                    type="button"
                    variant="outline"
                    size="sm"
                    onClick={() => document.getElementById('logo-upload-input')?.click()}
                  >
                    Pilih File Gambar
                  </Button>

                  {selectedFile && (
                    <Button
                      type="button"
                      variant="primary"
                      size="sm"
                      isLoading={isUploadingLogo}
                      onClick={handleUploadLogo}
                      leftIcon={<Upload className="w-3.5 h-3.5" />}
                    >
                      Unggah Logo Ini
                    </Button>
                  )}
                </div>
              </div>
            </Card>
          </div>
        </div>
      )}
    </div>
  );
};
