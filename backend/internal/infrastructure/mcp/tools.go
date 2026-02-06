package mcp

import (
	"github.com/udai-kiran/agentic-cash/internal/domain/entity"
)

// formatAccount converts an account to a map for JSON serialization
func formatAccount(acc *entity.Account) map[string]any {
	result := map[string]any{
		"guid":          acc.GUID,
		"name":          acc.Name,
		"type":          string(acc.AccountType),
		"description":   acc.Description,
		"code":          acc.Code,
		"hidden":        acc.Hidden,
		"placeholder":   acc.Placeholder,
		"commodity":     acc.CommodityMnemonic,
		"balance":       acc.Balance.StringFixed(2),
		"balance_num":   acc.BalanceNum,
		"balance_denom": acc.BalanceDenom,
	}

	if acc.ParentGUID != nil {
		result["parent_guid"] = *acc.ParentGUID
	}

	return result
}

// formatAccountHierarchy builds a hierarchical tree of accounts
func formatAccountHierarchy(accounts []*entity.Account) map[string]any {
	accountMap := make(map[string]*entity.Account)
	childMap := make(map[string][]*entity.Account)

	// Build maps
	for _, acc := range accounts {
		accountMap[acc.GUID] = acc
		if acc.ParentGUID != nil {
			childMap[*acc.ParentGUID] = append(childMap[*acc.ParentGUID], acc)
		}
	}

	// Build hierarchy
	var buildTree func(*entity.Account) map[string]any
	buildTree = func(acc *entity.Account) map[string]any {
		node := formatAccount(acc)
		if children, ok := childMap[acc.GUID]; ok {
			node["children"] = make([]map[string]any, 0, len(children))
			for _, child := range children {
				node["children"] = append(node["children"].([]map[string]any), buildTree(child))
			}
		}
		return node
	}

	// Find root accounts (those without parent or with parent not in list)
	var roots []*entity.Account
	for _, acc := range accounts {
		if acc.ParentGUID == nil || accountMap[*acc.ParentGUID] == nil {
			roots = append(roots, acc)
		}
	}

	result := make([]map[string]any, 0, len(roots))
	for _, root := range roots {
		result = append(result, buildTree(root))
	}

	return map[string]any{"hierarchy": result, "count": len(accounts)}
}

// formatTransaction converts a transaction to a map
func formatTransaction(tx *entity.Transaction) map[string]any {
	splits := make([]map[string]any, 0, len(tx.Splits))
	for _, s := range tx.Splits {
		splits = append(splits, map[string]any{
			"guid":           s.GUID,
			"account_guid":   s.AccountGUID,
			"account_name":   s.Account.Name,
			"account_type":   string(s.Account.AccountType),
			"value_num":      s.ValueNum,
			"value_denom":    s.ValueDenom,
			"quantity_num":   s.QuantityNum,
			"quantity_denom": s.QuantityDenom,
			"memo":           s.Memo,
			"action":         s.Action,
			"reconcile_state": s.ReconcileState,
		})
	}

	return map[string]any{
		"guid":         tx.GUID,
		"currency":     tx.CurrencyMnemonic,
		"currency_guid": tx.CurrencyGUID,
		"number":       tx.Num,
		"post_date":    tx.PostDate.Format("2006-01-02"),
		"enter_date":   tx.EnterDate.Format("2006-01-02"),
		"description":  tx.Description,
		"splits":       splits,
	}
}
