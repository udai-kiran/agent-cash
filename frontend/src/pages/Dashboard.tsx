import React from 'react';
import { Layout } from '../components/common/Layout';
import { LoadingSpinner } from '../components/common/LoadingSpinner';
import { ErrorMessage } from '../components/common/ErrorMessage';
import { IncomeExpenseChart } from '../components/analytics/IncomeExpenseChart';
import { NetWorthSummary } from '../components/analytics/NetWorthSummary';
import { useIncomeExpense, useNetWorth } from '../hooks/useAnalytics';

export const DashboardPage: React.FC = () => {
  const { data: incomeExpenseData, isLoading: ieLoading, error: ieError } = useIncomeExpense();
  const { data: netWorthData, isLoading: nwLoading, error: nwError } = useNetWorth();

  const isLoading = ieLoading || nwLoading;
  const error = ieError || nwError;

  return (
    <Layout>
      <div className="px-4 py-8 sm:px-0">
        <h1 className="text-3xl font-bold text-gray-900 mb-6">Dashboard</h1>

        {isLoading && <LoadingSpinner />}

        {error && (
          <ErrorMessage
            message={error instanceof Error ? error.message : 'Failed to load analytics'}
          />
        )}

        {!isLoading && !error && (
          <div className="space-y-6">
            {netWorthData && <NetWorthSummary data={netWorthData} />}
            {incomeExpenseData && <IncomeExpenseChart data={incomeExpenseData} />}
          </div>
        )}
      </div>
    </Layout>
  );
};
