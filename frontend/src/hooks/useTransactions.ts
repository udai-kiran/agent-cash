import { useQuery } from '@tanstack/react-query';
import { transactionsApi } from '../api/transactions';
import { TransactionFilter } from '../types/transaction';

export const useTransactions = (filter?: TransactionFilter) => {
  return useQuery({
    queryKey: ['transactions', filter],
    queryFn: () => transactionsApi.getAll(filter),
  });
};

export const useTransaction = (guid: string) => {
  return useQuery({
    queryKey: ['transactions', guid],
    queryFn: () => transactionsApi.getById(guid),
    enabled: !!guid,
  });
};
