import React from 'react';
import { PieChart, Pie, Cell, ResponsiveContainer, Legend, Tooltip } from 'recharts';
import { CategoryBreakdownItem } from '../../types/analytics';
import { formatCurrency } from '../../utils/currency';

interface CategoryPieChartProps {
  title: string;
  data: CategoryBreakdownItem[];
  currencyMnemonic?: string;
}

const COLORS = [
  '#3b82f6', '#10b981', '#f59e0b', '#ef4444', '#8b5cf6',
  '#ec4899', '#14b8a6', '#f97316', '#06b6d4', '#84cc16',
];

export const CategoryPieChart: React.FC<CategoryPieChartProps> = ({ title, data, currencyMnemonic }) => {
  const chartData = data.map((item) => ({
    name: item.category,
    value: parseFloat(item.amount),
    count: item.count,
  }));

  // Sort by value descending and take top 10
  const topData = chartData.sort((a, b) => b.value - a.value).slice(0, 10);

  return (
    <div className="bg-white shadow rounded-lg p-6">
      <h3 className="text-lg font-medium text-gray-900 mb-4">{title}</h3>

      {topData.length > 0 ? (
        <>
          <ResponsiveContainer width="100%" height={300}>
            <PieChart>
              <Pie
                data={topData}
                cx="50%"
                cy="50%"
                labelLine={false}
                label={({ name, percent }: { name?: string; percent?: number }) => `${name ?? ''}: ${((percent ?? 0) * 100).toFixed(0)}%`}
                outerRadius={80}
                fill="#8884d8"
                dataKey="value"
              >
                {topData.map((entry, index) => (
                  <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                ))}
              </Pie>
              <Tooltip formatter={(value) => formatCurrency(value as number, currencyMnemonic)} />
            </PieChart>
          </ResponsiveContainer>

          <div className="mt-4 space-y-2">
            {topData.slice(0, 5).map((item, index) => (
              <div key={index} className="flex justify-between items-center text-sm">
                <div className="flex items-center">
                  <div
                    className="w-3 h-3 rounded-full mr-2"
                    style={{ backgroundColor: COLORS[index % COLORS.length] }}
                  ></div>
                  <span className="text-gray-700">{item.name}</span>
                </div>
                <span className="font-medium">{formatCurrency(item.value, currencyMnemonic)}</span>
              </div>
            ))}
          </div>
        </>
      ) : (
        <div className="text-center py-8 text-gray-500">
          No data available
        </div>
      )}
    </div>
  );
};
