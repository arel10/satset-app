import React from 'react';
import { NavLink } from 'react-router-dom';
import {
  LayoutDashboard,
  PlusCircle,
  ShoppingBag,
  Users,
  Sparkles,
  ReceiptText,
  TrendingUp,
  Settings,
  X,
  Shirt,
} from 'lucide-react';
import { useAuth } from '@/context/AuthContext';
import { UserRole } from '@/types/api';

interface NavItem {
  to: string;
  label: string;
  icon: React.ReactNode;
  allowedRoles: UserRole[];
}

const NAV_ITEMS: NavItem[] = [
  {
    to: '/dashboard',
    label: 'Dashboard',
    icon: <LayoutDashboard className="w-5 h-5" />,
    allowedRoles: ['owner', 'kasir'],
  },
  {
    to: '/orders/new',
    label: 'Buat Pesanan (POS)',
    icon: <PlusCircle className="w-5 h-5" />,
    allowedRoles: ['owner', 'kasir'],
  },
  {
    to: '/orders',
    label: 'Daftar Pesanan',
    icon: <ShoppingBag className="w-5 h-5" />,
    allowedRoles: ['owner', 'kasir', 'operator'],
  },
  {
    to: '/customers',
    label: 'Pelanggan',
    icon: <Users className="w-5 h-5" />,
    allowedRoles: ['owner', 'kasir'],
  },
  {
    to: '/services',
    label: 'Layanan & Tarif',
    icon: <Shirt className="w-5 h-5" />,
    allowedRoles: ['owner'],
  },
  {
    to: '/expenses',
    label: 'Pengeluaran',
    icon: <ReceiptText className="w-5 h-5" />,
    allowedRoles: ['owner'],
  },
  {
    to: '/reports',
    label: 'Laporan Bisnis',
    icon: <TrendingUp className="w-5 h-5" />,
    allowedRoles: ['owner'],
  },
  {
    to: '/ai',
    label: 'AI Assistant',
    icon: <Sparkles className="w-5 h-5" />,
    allowedRoles: ['owner'],
  },
  {
    to: '/settings',
    label: 'Pengaturan',
    icon: <Settings className="w-5 h-5" />,
    allowedRoles: ['owner'],
  },
];

export interface SidebarProps {
  isOpen: boolean;
  onClose: () => void;
}

export const Sidebar: React.FC<SidebarProps> = ({ isOpen, onClose }) => {
  const { user } = useAuth();
  const role = user?.role || 'operator';

  const visibleItems = NAV_ITEMS.filter(item => item.allowedRoles.includes(role));

  return (
    <>
      {/* Mobile Backdrop */}
      {isOpen && (
        <div
          className="fixed inset-0 z-40 bg-slate-900/40 backdrop-blur-xs lg:hidden"
          onClick={onClose}
        />
      )}

      {/* Sidebar container */}
      <aside
        className={`fixed top-0 bottom-0 left-0 z-50 w-64 bg-white border-r border-slate-200/80 flex flex-col transition-transform duration-200 ease-in-out lg:translate-x-0 ${
          isOpen ? 'translate-x-0' : '-translate-x-full'
        }`}
      >
        {/* Brand header */}
        <div className="h-16 px-6 border-b border-slate-100 flex items-center justify-between">
          <div className="flex items-center gap-2.5">
            <div className="w-9 h-9 rounded-xl bg-teal-700 text-white flex items-center justify-center font-bold text-lg shadow-xs">
              S
            </div>
            <div>
              <span className="font-bold text-slate-900 tracking-tight text-lg">SatSet</span>
              <span className="text-[10px] font-semibold text-teal-700 bg-teal-50 px-1.5 py-0.5 rounded-sm ml-1.5 border border-teal-200 uppercase">
                {role}
              </span>
            </div>
          </div>
          <button
            onClick={onClose}
            className="p-1 rounded-lg text-slate-400 hover:text-slate-600 lg:hidden"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        {/* Navigation list */}
        <nav className="flex-1 px-3 py-4 space-y-1 overflow-y-auto">
          {visibleItems.map(item => (
            <NavLink
              key={item.to}
              to={item.to}
              end={item.to === '/orders' || item.to === '/dashboard'}
              onClick={() => {
                if (window.innerWidth < 1024) onClose();
              }}
              className={({ isActive }) =>
                `flex items-center gap-3 px-3.5 py-2.5 rounded-lg text-sm font-medium transition-colors ${
                  isActive
                    ? 'bg-teal-50 text-teal-800 font-semibold'
                    : 'text-slate-600 hover:bg-slate-50 hover:text-slate-900'
                }`
              }
            >
              {({ isActive }) => (
                <>
                  <span className={isActive ? 'text-teal-700' : 'text-slate-400'}>
                    {item.icon}
                  </span>
                  <span>{item.label}</span>
                </>
              )}
            </NavLink>
          ))}
        </nav>

        {/* User Card at footer */}
        <div className="p-3.5 border-t border-slate-100">
          <div className="flex items-center gap-3 px-2 py-1.5">
            <div className="w-8 h-8 rounded-full bg-slate-100 border border-slate-200 text-slate-700 font-semibold text-xs flex items-center justify-center">
              {user?.full_name?.charAt(0) || 'U'}
            </div>
            <div className="min-w-0 flex-1">
              <p className="text-xs font-semibold text-slate-900 truncate">
                {user?.full_name || user?.username}
              </p>
              <p className="text-[11px] text-slate-400 truncate capitalize">{role}</p>
            </div>
          </div>
        </div>
      </aside>
    </>
  );
};
