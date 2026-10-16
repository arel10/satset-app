import React from 'react';
import { NavLink } from 'react-router-dom';
import {
  LayoutDashboard,
  PlusCircle,
  ShoppingBag,
  Users,
  Shirt,
  ReceiptText,
  TrendingUp,
  Sparkles,
  Settings,
  X,
} from 'lucide-react';
import { useAuth } from '@/context/AuthContext';
import { UserRole } from '@/types/api';

interface NavItem {
  to: string;
  label: string;
  icon: React.ReactNode;
  allowedRoles: UserRole[];
  badge?: string;
}

interface NavCategory {
  title: string;
  items: NavItem[];
}

const NAV_CATEGORIES: NavCategory[] = [
  {
    title: 'Menu Utama',
    items: [
      {
        to: '/dashboard',
        label: 'Dashboard',
        icon: <LayoutDashboard className="w-4 h-4" />,
        allowedRoles: ['owner', 'kasir'],
      },
      {
        to: '/orders/new',
        label: 'Buat Pesanan (POS)',
        icon: <PlusCircle className="w-4 h-4" />,
        allowedRoles: ['owner', 'kasir'],
        badge: 'POS',
      },
      {
        to: '/orders',
        label: 'Daftar Pesanan',
        icon: <ShoppingBag className="w-4 h-4" />,
        allowedRoles: ['owner', 'kasir', 'operator'],
      },
    ],
  },
  {
    title: 'Operasional',
    items: [
      {
        to: '/customers',
        label: 'Pelanggan',
        icon: <Users className="w-4 h-4" />,
        allowedRoles: ['owner', 'kasir'],
      },
      {
        to: '/services',
        label: 'Layanan & Tarif',
        icon: <Shirt className="w-4 h-4" />,
        allowedRoles: ['owner'],
      },
    ],
  },
  {
    title: 'Keuangan & Laporan',
    items: [
      {
        to: '/expenses',
        label: 'Pengeluaran',
        icon: <ReceiptText className="w-4 h-4" />,
        allowedRoles: ['owner'],
      },
      {
        to: '/reports',
        label: 'Laporan Bisnis',
        icon: <TrendingUp className="w-4 h-4" />,
        allowedRoles: ['owner'],
      },
    ],
  },
  {
    title: 'Sistem & AI',
    items: [
      {
        to: '/ai',
        label: 'AI Assistant',
        icon: <Sparkles className="w-4 h-4" />,
        allowedRoles: ['owner'],
        badge: 'AI',
      },
      {
        to: '/settings',
        label: 'Pengaturan',
        icon: <Settings className="w-4 h-4" />,
        allowedRoles: ['owner'],
      },
    ],
  },
];

export interface SidebarProps {
  isOpen: boolean;
  onClose: () => void;
}

export const Sidebar: React.FC<SidebarProps> = ({ isOpen, onClose }) => {
  const { user } = useAuth();
  const role = user?.role || 'operator';

  // Filter categories and items according to user's role
  const visibleCategories = NAV_CATEGORIES.map(cat => ({
    ...cat,
    items: cat.items.filter(item => item.allowedRoles.includes(role)),
  })).filter(cat => cat.items.length > 0);

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
        className={`fixed top-0 bottom-0 left-0 z-50 w-60 bg-white border-r border-slate-200/80 flex flex-col transition-transform duration-200 ease-in-out lg:translate-x-0 ${
          isOpen ? 'translate-x-0' : '-translate-x-full'
        }`}
      >
        {/* Brand header */}
        <div className="h-14 px-4 border-b border-slate-100 flex items-center justify-between shrink-0">
          <div className="flex items-center gap-2.5">
            <div className="w-7 h-7 rounded-lg bg-teal-700 text-white flex items-center justify-center font-bold text-sm shadow-xs">
              S
            </div>
            <div className="flex items-center gap-1.5">
              <span className="font-bold text-slate-900 tracking-tight text-base">SatSet</span>
              <span className="text-[9px] font-semibold text-teal-700 bg-teal-50 px-1.5 py-0.5 rounded border border-teal-200 uppercase tracking-wide">
                {role}
              </span>
            </div>
          </div>
          <button
            onClick={onClose}
            className="p-1 rounded-md text-slate-400 hover:text-slate-600 lg:hidden"
          >
            <X className="w-4 h-4" />
          </button>
        </div>

        {/* Navigation list with Categories */}
        <nav className="flex-1 px-2.5 py-3 space-y-4 overflow-y-auto">
          {visibleCategories.map(cat => (
            <div key={cat.title} className="space-y-1">
              <div className="px-2.5 pt-1 pb-1 text-[10px] font-bold uppercase tracking-wider text-slate-400 select-none">
                {cat.title}
              </div>
              <div className="space-y-0.5">
                {cat.items.map(item => (
                  <NavLink
                    key={item.to}
                    to={item.to}
                    end={item.to === '/orders' || item.to === '/dashboard'}
                    onClick={() => {
                      if (window.innerWidth < 1024) onClose();
                    }}
                    className={({ isActive }) =>
                      `flex items-center justify-between px-2.5 py-1.5 rounded-lg text-xs font-medium transition-all group ${
                        isActive
                          ? 'bg-teal-50 text-teal-800 font-semibold shadow-2xs border-l-2 border-teal-600'
                          : 'text-slate-600 hover:bg-slate-50 hover:text-slate-900'
                      }`
                    }
                  >
                    {({ isActive }) => (
                      <>
                        <div className="flex items-center gap-2.5 min-w-0">
                          <span
                            className={`shrink-0 transition-colors ${
                              isActive ? 'text-teal-700' : 'text-slate-400 group-hover:text-slate-600'
                            }`}
                          >
                            {item.icon}
                          </span>
                          <span className="truncate">{item.label}</span>
                        </div>
                        {item.badge && (
                          <span
                            className={`text-[9px] font-bold px-1.5 py-0.2 rounded tracking-wide shrink-0 ${
                              isActive
                                ? 'bg-teal-700 text-white'
                                : 'bg-slate-100 text-slate-600 group-hover:bg-slate-200'
                            }`}
                          >
                            {item.badge}
                          </span>
                        )}
                      </>
                    )}
                  </NavLink>
                ))}
              </div>
            </div>
          ))}
        </nav>

        {/* User Card at footer */}
        <div className="p-3 border-t border-slate-100 shrink-0">
          <div className="flex items-center gap-2.5 px-1 py-1">
            <div className="w-7 h-7 rounded-full bg-slate-100 border border-slate-200 text-slate-700 font-semibold text-xs flex items-center justify-center shrink-0">
              {user?.full_name?.charAt(0) || 'U'}
            </div>
            <div className="min-w-0 flex-1">
              <p className="text-xs font-semibold text-slate-800 truncate leading-snug">
                {user?.full_name || user?.username}
              </p>
              <p className="text-[10px] text-slate-400 truncate capitalize">{role}</p>
            </div>
          </div>
        </div>
      </aside>
    </>
  );
};
