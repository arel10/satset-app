import { apiFetch } from './client';
import { Expense, ExpenseCategory, PaginationMeta } from '@/types/api';

export interface ExpenseFilterParams {
  page?: number;
  per_page?: number;
  category_id?: number;
  start_date?: string;
  end_date?: string;
}

export const expensesApi = {
  list: async (
    params: ExpenseFilterParams = {}
  ): Promise<{ expenses: Expense[]; meta?: PaginationMeta }> => {
    const searchParams = new URLSearchParams();
    if (params.page) searchParams.set('page', params.page.toString());
    if (params.per_page) searchParams.set('per_page', params.per_page.toString());
    if (params.category_id) searchParams.set('category_id', params.category_id.toString());
    if (params.start_date) searchParams.set('start_date', params.start_date);
    if (params.end_date) searchParams.set('end_date', params.end_date);

    const qs = searchParams.toString();
    const res = await apiFetch<Expense[]>(`/expenses${qs ? `?${qs}` : ''}`);
    return {
      expenses: res.data || [],
      meta: res.meta,
    };
  },

  create: async (data: {
    category_id: number;
    amount: number;
    expense_date: string;
    description: string;
    vendor?: string;
    receipt_url?: string;
  }): Promise<Expense> => {
    const res = await apiFetch<Expense>('/expenses', {
      method: 'POST',
      body: JSON.stringify(data),
    });
    return res.data!;
  },

  cancel: async (id: number, reason: string): Promise<Expense> => {
    const res = await apiFetch<Expense>(`/expenses/${id}/cancel`, {
      method: 'POST',
      body: JSON.stringify({ reason }),
    });
    return res.data!;
  },

  listCategories: async (): Promise<ExpenseCategory[]> => {
    const res = await apiFetch<ExpenseCategory[]>('/expenses/categories');
    return res.data || [];
  },

  createCategory: async (data: {
    name: string;
    description?: string;
  }): Promise<ExpenseCategory> => {
    const res = await apiFetch<ExpenseCategory>('/expenses/categories', {
      method: 'POST',
      body: JSON.stringify(data),
    });
    return res.data!;
  },

  updateCategory: async (
    id: number,
    data: {
      name?: string;
      description?: string;
      is_active?: boolean;
    }
  ): Promise<ExpenseCategory> => {
    const res = await apiFetch<ExpenseCategory>(`/expenses/categories/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    });
    return res.data!;
  },
};
