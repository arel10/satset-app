import { apiFetch } from './client';
import { LoginResponse, User } from '@/types/api';

export const authApi = {
  login: async (username: string, password: string): Promise<LoginResponse> => {
    const res = await apiFetch<LoginResponse>('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ identifier: username, username, password }),
      skipAuth: true,
    });
    return res.data!;
  },

  getMe: async (): Promise<User> => {
    const res = await apiFetch<User>('/auth/me');
    return res.data!;
  },

  logout: async (refreshToken?: string): Promise<void> => {
    try {
      await apiFetch('/auth/logout', {
        method: 'POST',
        body: JSON.stringify({ refresh_token: refreshToken }),
      });
    } catch {
      // Ignore error during logout
    }
  },
};
