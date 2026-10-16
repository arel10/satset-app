import React from 'react';
import { Menu, LogOut, Plus } from 'lucide-react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '@/context/AuthContext';
import { Button } from '@/components/ui/Button';

export interface HeaderProps {
  onMenuClick: () => void;
}

export const Header: React.FC<HeaderProps> = ({ onMenuClick }) => {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = async () => {
    if (window.confirm('Apakah Anda yakin ingin keluar dari sistem SatSet?')) {
      await logout();
      navigate('/login');
    }
  };

  const isPOSAvailable = user?.role === 'owner' || user?.role === 'kasir';

  return (
    <header className="h-16 bg-white border-b border-slate-200/80 px-4 lg:px-8 flex items-center justify-between sticky top-0 z-30">
      <div className="flex items-center gap-3">
        <button
          onClick={onMenuClick}
          className="p-2 rounded-lg text-slate-500 hover:bg-slate-100 lg:hidden"
          aria-label="Buka Menu"
        >
          <Menu className="w-5 h-5" />
        </button>
        <span className="text-sm font-semibold text-slate-700 hidden sm:inline-block">
          Sistem Operasional Laundry
        </span>
      </div>

      <div className="flex items-center gap-2.5">
        {isPOSAvailable && (
          <Button
            size="sm"
            variant="primary"
            onClick={() => navigate('/orders/new')}
            leftIcon={<Plus className="w-4 h-4" />}
          >
            Buat Pesanan
          </Button>
        )}

        <Button
          size="sm"
          variant="ghost"
          onClick={handleLogout}
          leftIcon={<LogOut className="w-4 h-4 text-slate-500" />}
          className="text-slate-600 hover:text-rose-600 hover:bg-rose-50"
        >
          Keluar
        </Button>
      </div>
    </header>
  );
};
