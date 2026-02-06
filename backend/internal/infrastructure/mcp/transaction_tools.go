package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/udai-kiran/agentic-cash/internal/domain/repository"
)

// TransactionsListParams defines parameters for transactions_list tool
type TransactionsListParams struct {
	AccountGUID string `json:"account_guid,omitempty" jsonschema:"Filter by account GUID"`
	StartDate   string `json:"start_date,omitempty" jsonschema:"Start date in YYYY-MM-DD format"`
	EndDate     string `json:"end_date,omitempty" jsonschema:"End date in YYYY-MM-DD format"`
	Description string `json:"description,omitempty" jsonschema:"Filter by description (partial match)"`
}

// TransactionsGetParams defines parameters for transactions_get tool
type TransactionsGetParams struct {
	GUID string `json:"guid" jsonschema:"required,Transaction GUID to retrieve"`
}

// handleTransactionsList handles the transactions_list tool
func (s *MCPServer) handleTransactionsList(ctx context.Context, req *mcp.CallToolRequest, params *TransactionsListParams) (*mcp.CallToolResult, any, error) {
	// Parse dates
	var startDate, endDate *time.Time
	if params.StartDate != "" {
		t, err := time.Parse("2006-01-02", params.StartDate)
		if err == nil {
			startDate = &t
		}
	}
	if params.EndDate != "" {
		t, err := time.Parse("2006-01-02", params.EndDate)
		if err == nil {
			endDate = &t
		}
	}

	// Create filter
	filter := &repository.TransactionFilter{}
	if params.AccountGUID != "" {
		filter.AccountGUID = &params.AccountGUID
	}
	if startDate != nil {
		filter.StartDate = startDate
	}
	if endDate != nil {
		filter.EndDate = endDate
	}
	if params.Description != "" {
		filter.Description = &params.Description
	}

	transactions, err := s.transactionRepo.FindAll(ctx, filter)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list transactions: %w", err)
	}

	if len(transactions) == 0 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "No transactions found."},
			},
		}, nil, nil
	}

	result := make([]map[string]any, 0, len(transactions))
	for _, tx := range transactions {
		result = append(result, formatTransaction(tx))
	}

	jsonData, _ := json.MarshalIndent(map[string]any{
		"transactions": result,
		"count":        len(transactions),
	}, "", "  ")

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(jsonData)},
		},
	}, nil, nil
}

// handleTransactionsGet handles the transactions_get tool
func (s *MCPServer) handleTransactionsGet(ctx context.Context, req *mcp.CallToolRequest, params *TransactionsGetParams) (*mcp.CallToolResult, any, error) {
	if params.GUID == "" {
		return nil, nil, fmt.Errorf("missing required parameter: guid")
	}

	transaction, err := s.transactionRepo.FindByGUID(ctx, params.GUID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	if transaction == nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Transaction not found."},
			},
		}, nil, nil
	}

	jsonData, _ := json.MarshalIndent(map[string]any{
		"transaction": formatTransaction(transaction),
	}, "", "  ")

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(jsonData)},
		},
	}, nil, nil
}
