import { apiFetch } from './client';
import { BusinessReportData } from '@/types/api';

export const reportsApi = {
  getBusinessReport: async (
    startDate?: string,
    endDate?: string
  ): Promise<BusinessReportData> => {
    const searchParams = new URLSearchParams();
    if (startDate) searchParams.set('start_date', startDate);
    if (endDate) searchParams.set('end_date', endDate);

    const qs = searchParams.toString();
    const res = await apiFetch<BusinessReportData>(`/reports/business${qs ? `?${qs}` : ''}`);
    return res.data!;
  },
};
