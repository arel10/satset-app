import { ApiResponse } from '@/types/api';

const TOKEN_KEY = 'satset_access_token';
const REFRESH_TOKEN_KEY = 'satset_refresh_token';

export function getStoredAccessToken(): string | null {
  return localStorage.getItem(TOKEN_KEY);
}

export function getStoredRefreshToken(): string | null {
  return localStorage.getItem(REFRESH_TOKEN_KEY);
}

export function setStoredTokens(accessToken: string, refreshToken?: string): void {
  localStorage.setItem(TOKEN_KEY, accessToken);
  if (refreshToken) {
    localStorage.setItem(REFRESH_TOKEN_KEY, refreshToken);
  }
}

export function clearStoredTokens(): void {
  localStorage.removeItem(TOKEN_KEY);
  localStorage.removeItem(REFRESH_TOKEN_KEY);
}

interface RequestOptions extends RequestInit {
  idempotencyKey?: string;
  skipAuth?: boolean;
}

let isRefreshing = false;
let failedQueue: Array<{
  resolve: (value: any) => void;
  reject: (reason?: any) => void;
}> = [];

function processQueue(error: any, token: string | null = null) {
  failedQueue.forEach(prom => {
    if (error) {
      prom.reject(error);
    } else {
      prom.resolve(token);
    }
  });
  failedQueue = [];
}

const API_BASE_URL = (import.meta.env.VITE_API_BASE_URL || '/api').replace(/\/+$/, '');

export function getApiBaseUrl(): string {
  return API_BASE_URL;
}

/**
 * Base HTTP fetch wrapper with auto-retry on 401 using refresh token.
 */
export async function apiFetch<T = any>(
  endpoint: string,
  options: RequestOptions = {}
): Promise<ApiResponse<T>> {
  const normalizedEndpoint = endpoint.startsWith('/') ? endpoint : `/${endpoint}`;
  const url = endpoint.startsWith('http') ? endpoint : `${API_BASE_URL}${normalizedEndpoint}`;
  const headers = new Headers(options.headers || {});

  if (!headers.has('Content-Type') && !(options.body instanceof FormData)) {
    headers.set('Content-Type', 'application/json');
  }

  // Attach access token
  if (!options.skipAuth) {
    const token = getStoredAccessToken();
    if (token) {
      headers.set('Authorization', `Bearer ${token}`);
    }
  }

  // Attach idempotency key if provided
  if (options.idempotencyKey) {
    headers.set('Idempotency-Key', options.idempotencyKey);
  }

  let response = await fetch(url, {
    ...options,
    headers,
  });

  // Handle 401 Unauthorized with token refresh
  if (response.status === 401 && !options.skipAuth && !endpoint.includes('/auth/login') && !endpoint.includes('/auth/refresh')) {
    const refreshToken = getStoredRefreshToken();

    if (!refreshToken) {
      clearStoredTokens();
      window.dispatchEvent(new Event('auth:unauthorized'));
      throw new Error('Sesi telah berakhir. Silakan login kembali.');
    }

    if (isRefreshing) {
      return new Promise((resolve, reject) => {
        failedQueue.push({ resolve, reject });
      })
        .then(token => {
          headers.set('Authorization', `Bearer ${token}`);
          return fetch(url, { ...options, headers }).then(res => res.json());
        });
    }

    isRefreshing = true;

    try {
      const refreshRes = await fetch(`${API_BASE_URL}/auth/refresh`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ refresh_token: refreshToken }),
      });

      const refreshData: ApiResponse<{ access_token: string }> = await refreshRes.json();

      if (!refreshRes.ok || !refreshData.success || !refreshData.data?.access_token) {
        throw new Error('Gagal memperbarui token');
      }

      const newAccessToken = refreshData.data.access_token;
      setStoredTokens(newAccessToken);
      processQueue(null, newAccessToken);

      // Retry original request
      headers.set('Authorization', `Bearer ${newAccessToken}`);
      response = await fetch(url, { ...options, headers });
    } catch (refreshErr) {
      processQueue(refreshErr, null);
      clearStoredTokens();
      window.dispatchEvent(new Event('auth:unauthorized'));
      throw new Error('Sesi login telah habis. Silakan masuk kembali.');
    } finally {
      isRefreshing = false;
    }
  }

  const data: ApiResponse<T> = await response.json().catch(() => ({
    success: false,
    error: `Server mengembalikan status HTTP ${response.status}`,
  }));

  if (!response.ok || !data.success) {
    const errorMsg = data.error || data.message || `Terjadi kesalahan (HTTP ${response.status})`;
    throw new Error(errorMsg);
  }

  return data;
}
