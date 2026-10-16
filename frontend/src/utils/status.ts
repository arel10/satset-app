import { WorkStatus, PaymentStatus } from '@/types/api';

export interface StatusMeta {
  label: string;
  badgeClass: string;
}

export const WORK_STATUS_META: Record<WorkStatus, StatusMeta> = {
  DITERIMA: {
    label: 'Diterima',
    badgeClass: 'bg-blue-50 text-blue-700 border-blue-200',
  },
  DICUCI: {
    label: 'Dicuci',
    badgeClass: 'bg-cyan-50 text-cyan-700 border-cyan-200',
  },
  DIKERINGKAN: {
    label: 'Dikeringkan',
    badgeClass: 'bg-amber-50 text-amber-700 border-amber-200',
  },
  DISETRIKA: {
    label: 'Disetrika',
    badgeClass: 'bg-indigo-50 text-indigo-700 border-indigo-200',
  },
  DIKEMAS: {
    label: 'Dikemas',
    badgeClass: 'bg-purple-50 text-purple-700 border-purple-200',
  },
  SIAP_DIAMBIL: {
    label: 'Siap Diambil',
    badgeClass: 'bg-teal-50 text-teal-700 border-teal-300 font-medium',
  },
  SELESAI: {
    label: 'Selesai',
    badgeClass: 'bg-emerald-50 text-emerald-700 border-emerald-200',
  },
  DIBATALKAN: {
    label: 'Dibatalkan',
    badgeClass: 'bg-rose-50 text-rose-700 border-rose-200',
  },
};

export const PAYMENT_STATUS_META: Record<PaymentStatus, StatusMeta> = {
  LUNAS: {
    label: 'LUNAS',
    badgeClass: 'bg-emerald-50 text-emerald-700 border-emerald-300 font-semibold',
  },
  BELUM_LUNAS: {
    label: 'BELUM LUNAS',
    badgeClass: 'bg-amber-50 text-amber-700 border-amber-300 font-semibold',
  },
};

/**
 * Returns allowed normal forward transitions for Operator/Kasir
 */
export function getNextWorkStatuses(current: WorkStatus): WorkStatus[] {
  switch (current) {
    case 'DITERIMA':
      return ['DICUCI', 'DISETRIKA'];
    case 'DICUCI':
      return ['DIKERINGKAN', 'DIKEMAS'];
    case 'DIKERINGKAN':
      return ['DISETRIKA', 'DIKEMAS'];
    case 'DISETRIKA':
      return ['DIKEMAS'];
    case 'DIKEMAS':
      return ['SIAP_DIAMBIL'];
    case 'SIAP_DIAMBIL':
      return ['SELESAI'];
    default:
      return [];
  }
}
