import React, { useState } from 'react';
import { Layout } from '../components/common/Layout';
import { LoadingSpinner } from '../components/common/LoadingSpinner';
import { ErrorMessage } from '../components/common/ErrorMessage';
import { IncomeExpenseChart } from '../components/analytics/IncomeExpenseChart';
import { CategoryPieChart } from '../components/analytics/CategoryPieChart';
import { NetWorthSummary } from '../components/analytics/NetWorthSummary';
import {
  useIncomeExpense,
  useCategoryBreakdown,
  useNetWorth,
} from '../hooks/useAnalytics';

export const AnalyticsPage: React.FC = () => {
  const [dateRange, setDateRange] = useState({
    startDate: '',
    endDate: '',
  });

  const { data: incomeExpenseData, isLoading: ieLoading, error: ieError } = useIncomeExpense(
    dateRange.startDate,
    dateRange.endDate
  );

  const { data: categoryData, isLoading: catLoading, error: catError } = useCategoryBreakdown(
    dateRange.startDate,
    dateRange.endDate
  );

  const { data: netWorthData, isLoading: nwLoading, error: nwError } = useNetWorth();

  const isLoading = ieLoading || catLoading || nwLoading;
  const error = ieError || catError || nwError;

  const handleDateRangeSubmit = (e: React.FormEvent) => {
    e.preventDefault();
  };

  return (
    <Layout>
      <div className="px-4 py-8 sm:px-0">
        <h1 className="text-3xl font-bold text-gray-900 mb-6">Analytics</h1>

        <div className="bg-white shadow rounded-lg p-4 mb-6">
          <form onSubmit={handleDateRangeSubmit} className="flex items-end space-x-4">
            <div>
              <label htmlFor="startDate" className="block text-sm font-medium text-gray-700">
                Start Date
              </label>
              <input
                type="date"
                id="startDate"
                value={dateRange.startDate}
                onChange={(e) => setDateRange({ ...dateRange, startDate: e.target.value })}
                className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
              />
            </div>

            <div>
              <label htmlFor="endDate" className="block text-sm font-medium text-gray-700">
                End Date
              </label>
              <input
                type="date"
                id="endDate"
                value={dateRange.endDate}
                onChange={(e) => setDateRange({ ...dateRange, endDate: e.target.value })}
                className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
              />
            </div>

            <button
              type="button"
              onClick={() => setDateRange({ startDate: '', endDate: '' })}
              className="px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50"
            >
              Reset
            </button>
          </form>
        </div>

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

            {categoryData && (
              <div className="grid md:grid-cols-2 gap-6">
                <CategoryPieChart title="Income by Category" data={categoryData.income} currencyMnemonic={categoryData.currency_mnemonic} />
                <CategoryPieChart title="Expenses by Category" data={categoryData.expense} currencyMnemonic={categoryData.currency_mnemonic} />
              </div>
            )}
          </div>
        )}
      </div>
    </Layout>
  );
};
