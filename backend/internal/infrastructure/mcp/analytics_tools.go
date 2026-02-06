package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// AnalyticsDateRangeParams defines parameters for analytics tools with date range
type AnalyticsDateRangeParams struct {
	StartDate string `json:"start_date,omitempty" jsonschema:"Start date in YYYY-MM-DD format (defaults to 1 month ago)"`
	EndDate   string `json:"end_date,omitempty" jsonschema:"End date in YYYY-MM-DD format (defaults to today)"`
}

// handleAnalyticsExpenses handles the analytics_expenses tool
func (s *MCPServer) handleAnalyticsExpenses(ctx context.Context, req *mcp.CallToolRequest, params *AnalyticsDateRangeParams) (*mcp.CallToolResult, any, error) {
	startDate, endDate := parseDateRange(params.StartDate, params.EndDate)

	result, err := s.analyticsService.GetIncomeExpense(ctx, startDate, endDate)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get expenses: %w", err)
	}

	jsonData, _ := json.MarshalIndent(map[string]any{
		"currency":      result.CurrencyMnemonic,
		"total_expense": result.TotalExpense,
		"data":          result.Data,
		"start_date":    startDate.Format("2006-01-02"),
		"end_date":      endDate.Format("2006-01-02"),
	}, "", "  ")

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(jsonData)},
		},
	}, nil, nil
}

// handleAnalyticsIncome handles the analytics_income tool
func (s *MCPServer) handleAnalyticsIncome(ctx context.Context, req *mcp.CallToolRequest, params *AnalyticsDateRangeParams) (*mcp.CallToolResult, any, error) {
	startDate, endDate := parseDateRange(params.StartDate, params.EndDate)

	result, err := s.analyticsService.GetIncomeExpense(ctx, startDate, endDate)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get income: %w", err)
	}

	jsonData, _ := json.MarshalIndent(map[string]any{
		"currency":     result.CurrencyMnemonic,
		"total_income": result.TotalIncome,
		"data":         result.Data,
		"start_date":   startDate.Format("2006-01-02"),
		"end_date":     endDate.Format("2006-01-02"),
	}, "", "  ")

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(jsonData)},
		},
	}, nil, nil
}

// handleAnalyticsCashflow handles the analytics_cashflow tool
func (s *MCPServer) handleAnalyticsCashflow(ctx context.Context, req *mcp.CallToolRequest, params *AnalyticsDateRangeParams) (*mcp.CallToolResult, any, error) {
	startDate, endDate := parseDateRange(params.StartDate, params.EndDate)

	result, err := s.analyticsService.GetIncomeExpense(ctx, startDate, endDate)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get cashflow: %w", err)
	}

	jsonData, _ := json.MarshalIndent(map[string]any{
		"currency":      result.CurrencyMnemonic,
		"total_income":  result.TotalIncome,
		"total_expense": result.TotalExpense,
		"net_cashflow":  result.NetTotal,
		"data":          result.Data,
		"start_date":    startDate.Format("2006-01-02"),
		"end_date":      endDate.Format("2006-01-02"),
	}, "", "  ")

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(jsonData)},
		},
	}, nil, nil
}

// parseDateRange parses date strings or provides defaults
func parseDateRange(startDateStr, endDateStr string) (time.Time, time.Time) {
	var startDate, endDate time.Time

	if startDateStr != "" {
		t, err := time.Parse("2006-01-02", startDateStr)
		if err == nil {
			startDate = t
		}
	}

	if endDateStr != "" {
		t, err := time.Parse("2006-01-02", endDateStr)
		if err == nil {
			endDate = t
		}
	}

	// Default to current month if no dates provided
	if startDate.IsZero() {
		startDate = time.Now().AddDate(0, -1, 0) // One month ago
	}
	if endDate.IsZero() {
		endDate = time.Now()
	}

	return startDate, endDate
}
