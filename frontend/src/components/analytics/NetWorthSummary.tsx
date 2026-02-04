import React from 'react';
import { NetWorthResponse } from '../../types/analytics';
import { formatCurrency } from '../../utils/currency';

interface NetWorthSummaryProps {
  data: NetWorthResponse;
}

export const NetWorthSummary: React.FC<NetWorthSummaryProps> = ({ data }) => {
  return (
    <div className="bg-white shadow rounded-lg p-6">
      <h3 className="text-lg font-medium text-gray-900 mb-4">Net Worth</h3>

      <div className="grid grid-cols-3 gap-4 mb-6">
        <div className="text-center p-4 bg-green-50 rounded-lg">
          <div className="text-sm text-gray-600">Total Assets</div>
          <div className="text-2xl font-bold text-green-600">{formatCurrency(data.total_assets, data.currency_mnemonic)}</div>
        </div>
        <div className="text-center p-4 bg-red-50 rounded-lg">
          <div className="text-sm text-gray-600">Total Liabilities</div>
          <div className="text-2xl font-bold text-red-600">{formatCurrency(data.total_liabilities, data.currency_mnemonic)}</div>
        </div>
        <div className="text-center p-4 bg-blue-50 rounded-lg">
          <div className="text-sm text-gray-600">Net Worth</div>
          <div className={`text-2xl font-bold ${
            parseFloat(data.net_worth) >= 0 ? 'text-blue-600' : 'text-red-600'
          }`}>
            {formatCurrency(data.net_worth, data.currency_mnemonic)}
          </div>
        </div>
      </div>

      <div className="grid md:grid-cols-2 gap-6">
        <div>
          <h4 className="text-sm font-medium text-gray-700 mb-3">Assets</h4>
          <div className="space-y-2">
            {data.assets.length > 0 ? (
              data.assets.slice(0, 5).map((item, index) => (
                <div key={index} className="flex justify-between text-sm">
                  <span className="text-gray-600 truncate">{item.account_name}</span>
                  <span className="font-medium text-green-600">{formatCurrency(item.balance, data.currency_mnemonic)}</span>
                </div>
              ))
            ) : (
              <div className="text-sm text-gray-500">No assets</div>
            )}
          </div>
        </div>

        <div>
          <h4 className="text-sm font-medium text-gray-700 mb-3">Liabilities</h4>
          <div className="space-y-2">
            {data.liabilities.length > 0 ? (
              data.liabilities.slice(0, 5).map((item, index) => (
                <div key={index} className="flex justify-between text-sm">
                  <span className="text-gray-600 truncate">{item.account_name}</span>
                  <span className="font-medium text-red-600">{formatCurrency(item.balance, data.currency_mnemonic)}</span>
                </div>
              ))
            ) : (
              <div className="text-sm text-gray-500">No liabilities</div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};
