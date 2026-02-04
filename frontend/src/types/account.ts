export interface Account {
  guid: string;
  name: string;
  type: string;
  code?: string;
  description?: string;
  hidden: boolean;
  placeholder: boolean;
  parent_guid?: string;
  balance: string;
  balance_num: number;
  balance_denom: number;
  commodity_mnemonic?: string;
  children?: Account[];
}

export interface AccountBalance {
  guid: string;
  name: string;
  balance: string;
  balance_num: number;
  balance_denom: number;
}
