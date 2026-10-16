import { apiFetch } from './client';

export interface BusinessSettings {
  business_name?: string;
  business_address?: string;
  business_phone?: string;
  receipt_footer?: string;
  business_logo_url?: string;
}

export const settingsApi = {
  getAll: async (): Promise<Record<string, string>> => {
    const res = await apiFetch<Record<string, string>>('/settings');
    return res.data || {};
  },

  update: async (data: BusinessSettings): Promise<void> => {
    await apiFetch('/settings', {
      method: 'PUT',
      body: JSON.stringify(data),
    });
  },

  uploadLogo: async (file: File): Promise<{ logo_url: string }> => {
    const formData = new FormData();
    formData.append('logo', file);

    const res = await apiFetch<{ logo_url: string }>('/settings/logo', {
      method: 'POST',
      body: formData,
    });
    return res.data!;
  },
};
