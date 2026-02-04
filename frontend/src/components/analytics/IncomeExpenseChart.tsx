import React from 'react';
import {
  AreaChart,
  Area,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from 'recharts';
import { IncomeExpenseResponse } from '../../types/analytics';
import { formatCurrency } from '../../utils/currency';

interface IncomeExpenseChartProps {
  data: IncomeExpenseResponse;
}

export const IncomeExpenseChart: React.FC<IncomeExpenseChartProps> = ({ data }) => {
  const chartData = data.data.map((item) => ({
    period: item.period,
    income: parseFloat(item.income),
    expense: parseFloat(item.expense),
    net: parseFloat(item.net),
  }));

  return (
    <div className="bg-white shadow rounded-lg p-6">
      <h3 className="text-lg font-medium text-gray-900 mb-4">Income vs Expense</h3>

      <div className="grid grid-cols-3 gap-4 mb-6">
        <div className="text-center">
          <div className="text-sm text-gray-500">Total Income</div>
          <div className="text-2xl font-bold text-green-600">{formatCurrency(data.total_income, data.currency_mnemonic)}</div>
        </div>
        <div className="text-center">
          <div className="text-sm text-gray-500">Total Expense</div>
          <div className="text-2xl font-bold text-red-600">{formatCurrency(data.total_expense, data.currency_mnemonic)}</div>
        </div>
        <div className="text-center">
          <div className="text-sm text-gray-500">Net</div>
          <div className={`text-2xl font-bold ${
            parseFloat(data.net_total) >= 0 ? 'text-green-600' : 'text-red-600'
          }`}>
            {formatCurrency(data.net_total, data.currency_mnemonic)}
          </div>
        </div>
      </div>

      <ResponsiveContainer width="100%" height={300}>
        <AreaChart data={chartData}>
          <CartesianGrid strokeDasharray="3 3" />
          <XAxis dataKey="period" />
          <YAxis />
          <Tooltip formatter={(value) => formatCurrency(value as number, data.currency_mnemonic)} />
          <Legend />
          <Area
            type="monotone"
            dataKey="income"
            stackId="1"
            stroke="#10b981"
            fill="#10b981"
            name="Income"
          />
          <Area
            type="monotone"
            dataKey="expense"
            stackId="2"
            stroke="#ef4444"
            fill="#ef4444"
            name="Expense"
          />
        </AreaChart>
      </ResponsiveContainer>
    </div>
  );
};
