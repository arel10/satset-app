import { apiFetch } from './client';
import { LaundryService, PaginationMeta } from '@/types/api';

export const servicesApi = {
  getActive: async (): Promise<LaundryService[]> => {
    const res = await apiFetch<LaundryService[]>('/services/active');
    return res.data || [];
  },

  list: async (
    page: number = 1,
    perPage: number = 50
  ): Promise<{ services: LaundryService[]; meta?: PaginationMeta }> => {
    const res = await apiFetch<LaundryService[]>(`/services?page=${page}&per_page=${perPage}`);
    return {
      services: res.data || [],
      meta: res.meta,
    };
  },

  getByID: async (id: number): Promise<LaundryService> => {
    const res = await apiFetch<LaundryService>(`/services/${id}`);
    return res.data!;
  },

  create: async (data: {
    name: string;
    category: string;
    unit: string;
    price_per_unit: number;
    estimated_duration: number;
    is_active?: boolean;
  }): Promise<LaundryService> => {
    const res = await apiFetch<LaundryService>('/services', {
      method: 'POST',
      body: JSON.stringify(data),
    });
    return res.data!;
  },

  update: async (
    id: number,
    data: {
      name?: string;
      category?: string;
      unit?: string;
      price_per_unit?: number;
      estimated_duration?: number;
      is_active?: boolean;
    }
  ): Promise<LaundryService> => {
    const res = await apiFetch<LaundryService>(`/services/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    });
    return res.data!;
  },
};
