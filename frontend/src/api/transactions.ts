import { apiClient } from './client';
import { Transaction, TransactionFilter } from '../types/transaction';

interface TransactionListResponse {
  transactions: Transaction[];
  total: number;
  limit: number;
  offset: number;
}

export const transactionsApi = {
  getAll: async (filter?: TransactionFilter): Promise<TransactionListResponse> => {
    const params = new URLSearchParams();

    if (filter) {
      if (filter.account_guid) params.append('account_guid', filter.account_guid);
      if (filter.start_date) params.append('start_date', filter.start_date);
      if (filter.end_date) params.append('end_date', filter.end_date);
      if (filter.description) params.append('description', filter.description);
      if (filter.limit) params.append('limit', filter.limit.toString());
      if (filter.offset) params.append('offset', filter.offset.toString());
    }

    const response = await apiClient.get<TransactionListResponse>(
      `/transactions?${params.toString()}`
    );
    return response.data;
  },

  getById: async (guid: string): Promise<Transaction> => {
    const response = await apiClient.get<Transaction>(`/transactions/${guid}`);
    return response.data;
  },
};
