export interface Split {
  guid: string;
  tx_guid: string;
  account_guid: string;
  memo?: string;
  action?: string;
  reconcile_state: string;
  value_num: number;
  value_denom: number;
  quantity_num: number;
  quantity_denom: number;
  value: string;
  quantity: string;
  account?: {
    guid: string;
    name: string;
    type: string;
  };
}

export interface Transaction {
  guid: string;
  currency_guid: string;
  currency_mnemonic?: string;
  num?: string;
  post_date: string;
  enter_date: string;
  description?: string;
  splits: Split[];
}

export interface TransactionFilter {
  account_guid?: string;
  start_date?: string;
  end_date?: string;
  description?: string;
  limit?: number;
  offset?: number;
}
