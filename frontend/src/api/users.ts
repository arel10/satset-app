import { apiFetch } from './client';
import { User, PaginationMeta, UserRole } from '@/types/api';

export const usersApi = {
  list: async (
    page: number = 1,
    perPage: number = 20
  ): Promise<{ users: User[]; meta?: PaginationMeta }> => {
    const res = await apiFetch<User[]>(`/users?page=${page}&per_page=${perPage}`);
    return {
      users: res.data || [],
      meta: res.meta,
    };
  },

  getByID: async (id: number): Promise<User> => {
    const res = await apiFetch<User>(`/users/${id}`);
    return res.data!;
  },

  create: async (data: {
    username: string;
    email: string;
    password: string;
    full_name: string;
    role: UserRole;
  }): Promise<User> => {
    const res = await apiFetch<User>('/users', {
      method: 'POST',
      body: JSON.stringify(data),
    });
    return res.data!;
  },

  update: async (
    id: number,
    data: {
      email?: string;
      full_name?: string;
      role?: UserRole;
      is_active?: boolean;
      password?: string;
    }
  ): Promise<User> => {
    const res = await apiFetch<User>(`/users/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    });
    return res.data!;
  },

  delete: async (id: number): Promise<void> => {
    await apiFetch(`/users/${id}`, {
      method: 'DELETE',
    });
  },
};
