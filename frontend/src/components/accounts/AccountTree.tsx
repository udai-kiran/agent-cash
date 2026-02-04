import React, { useState } from 'react';
import { Account } from '../../types/account';
import { formatCurrency } from '../../utils/currency';

interface AccountTreeProps {
  accounts: Account[];
}

interface AccountNodeProps {
  account: Account;
  level: number;
}

const AccountNode: React.FC<AccountNodeProps> = ({ account, level }) => {
  const [isExpanded, setIsExpanded] = useState(true);
  const hasChildren = account.children && account.children.length > 0;

  const indentClass = `pl-${level * 4}`;

  return (
    <div>
      <div
        className={`flex items-center py-2 px-4 hover:bg-gray-50 cursor-pointer ${indentClass}`}
        onClick={() => hasChildren && setIsExpanded(!isExpanded)}
      >
        <div className="flex-1 flex items-center">
          {hasChildren && (
            <span className="mr-2 text-gray-500">
              {isExpanded ? '▼' : '▶'}
            </span>
          )}
          {!hasChildren && <span className="mr-2 w-4"></span>}
          <div className="flex-1">
            <span className="font-medium text-gray-900">{account.name}</span>
            {account.code && (
              <span className="ml-2 text-sm text-gray-500">({account.code})</span>
            )}
            <div className="text-xs text-gray-500">{account.type}</div>
          </div>
        </div>
        <div className="text-right">
          <span className={`font-medium ${
            parseFloat(account.balance) >= 0 ? 'text-green-600' : 'text-red-600'
          }`}>
            {formatCurrency(account.balance, account.commodity_mnemonic)}
          </span>
        </div>
      </div>
      {hasChildren && isExpanded && (
        <div>
          {account.children!.map((child) => (
            <AccountNode key={child.guid} account={child} level={level + 1} />
          ))}
        </div>
      )}
    </div>
  );
};

export const AccountTree: React.FC<AccountTreeProps> = ({ accounts }) => {
  return (
    <div className="bg-white shadow overflow-hidden sm:rounded-lg">
      <div className="px-4 py-5 sm:px-6 border-b border-gray-200">
        <h3 className="text-lg leading-6 font-medium text-gray-900">
          Account Hierarchy
        </h3>
      </div>
      <div className="divide-y divide-gray-200">
        {accounts.map((account) => (
          <AccountNode key={account.guid} account={account} level={0} />
        ))}
      </div>
    </div>
  );
};
