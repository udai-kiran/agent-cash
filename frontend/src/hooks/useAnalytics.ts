import { useQuery } from '@tanstack/react-query';
import { analyticsApi } from '../api/analytics';

export const useIncomeExpense = (startDate?: string, endDate?: string) => {
  return useQuery({
    queryKey: ['analytics', 'income-expense', startDate, endDate],
    queryFn: () => analyticsApi.getIncomeExpense(startDate, endDate),
  });
};

export const useCategoryBreakdown = (startDate?: string, endDate?: string) => {
  return useQuery({
    queryKey: ['analytics', 'category-breakdown', startDate, endDate],
    queryFn: () => analyticsApi.getCategoryBreakdown(startDate, endDate),
  });
};

export const useNetWorth = () => {
  return useQuery({
    queryKey: ['analytics', 'net-worth'],
    queryFn: analyticsApi.getNetWorth,
  });
};
