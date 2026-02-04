export interface IncomeExpenseData {
  period: string;
  income: string;
  expense: string;
  net: string;
}

export interface IncomeExpenseResponse {
  data: IncomeExpenseData[];
  total_income: string;
  total_expense: string;
  net_total: string;
  currency_mnemonic?: string;
}

export interface CategoryBreakdownItem {
  category: string;
  amount: string;
  count: number;
}

export interface CategoryBreakdownResponse {
  income: CategoryBreakdownItem[];
  expense: CategoryBreakdownItem[];
  currency_mnemonic?: string;
}

export interface NetWorthItem {
  account_name: string;
  account_type: string;
  balance: string;
}

export interface NetWorthResponse {
  assets: NetWorthItem[];
  liabilities: NetWorthItem[];
  total_assets: string;
  total_liabilities: string;
  net_worth: string;
  currency_mnemonic?: string;
}
