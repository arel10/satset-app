import React from 'react';
import { Navigate, Outlet } from 'react-router-dom';
import { useAuth } from '@/context/AuthContext';
import { UserRole } from '@/types/api';
import { ShieldAlert, Loader2 } from 'lucide-react';
import { Button } from '@/components/ui/Button';

export interface ProtectedRouteProps {
  allowedRoles?: UserRole[];
}

export const ProtectedRoute: React.FC<ProtectedRouteProps> = ({ allowedRoles }) => {
  const { user, role, isLoading, isAuthenticated } = useAuth();

  if (isLoading) {
    return (
      <div className="min-h-screen flex flex-col items-center justify-center bg-slate-50 gap-3">
        <Loader2 className="w-8 h-8 text-teal-700 animate-spin" />
        <p className="text-xs text-slate-500 font-medium">Memuat sistem SatSet...</p>
      </div>
    );
  }

  if (!isAuthenticated || !user) {
    return <Navigate to="/login" replace />;
  }

  if (allowedRoles && role && !allowedRoles.includes(role)) {
    return (
      <div className="min-h-[60vh] flex flex-col items-center justify-center text-center p-6">
        <div className="w-16 h-16 rounded-full bg-rose-50 text-rose-600 flex items-center justify-center mb-4">
          <ShieldAlert className="w-8 h-8" />
        </div>
        <h2 className="text-xl font-bold text-slate-900">Akses Dibatasi</h2>
        <p className="text-sm text-slate-500 max-w-md mt-1.5 mb-6">
          Halaman ini khusus untuk peran yang memiliki otorisasi. Akun Anda saat ini memiliki peran{' '}
          <span className="font-semibold text-slate-700 uppercase">{role}</span>.
        </p>
        <Button variant="primary" onClick={() => window.history.back()}>
          Kembali ke Halaman Sebelumnya
        </Button>
      </div>
    );
  }

  return <Outlet />;
};
