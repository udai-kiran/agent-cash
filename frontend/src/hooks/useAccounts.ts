import { useQuery } from '@tanstack/react-query';
import { accountsApi } from '../api/accounts';

export const useAccounts = () => {
  return useQuery({
    queryKey: ['accounts'],
    queryFn: accountsApi.getAll,
  });
};

export const useAccountHierarchy = () => {
  return useQuery({
    queryKey: ['accounts', 'hierarchy'],
    queryFn: accountsApi.getHierarchy,
  });
};

export const useAccount = (guid: string) => {
  return useQuery({
    queryKey: ['accounts', guid],
    queryFn: () => accountsApi.getById(guid),
    enabled: !!guid,
  });
};

export const useAccountBalance = (guid: string) => {
  return useQuery({
    queryKey: ['accounts', guid, 'balance'],
    queryFn: () => accountsApi.getBalance(guid),
    enabled: !!guid,
  });
};
