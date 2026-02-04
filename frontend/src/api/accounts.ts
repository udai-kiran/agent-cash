import { apiClient } from './client';
import { Account, AccountBalance } from '../types/account';

export const accountsApi = {
  getAll: async (): Promise<Account[]> => {
    const response = await apiClient.get<Account[]>('/accounts');
    return response.data;
  },

  getHierarchy: async (): Promise<Account[]> => {
    const response = await apiClient.get<Account[]>('/accounts/hierarchy');
    return response.data;
  },

  getById: async (guid: string): Promise<Account> => {
    const response = await apiClient.get<Account>(`/accounts/${guid}`);
    return response.data;
  },

  getBalance: async (guid: string): Promise<AccountBalance> => {
    const response = await apiClient.get<AccountBalance>(`/accounts/${guid}/balance`);
    return response.data;
  },
};
