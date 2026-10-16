import { apiFetch } from './client';
import { Payment } from '@/types/api';

export const paymentsApi = {
  record: async (
    orderID: number,
    amount: number,
    paymentMethod: string,
    notes?: string
  ): Promise<Payment> => {
    const res = await apiFetch<Payment>('/payments', {
      method: 'POST',
      body: JSON.stringify({
        order_id: orderID,
        amount,
        payment_method: paymentMethod,
        notes: notes || '',
      }),
    });
    return res.data!;
  },

  cancel: async (paymentID: number, reason: string): Promise<Payment> => {
    const res = await apiFetch<Payment>(`/payments/${paymentID}/cancel`, {
      method: 'PUT',
      body: JSON.stringify({ reason }),
    });
    return res.data!;
  },
};
