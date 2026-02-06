package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// CommoditiesGetParams defines parameters for commodities_get tool
type CommoditiesGetParams struct {
	GUID string `json:"guid" jsonschema:"required,Commodity GUID to retrieve"`
}

// handleCommoditiesList handles the commodities_list tool
func (s *MCPServer) handleCommoditiesList(ctx context.Context, req *mcp.CallToolRequest, params *struct{}) (*mcp.CallToolResult, any, error) {
	commodities, err := s.commodityRepo.FindCurrencies(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list commodities: %w", err)
	}

	if len(commodities) == 0 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "No commodities found."},
			},
		}, nil, nil
	}

	result := make([]map[string]any, 0, len(commodities))
	for _, c := range commodities {
		result = append(result, map[string]any{
			"guid":      c.GUID,
			"namespace": c.Namespace,
			"mnemonic":  c.Mnemonic,
			"full_name": c.Fullname,
			"fraction":  c.Fraction,
		})
	}

	jsonData, _ := json.MarshalIndent(map[string]any{
		"commodities": result,
		"count":       len(commodities),
	}, "", "  ")

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(jsonData)},
		},
	}, nil, nil
}

// handleCommoditiesGet handles the commodities_get tool
func (s *MCPServer) handleCommoditiesGet(ctx context.Context, req *mcp.CallToolRequest, params *CommoditiesGetParams) (*mcp.CallToolResult, any, error) {
	if params.GUID == "" {
		return nil, nil, fmt.Errorf("missing required parameter: guid")
	}

	commodity, err := s.commodityRepo.FindByGUID(ctx, params.GUID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get commodity: %w", err)
	}

	if commodity == nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Commodity not found."},
			},
		}, nil, nil
	}

	jsonData, _ := json.MarshalIndent(map[string]any{
		"guid":      commodity.GUID,
		"namespace": commodity.Namespace,
		"mnemonic":  commodity.Mnemonic,
		"full_name": commodity.Fullname,
		"fraction":  commodity.Fraction,
	}, "", "  ")

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(jsonData)},
		},
	}, nil, nil
}
