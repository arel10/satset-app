import React, { useEffect, useState } from 'react';
import { NavLink } from 'react-router-dom';
import { usersApi } from '@/api/users';
import { User, UserRole } from '@/types/api';
import { formatJakartaDate } from '@/utils/date';
import { Card } from '@/components/ui/Card';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';
import { Select } from '@/components/ui/Select';
import { Modal } from '@/components/ui/Modal';
import { Skeleton, EmptyState } from '@/components/ui/EmptyState';
import {
  Users,
  Store,
  FileText,
  Plus,
  Edit2,
  Trash2,
  CheckCircle2,
  XCircle,
  AlertCircle,
} from 'lucide-react';

export const UsersPage: React.FC = () => {
  const [users, setUsers] = useState<User[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Add/Edit Modal
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingUser, setEditingUser] = useState<User | null>(null);
  const [formData, setFormData] = useState({
    username: '',
    email: '',
    full_name: '',
    role: 'kasir' as UserRole,
    password: '',
    is_active: true,
  });
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [formError, setFormError] = useState<string | null>(null);

  const fetchUsers = async () => {
    setIsLoading(true);
    setError(null);
    try {
      const res = await usersApi.list(1, 50);
      setUsers(res.users);
    } catch (err: any) {
      setError(err.message || 'Gagal memuat pengguna.');
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchUsers();
  }, []);

  const handleOpenAdd = () => {
    setEditingUser(null);
    setFormError(null);
    setFormData({
      username: '',
      email: '',
      full_name: '',
      role: 'kasir',
      password: '',
      is_active: true,
    });
    setIsModalOpen(true);
  };

  const handleOpenEdit = (u: User) => {
    setEditingUser(u);
    setFormError(null);
    setFormData({
      username: u.username,
      email: u.email,
      full_name: u.full_name,
      role: u.role,
      password: '',
      is_active: u.is_active,
    });
    setIsModalOpen(true);
  };

  const handleSave = async (e: React.FormEvent) => {
    e.preventDefault();
    setFormError(null);
    setIsSubmitting(true);

    try {
      if (editingUser) {
        await usersApi.update(editingUser.id, {
          email: formData.email,
          full_name: formData.full_name,
          role: formData.role,
          is_active: formData.is_active,
          password: formData.password || undefined,
        });
      } else {
        if (!formData.password) {
          throw new Error('Kata sandi wajib diisi untuk pengguna baru.');
        }
        await usersApi.create({
          username: formData.username,
          email: formData.email,
          full_name: formData.full_name,
          role: formData.role,
          password: formData.password,
        });
      }
      setIsModalOpen(false);
      fetchUsers();
    } catch (err: any) {
      setFormError(err.message || 'Gagal menyimpan akun pengguna.');
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleDelete = async (u: User) => {
    if (!window.confirm(`Hapus pengguna '${u.full_name || u.username}'?`)) {
      return;
    }
    try {
      await usersApi.delete(u.id);
      fetchUsers();
    } catch (err: any) {
      alert(err.message || 'Gagal menghapus pengguna.');
    }
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-slate-900 tracking-tight">Manajemen Pengguna</h1>
          <p className="text-xs text-slate-500 mt-1">
            Kelola akses Owner, Kasir, dan Operator sistem laundry.
          </p>
        </div>
        <Button variant="primary" onClick={handleOpenAdd} leftIcon={<Plus className="w-4 h-4" />}>
          Tambah Pengguna
        </Button>
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

      {/* Users Table */}
      <Card className="overflow-hidden p-0">
        {isLoading ? (
          <div className="p-6 space-y-3">
            {Array.from({ length: 4 }).map((_, i) => (
              <Skeleton key={i} className="h-12 w-full" />
            ))}
          </div>
        ) : error ? (
          <div className="p-8 text-center text-xs text-rose-700 bg-rose-50">
            {error}
          </div>
        ) : users.length === 0 ? (
          <div className="p-8">
            <EmptyState
              icon={<Users className="w-6 h-6 text-slate-400" />}
              title="Belum Ada Pengguna"
              description="Tambahkan akun kasir atau operator untuk operasional harian."
              action={
                <Button variant="primary" size="sm" onClick={handleOpenAdd}>
                  Tambah Pengguna
                </Button>
              }
            />
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-left text-xs divide-y divide-slate-200">
              <thead className="bg-slate-50 text-[11px] font-bold text-slate-700 uppercase tracking-wider">
                <tr>
                  <th className="px-5 py-3.5">Nama &amp; Username</th>
                  <th className="px-5 py-3.5">Email</th>
                  <th className="px-5 py-3.5">Peran (Role)</th>
                  <th className="px-5 py-3.5 text-center">Status</th>
                  <th className="px-5 py-3.5">Terdaftar</th>
                  <th className="px-5 py-3.5 text-right">Aksi</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-100 bg-white">
                {users.map(u => (
                  <tr key={u.id} className="hover:bg-slate-50/80 transition-colors">
                    <td className="px-5 py-3.5 font-bold text-slate-900">
                      {u.full_name}
                      <span className="block text-[11px] font-normal text-slate-400">
                        @{u.username}
                      </span>
                    </td>
                    <td className="px-5 py-3.5 text-slate-600">{u.email}</td>
                    <td className="px-5 py-3.5">
                      <span
                        className={`inline-flex items-center px-2 py-0.5 rounded text-[11px] font-bold uppercase ${
                          u.role === 'owner'
                            ? 'bg-purple-50 text-purple-700 border border-purple-200'
                            : u.role === 'kasir'
                            ? 'bg-teal-50 text-teal-700 border border-teal-200'
                            : 'bg-blue-50 text-blue-700 border border-blue-200'
                        }`}
                      >
                        {u.role}
                      </span>
                    </td>
                    <td className="px-5 py-3.5 text-center">
                      <span
                        className={`inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-[11px] font-medium border ${
                          u.is_active
                            ? 'bg-emerald-50 text-emerald-700 border-emerald-200'
                            : 'bg-slate-100 text-slate-500 border-slate-200'
                        }`}
                      >
                        {u.is_active ? (
                          <>
                            <CheckCircle2 className="w-3 h-3 text-emerald-600" /> Aktif
                          </>
                        ) : (
                          <>
                            <XCircle className="w-3 h-3 text-slate-400" /> Nonaktif
                          </>
                        )}
                      </span>
                    </td>
                    <td className="px-5 py-3.5 text-slate-500 whitespace-nowrap">
                      {formatJakartaDate(u.created_at)}
                    </td>
                    <td className="px-5 py-3.5 text-right whitespace-nowrap space-x-1.5">
                      <Button
                        size="sm"
                        variant="outline"
                        onClick={() => handleOpenEdit(u)}
                        leftIcon={<Edit2 className="w-3.5 h-3.5" />}
                      >
                        Edit
                      </Button>
                      <Button
                        size="sm"
                        variant="ghost"
                        onClick={() => handleDelete(u)}
                        className="text-rose-600 hover:text-rose-800 hover:bg-rose-50"
                        leftIcon={<Trash2 className="w-3.5 h-3.5" />}
                      >
                        Hapus
                      </Button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </Card>

      {/* Modal: Tambah/Edit Pengguna */}
      <Modal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        title={editingUser ? 'Edit Akun Pengguna' : 'Tambah Pengguna Baru'}
      >
        <form onSubmit={handleSave} className="space-y-3.5">
          {formError && (
            <div className="p-3 rounded-lg bg-rose-50 border border-rose-200 text-xs text-rose-800 flex items-start gap-2">
              <AlertCircle className="w-4 h-4 text-rose-600 shrink-0 mt-0.5" />
              <span>{formError}</span>
            </div>
          )}

          {!editingUser && (
            <Input
              label="Nama Pengguna (Username)"
              required
              autoFocus
              value={formData.username}
              onChange={e => setFormData({ ...formData, username: e.target.value })}
              placeholder="Contoh: kasir1"
            />
          )}

          <Input
            label="Nama Lengkap"
            required
            value={formData.full_name}
            onChange={e => setFormData({ ...formData, full_name: e.target.value })}
            placeholder="Contoh: Rian Pratama"
          />

          <Input
            label="Alamat Email"
            type="email"
            required
            value={formData.email}
            onChange={e => setFormData({ ...formData, email: e.target.value })}
            placeholder="Contoh: rian@satset.local"
          />

          <Select
            label="Peran Akun (Role)"
            required
            value={formData.role}
            onChange={e => setFormData({ ...formData, role: e.target.value as UserRole })}
            options={[
              { value: 'kasir', label: 'Kasir — POS & Pelanggan' },
              { value: 'operator', label: 'Operator — Pengerjaan Cucian' },
              { value: 'owner', label: 'Owner — Akses Penuh' },
            ]}
          />

          <Input
            label={editingUser ? 'Ganti Kata Sandi (Kosongkan jika tidak diganti)' : 'Kata Sandi'}
            type="password"
            required={!editingUser}
            value={formData.password}
            onChange={e => setFormData({ ...formData, password: e.target.value })}
            placeholder="••••••••"
          />

          {editingUser && (
            <div className="pt-1">
              <label className="flex items-center gap-2 cursor-pointer text-xs font-semibold text-slate-700">
                <input
                  type="checkbox"
                  checked={formData.is_active}
                  onChange={e => setFormData({ ...formData, is_active: e.target.checked })}
                  className="w-4 h-4 rounded text-teal-700 focus:ring-teal-600 border-slate-300"
                />
                <span>Akun Aktif (Dapat Login)</span>
              </label>
            </div>
          )}

          <div className="flex justify-end gap-2 pt-2">
            <Button type="button" variant="outline" onClick={() => setIsModalOpen(false)}>
              Batal
            </Button>
            <Button type="submit" variant="primary" isLoading={isSubmitting}>
              {editingUser ? 'Simpan Perubahan' : 'Buat Pengguna'}
            </Button>
          </div>
        </form>
      </Modal>
    </div>
  );
};
