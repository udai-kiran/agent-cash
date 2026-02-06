package mcp

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/udai-kiran/agentic-cash/internal/application/service"
	"github.com/udai-kiran/agentic-cash/internal/domain/repository"
)

// MCPServer represents the MCP server for GnuCash
type MCPServer struct {
	accountRepo      repository.AccountRepository
	transactionRepo  repository.TransactionRepository
	commodityRepo    repository.CommodityRepository
	analyticsService *service.AnalyticsService
	server           *mcp.Server
	httpServer       *http.Server
	port             int
}

// NewMCPServer creates a new MCP server instance
func NewMCPServer(
	accountRepo repository.AccountRepository,
	transactionRepo repository.TransactionRepository,
	commodityRepo repository.CommodityRepository,
	analyticsService *service.AnalyticsService,
) *MCPServer {
	port := 8081
	if envPort := os.Getenv("MCP_PORT"); envPort != "" {
		fmt.Sscanf(envPort, "%d", &port)
	}

	serverName := os.Getenv("MCP_SERVER_NAME")
	if serverName == "" {
		serverName = "gnucash-mcp-server"
	}

	serverVersion := os.Getenv("MCP_SERVER_VERSION")
	if serverVersion == "" {
		serverVersion = "1.0.0"
	}

	// Create the MCP server
	mcpServer := mcp.NewServer(&mcp.Implementation{
		Name:    serverName,
		Version: serverVersion,
	}, nil)

	s := &MCPServer{
		accountRepo:      accountRepo,
		transactionRepo:  transactionRepo,
		commodityRepo:    commodityRepo,
		analyticsService: analyticsService,
		server:           mcpServer,
		port:             port,
	}

	// Register all tools
	s.registerTools()

	return s
}

// registerTools registers all available MCP tools
func (s *MCPServer) registerTools() {
	// Account tools
	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "accounts_list",
		Description: "List all accounts or filter by account type (ASSET, LIABILITY, EQUITY, INCOME, EXPENSE, etc.)",
	}, s.handleAccountsList)

	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "accounts_get",
		Description: "Get detailed information about a specific account by GUID",
	}, s.handleAccountsGet)

	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "accounts_hierarchy",
		Description: "Get the complete account hierarchy tree",
	}, s.handleAccountsHierarchy)

	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "accounts_balance",
		Description: "Get the current balance of a specific account by GUID",
	}, s.handleAccountsBalance)

	// Transaction tools
	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "transactions_list",
		Description: "List transactions, optionally filtered by account GUID and date range",
	}, s.handleTransactionsList)

	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "transactions_get",
		Description: "Get detailed information about a specific transaction by GUID",
	}, s.handleTransactionsGet)

	// Analytics tools
	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "analytics_expenses",
		Description: "Get expense analysis for a date range, grouped by account",
	}, s.handleAnalyticsExpenses)

	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "analytics_income",
		Description: "Get income analysis for a date range, grouped by account",
	}, s.handleAnalyticsIncome)

	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "analytics_cashflow",
		Description: "Get cash flow analysis showing income and expenses over time",
	}, s.handleAnalyticsCashflow)

	// Commodity tools
	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "commodities_list",
		Description: "List all commodities (currencies) in the database",
	}, s.handleCommoditiesList)

	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "commodities_get",
		Description: "Get detailed information about a specific commodity by GUID",
	}, s.handleCommoditiesGet)

	log.Printf("Registered %d MCP tools", 11)
}

// Start runs the MCP server with HTTP transport
func (s *MCPServer) Start(ctx context.Context) error {
	// Create the streamable HTTP handler
	handler := mcp.NewStreamableHTTPHandler(func(req *http.Request) *mcp.Server {
		return s.server
	}, nil)

	// Create HTTP server
	addr := fmt.Sprintf(":%d", s.port)
	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	log.Printf("GnuCash MCP Server starting on http://0.0.0.0%s", addr)
	log.Printf("Available tools: accounts_*, transactions_*, analytics_*, commodities_*")

	// Start the HTTP server in a goroutine
	errChan := make(chan error, 1)
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	// Wait for context cancellation or error
	select {
	case <-ctx.Done():
		return s.Shutdown(context.Background())
	case err := <-errChan:
		return err
	}
}

// Shutdown gracefully shuts down the MCP server
func (s *MCPServer) Shutdown(ctx context.Context) error {
	log.Println("Shutting down MCP server...")
	if s.httpServer != nil {
		return s.httpServer.Shutdown(ctx)
	}
	return nil
}
