package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/udai-kiran/agentic-cash/internal/application/dto"
	"github.com/udai-kiran/agentic-cash/internal/domain/entity"
	"github.com/udai-kiran/agentic-cash/internal/domain/repository"
	"github.com/udai-kiran/agentic-cash/pkg/gnucash"
)

// TransactionHandler handles transaction-related HTTP requests
type TransactionHandler struct {
	transactionRepo repository.TransactionRepository
}

// NewTransactionHandler creates a new transaction handler
func NewTransactionHandler(transactionRepo repository.TransactionRepository) *TransactionHandler {
	return &TransactionHandler{
		transactionRepo: transactionRepo,
	}
}

// GetTransactions retrieves transactions with optional filtering
func (h *TransactionHandler) GetTransactions(c *gin.Context) {
	filter := &repository.TransactionFilter{
		Limit:  50, // default
		Offset: 0,
	}

	// Parse query parameters
	if accountGUID := c.Query("account_guid"); accountGUID != "" {
		filter.AccountGUID = &accountGUID
	}

	if startDateStr := c.Query("start_date"); startDateStr != "" {
		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err == nil {
			filter.StartDate = &startDate
		}
	}

	if endDateStr := c.Query("end_date"); endDateStr != "" {
		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err == nil {
			filter.EndDate = &endDate
		}
	}

	if description := c.Query("description"); description != "" {
		filter.Description = &description
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			filter.Limit = limit
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			filter.Offset = offset
		}
	}

	// Get transactions
	transactions, err := h.transactionRepo.FindAll(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Failed to retrieve transactions",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	// Get total count
	total, err := h.transactionRepo.Count(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Failed to count transactions",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	// Convert to response
	response := dto.TransactionListResponse{
		Transactions: make([]dto.TransactionResponse, len(transactions)),
		Total:        total,
		Limit:        filter.Limit,
		Offset:       filter.Offset,
	}

	for i, tx := range transactions {
		response.Transactions[i] = h.toTransactionResponse(tx)
	}

	c.JSON(http.StatusOK, response)
}

// GetTransaction retrieves a single transaction by GUID
func (h *TransactionHandler) GetTransaction(c *gin.Context) {
	guid := c.Param("guid")

	transaction, err := h.transactionRepo.FindByGUID(c.Request.Context(), guid)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "Not Found",
			Message: "Transaction not found",
			Code:    http.StatusNotFound,
		})
		return
	}

	c.JSON(http.StatusOK, h.toTransactionResponse(transaction))
}

// toTransactionResponse converts entity.Transaction to dto.TransactionResponse
func (h *TransactionHandler) toTransactionResponse(tx *entity.Transaction) dto.TransactionResponse {
	splits := make([]dto.SplitResponse, len(tx.Splits))
	for i, split := range tx.Splits {
		splits[i] = dto.SplitResponse{
			GUID:           split.GUID,
			TxGUID:         split.TxGUID,
			AccountGUID:    split.AccountGUID,
			Memo:           split.Memo,
			Action:         split.Action,
			ReconcileState: split.ReconcileState,
			ValueNum:       split.ValueNum,
			ValueDenom:     split.ValueDenom,
			QuantityNum:    split.QuantityNum,
			QuantityDenom:  split.QuantityDenom,
			Value:          gnucash.FormatAmount(split.ValueNum, split.ValueDenom),
			Quantity:       gnucash.FormatAmount(split.QuantityNum, split.QuantityDenom),
		}

		if split.Account != nil {
			splits[i].Account = &dto.AccountSummary{
				GUID: split.Account.GUID,
				Name: split.Account.Name,
				Type: string(split.Account.AccountType),
			}
		}
	}

	return dto.TransactionResponse{
		GUID:             tx.GUID,
		CurrencyGUID:     tx.CurrencyGUID,
		CurrencyMnemonic: tx.CurrencyMnemonic,
		Num:              tx.Num,
		PostDate:         tx.PostDate,
		EnterDate:        tx.EnterDate,
		Description:      tx.Description,
		Splits:           splits,
	}
}
