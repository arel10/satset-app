import { apiFetch } from './client';
import { OwnerDashboardData, KasirDashboardData } from '@/types/api';

export const dashboardApi = {
  getOwnerDashboard: async (): Promise<OwnerDashboardData> => {
    const res = await apiFetch<OwnerDashboardData>('/dashboard/owner');
    return res.data!;
  },

  getKasirDashboard: async (): Promise<KasirDashboardData> => {
    const res = await apiFetch<KasirDashboardData>('/dashboard/kasir');
    return res.data!;
  },
};
