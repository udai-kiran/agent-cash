import React from 'react';
import { useAccountHierarchy } from '../hooks/useAccounts';
import { AccountTree } from '../components/accounts/AccountTree';
import { LoadingSpinner } from '../components/common/LoadingSpinner';
import { ErrorMessage } from '../components/common/ErrorMessage';
import { Layout } from '../components/common/Layout';

export const AccountsPage: React.FC = () => {
  const { data: accounts, isLoading, error } = useAccountHierarchy();

  return (
    <Layout>
      <div className="px-4 py-8 sm:px-0">
        <h1 className="text-3xl font-bold text-gray-900 mb-6">Accounts</h1>

        {isLoading && <LoadingSpinner />}

        {error && (
          <ErrorMessage
            message={error instanceof Error ? error.message : 'Failed to load accounts'}
          />
        )}

        {accounts && <AccountTree accounts={accounts} />}
      </div>
    </Layout>
  );
};
