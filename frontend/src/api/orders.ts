import { apiFetch } from './client';
import { CreateOrderInput, Order, Payment, PaginationMeta } from '@/types/api';

export interface OrderFilterParams {
  page?: number;
  per_page?: number;
  q?: string;
  status?: string;
  payment_status?: string;
  start_date?: string;
  end_date?: string;
  overdue?: boolean;
  customer_id?: number;
}

export const ordersApi = {
  create: async (data: CreateOrderInput, idempotencyKey: string): Promise<Order> => {
    const res = await apiFetch<Order>('/orders', {
      method: 'POST',
      body: JSON.stringify(data),
      idempotencyKey,
    });
    return res.data!;
  },

  list: async (params: OrderFilterParams = {}): Promise<{ orders: Order[]; meta?: PaginationMeta }> => {
    const searchParams = new URLSearchParams();
    if (params.page) searchParams.set('page', params.page.toString());
    if (params.per_page) searchParams.set('per_page', params.per_page.toString());
    if (params.q) searchParams.set('q', params.q);
    if (params.status) searchParams.set('status', params.status);
    if (params.payment_status) searchParams.set('payment_status', params.payment_status);
    if (params.start_date) searchParams.set('start_date', params.start_date);
    if (params.end_date) searchParams.set('end_date', params.end_date);
    if (params.overdue) searchParams.set('overdue', 'true');
    if (params.customer_id) searchParams.set('customer_id', params.customer_id.toString());

    const qs = searchParams.toString();
    const res = await apiFetch<Order[]>(`/orders${qs ? `?${qs}` : ''}`);
    return {
      orders: res.data || [],
      meta: res.meta,
    };
  },

  getByID: async (id: number): Promise<Order> => {
    const res = await apiFetch<Order>(`/orders/${id}`);
    return res.data!;
  },

  getByNumber: async (orderNumber: string): Promise<Order> => {
    const res = await apiFetch<Order>(`/orders/number/${orderNumber}`);
    return res.data!;
  },

  updateStatus: async (id: number, status: string, notes?: string): Promise<Order> => {
    const res = await apiFetch<Order>(`/orders/${id}/status`, {
      method: 'PATCH',
      body: JSON.stringify({ status, notes }),
    });
    return res.data!;
  },

  cancel: async (id: number, reason: string): Promise<Order> => {
    const res = await apiFetch<Order>(`/orders/${id}/cancel`, {
      method: 'POST',
      body: JSON.stringify({ reason }),
    });
    return res.data!;
  },

  getPayments: async (orderID: number): Promise<Payment[]> => {
    const res = await apiFetch<Payment[]>(`/orders/${orderID}/payments`);
    return res.data || [];
  },
};
