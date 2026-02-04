import { apiClient } from './client';
import {
  IncomeExpenseResponse,
  CategoryBreakdownResponse,
  NetWorthResponse,
} from '../types/analytics';

export const analyticsApi = {
  getIncomeExpense: async (startDate?: string, endDate?: string): Promise<IncomeExpenseResponse> => {
    const params = new URLSearchParams();
    if (startDate) params.append('start_date', startDate);
    if (endDate) params.append('end_date', endDate);

    const response = await apiClient.get<IncomeExpenseResponse>(
      `/analytics/income-expense?${params.toString()}`
    );
    return response.data;
  },

  getCategoryBreakdown: async (startDate?: string, endDate?: string): Promise<CategoryBreakdownResponse> => {
    const params = new URLSearchParams();
    if (startDate) params.append('start_date', startDate);
    if (endDate) params.append('end_date', endDate);

    const response = await apiClient.get<CategoryBreakdownResponse>(
      `/analytics/category-breakdown?${params.toString()}`
    );
    return response.data;
  },

  getNetWorth: async (): Promise<NetWorthResponse> => {
    const response = await apiClient.get<NetWorthResponse>('/analytics/net-worth');
    return response.data;
  },
};
