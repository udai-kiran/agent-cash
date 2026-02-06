package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/udai-kiran/agentic-cash/internal/domain/entity"
)

// AccountsListParams defines parameters for accounts_list tool
type AccountsListParams struct {
	Type string `json:"type,omitempty" jsonschema:"Account type filter (ASSET, LIABILITY, EQUITY, INCOME, EXPENSE, etc.)"`
}

// AccountsGetParams defines parameters for accounts_get tool
type AccountsGetParams struct {
	GUID string `json:"guid" jsonschema:"required,Account GUID to retrieve"`
}

// AccountsBalanceParams defines parameters for accounts_balance tool
type AccountsBalanceParams struct {
	GUID string `json:"guid" jsonschema:"required,Account GUID to get balance for"`
}

// handleAccountsList handles the accounts_list tool
func (s *MCPServer) handleAccountsList(ctx context.Context, req *mcp.CallToolRequest, params *AccountsListParams) (*mcp.CallToolResult, any, error) {
	var accountType entity.AccountType
	if params.Type != "" {
		accountType = entity.AccountType(params.Type)
	}

	var accounts []*entity.Account
	var err error

	if accountType != "" {
		accounts, err = s.accountRepo.FindByType(ctx, accountType)
	} else {
		accounts, err = s.accountRepo.FindAll(ctx)
	}

	if err != nil {
		return nil, nil, fmt.Errorf("failed to list accounts: %w", err)
	}

	if len(accounts) == 0 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "No accounts found."},
			},
		}, nil, nil
	}

	result := make([]map[string]any, 0, len(accounts))
	for _, acc := range accounts {
		result = append(result, formatAccount(acc))
	}

	jsonData, _ := json.MarshalIndent(map[string]any{
		"accounts": result,
		"count":    len(accounts),
	}, "", "  ")

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(jsonData)},
		},
	}, nil, nil
}

// handleAccountsGet handles the accounts_get tool
func (s *MCPServer) handleAccountsGet(ctx context.Context, req *mcp.CallToolRequest, params *AccountsGetParams) (*mcp.CallToolResult, any, error) {
	if params.GUID == "" {
		return nil, nil, fmt.Errorf("missing required parameter: guid")
	}

	account, err := s.accountRepo.FindByGUID(ctx, params.GUID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get account: %w", err)
	}

	if account == nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Account not found."},
			},
		}, nil, nil
	}

	jsonData, _ := json.MarshalIndent(map[string]any{
		"account": formatAccount(account),
	}, "", "  ")

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(jsonData)},
		},
	}, nil, nil
}

// handleAccountsHierarchy handles the accounts_hierarchy tool
func (s *MCPServer) handleAccountsHierarchy(ctx context.Context, req *mcp.CallToolRequest, params *struct{}) (*mcp.CallToolResult, any, error) {
	accounts, err := s.accountRepo.FindHierarchy(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get account hierarchy: %w", err)
	}

	if len(accounts) == 0 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "No accounts found."},
			},
		}, nil, nil
	}

	jsonData, _ := json.MarshalIndent(formatAccountHierarchy(accounts), "", "  ")

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(jsonData)},
		},
	}, nil, nil
}

// handleAccountsBalance handles the accounts_balance tool
func (s *MCPServer) handleAccountsBalance(ctx context.Context, req *mcp.CallToolRequest, params *AccountsBalanceParams) (*mcp.CallToolResult, any, error) {
	if params.GUID == "" {
		return nil, nil, fmt.Errorf("missing required parameter: guid")
	}

	balanceNum, balanceDenom, err := s.accountRepo.GetBalance(ctx, params.GUID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get account balance: %w", err)
	}

	// Calculate human-readable balance
	balance := float64(balanceNum) / float64(balanceDenom)

	jsonData, _ := json.MarshalIndent(map[string]any{
		"guid":            params.GUID,
		"balance_num":     balanceNum,
		"balance_denom":   balanceDenom,
		"balance_decimal": balance,
	}, "", "  ")

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(jsonData)},
		},
	}, nil, nil
}