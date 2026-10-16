import { apiFetch } from './client';
import { Customer, Order, PaginationMeta } from '@/types/api';

export const customersApi = {
  search: async (
    q: string = '',
    page: number = 1,
    perPage: number = 20
  ): Promise<{ customers: Customer[]; meta?: PaginationMeta }> => {
    const searchParams = new URLSearchParams({
      page: page.toString(),
      per_page: perPage.toString(),
    });
    if (q) searchParams.set('q', q);

    const res = await apiFetch<Customer[]>(`/customers?${searchParams.toString()}`);
    return {
      customers: res.data || [],
      meta: res.meta,
    };
  },

  getByID: async (id: number): Promise<Customer> => {
    const res = await apiFetch<Customer>(`/customers/${id}`);
    return res.data!;
  },

  create: async (data: {
    name: string;
    phone: string;
    address?: string;
    notes?: string;
  }): Promise<Customer> => {
    const res = await apiFetch<Customer>('/customers', {
      method: 'POST',
      body: JSON.stringify(data),
    });
    return res.data!;
  },

  update: async (
    id: number,
    data: {
      name?: string;
      phone?: string;
      address?: string;
      notes?: string;
    }
  ): Promise<Customer> => {
    const res = await apiFetch<Customer>(`/customers/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    });
    return res.data!;
  },

  getOrderHistory: async (id: number, limit: number = 10): Promise<Order[]> => {
    const res = await apiFetch<Order[]>(`/customers/${id}/orders?limit=${limit}`);
    return res.data || [];
  },
};
