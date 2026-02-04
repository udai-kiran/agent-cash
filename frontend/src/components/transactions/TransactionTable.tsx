import React from 'react';
import { Transaction } from '../../types/transaction';
import { formatCurrency } from '../../utils/currency';

interface TransactionTableProps {
  transactions: Transaction[];
}

export const TransactionTable: React.FC<TransactionTableProps> = ({ transactions }) => {
  const formatDate = (dateStr: string) => {
    return new Date(dateStr).toLocaleDateString();
  };

  return (
    <div className="bg-white shadow overflow-hidden sm:rounded-lg">
      <table className="min-w-full divide-y divide-gray-200">
        <thead className="bg-gray-50">
          <tr>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Date
            </th>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Description
            </th>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Account
            </th>
            <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
              Amount
            </th>
          </tr>
        </thead>
        <tbody className="bg-white divide-y divide-gray-200">
          {transactions.map((transaction) => {
            // For display purposes, show the first two splits
            const primarySplits = transaction.splits.slice(0, 2);

            return primarySplits.map((split, index) => (
              <tr key={`${transaction.guid}-${split.guid}`} className="hover:bg-gray-50">
                {index === 0 && (
                  <>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900" rowSpan={primarySplits.length}>
                      {formatDate(transaction.post_date)}
                    </td>
                    <td className="px-6 py-4 text-sm text-gray-900" rowSpan={primarySplits.length}>
                      <div>
                        {transaction.description || '(No description)'}
                        {transaction.num && (
                          <span className="ml-2 text-xs text-gray-500">#{transaction.num}</span>
                        )}
                      </div>
                      {split.memo && (
                        <div className="text-xs text-gray-500 mt-1">{split.memo}</div>
                      )}
                    </td>
                  </>
                )}
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {split.account?.name || split.account_guid}
                  <div className="text-xs text-gray-500">{split.account?.type}</div>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-right">
                  <span className={`font-medium ${
                    parseFloat(split.value) >= 0 ? 'text-green-600' : 'text-red-600'
                  }`}>
                    {formatCurrency(split.value, transaction.currency_mnemonic)}
                  </span>
                </td>
              </tr>
            ));
          })}
        </tbody>
      </table>
      {transactions.length === 0 && (
        <div className="text-center py-8 text-gray-500">
          No transactions found
        </div>
      )}
    </div>
  );
};
