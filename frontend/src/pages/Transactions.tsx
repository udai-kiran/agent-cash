import React, { useState } from 'react';
import { useTransactions } from '../hooks/useTransactions';
import { TransactionTable } from '../components/transactions/TransactionTable';
import { TransactionFilter } from '../components/transactions/TransactionFilter';
import { LoadingSpinner } from '../components/common/LoadingSpinner';
import { ErrorMessage } from '../components/common/ErrorMessage';
import { Layout } from '../components/common/Layout';
import { TransactionFilter as Filter } from '../types/transaction';

export const TransactionsPage: React.FC = () => {
  const [filter, setFilter] = useState<Filter>({ limit: 50, offset: 0 });
  const { data, isLoading, error } = useTransactions(filter);

  return (
    <Layout>
      <div className="px-4 py-8 sm:px-0">
        <h1 className="text-3xl font-bold text-gray-900 mb-6">Transactions</h1>

        <TransactionFilter onFilterChange={setFilter} />

        {isLoading && <LoadingSpinner />}

        {error && (
          <ErrorMessage
            message={error instanceof Error ? error.message : 'Failed to load transactions'}
          />
        )}

        {data && (
          <>
            <div className="mb-4 text-sm text-gray-600">
              Showing {data.transactions.length} of {data.total} transactions
            </div>
            <TransactionTable transactions={data.transactions} />
          </>
        )}
      </div>
    </Layout>
  );
};
