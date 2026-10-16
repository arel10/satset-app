import { apiFetch } from './client';
import { AuditLog, PaginationMeta } from '@/types/api';

export interface AuditLogFilterParams {
  page?: number;
  per_page?: number;
  user_id?: number;
  action?: string;
  entity_type?: string;
  entity_id?: number;
  start_date?: string;
  end_date?: string;
}

export const auditApi = {
  list: async (
    params: AuditLogFilterParams = {}
  ): Promise<{ logs: AuditLog[]; meta?: PaginationMeta }> => {
    const searchParams = new URLSearchParams();
    if (params.page) searchParams.set('page', params.page.toString());
    if (params.per_page) searchParams.set('per_page', params.per_page.toString());
    if (params.user_id) searchParams.set('user_id', params.user_id.toString());
    if (params.action) searchParams.set('action', params.action);
    if (params.entity_type) searchParams.set('entity_type', params.entity_type);
    if (params.entity_id) searchParams.set('entity_id', params.entity_id.toString());
    if (params.start_date) searchParams.set('start_date', params.start_date);
    if (params.end_date) searchParams.set('end_date', params.end_date);

    const qs = searchParams.toString();
    const res = await apiFetch<AuditLog[]>(`/audit-logs${qs ? `?${qs}` : ''}`);
    return {
      logs: res.data || [],
      meta: res.meta,
    };
  },
};
