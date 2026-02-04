package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/udai-kiran/agentic-cash/internal/application/dto"
	"github.com/udai-kiran/agentic-cash/internal/domain/entity"
	"github.com/udai-kiran/agentic-cash/internal/domain/repository"
	"github.com/udai-kiran/agentic-cash/pkg/gnucash"
)

// AccountHandler handles account-related HTTP requests
type AccountHandler struct {
	accountRepo   repository.AccountRepository
	commodityRepo repository.CommodityRepository
}

// NewAccountHandler creates a new account handler
func NewAccountHandler(accountRepo repository.AccountRepository, commodityRepo repository.CommodityRepository) *AccountHandler {
	return &AccountHandler{
		accountRepo:   accountRepo,
		commodityRepo: commodityRepo,
	}
}

// GetAccounts retrieves all accounts
func (h *AccountHandler) GetAccounts(c *gin.Context) {
	accounts, err := h.accountRepo.FindAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Failed to retrieve accounts",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	response := make([]dto.AccountResponse, len(accounts))
	for i, account := range accounts {
		response[i] = h.toAccountResponse(account)
	}

	c.JSON(http.StatusOK, response)
}

// GetAccountHierarchy retrieves the account hierarchy
func (h *AccountHandler) GetAccountHierarchy(c *gin.Context) {
	roots, err := h.accountRepo.FindHierarchy(c.Request.Context())
	if err != nil {
		log.Printf("ERROR: GetAccountHierarchy: %v", err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Failed to retrieve account hierarchy",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	response := make([]dto.AccountResponse, len(roots))
	for i, root := range roots {
		response[i] = h.toAccountResponseWithChildren(root)
	}

	c.JSON(http.StatusOK, response)
}

// GetAccount retrieves a single account by GUID
func (h *AccountHandler) GetAccount(c *gin.Context) {
	guid := c.Param("guid")

	account, err := h.accountRepo.FindByGUID(c.Request.Context(), guid)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "Not Found",
			Message: "Account not found",
			Code:    http.StatusNotFound,
		})
		return
	}

	c.JSON(http.StatusOK, h.toAccountResponse(account))
}

// GetAccountBalance retrieves the balance for an account
func (h *AccountHandler) GetAccountBalance(c *gin.Context) {
	guid := c.Param("guid")

	account, err := h.accountRepo.FindByGUID(c.Request.Context(), guid)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "Not Found",
			Message: "Account not found",
			Code:    http.StatusNotFound,
		})
		return
	}

	balanceNum, balanceDenom, err := h.accountRepo.GetBalance(c.Request.Context(), guid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Failed to calculate balance",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, dto.AccountBalanceResponse{
		GUID:         account.GUID,
		Name:         account.Name,
		Balance:      gnucash.FormatAmount(balanceNum, balanceDenom),
		BalanceNum:   balanceNum,
		BalanceDenom: balanceDenom,
	})
}

// toAccountResponse converts an entity.Account to dto.AccountResponse
func (h *AccountHandler) toAccountResponse(account *entity.Account) dto.AccountResponse {
	return dto.AccountResponse{
		GUID:              account.GUID,
		Name:              account.Name,
		Type:              string(account.AccountType),
		Code:              account.Code,
		Description:       account.Description,
		Hidden:            account.Hidden,
		Placeholder:       account.Placeholder,
		ParentGUID:        account.ParentGUID,
		Balance:           gnucash.FormatAmount(account.BalanceNum, account.BalanceDenom),
		BalanceNum:        account.BalanceNum,
		BalanceDenom:      account.BalanceDenom,
		CommodityMnemonic: account.CommodityMnemonic,
	}
}

// toAccountResponseWithChildren converts an entity.Account to dto.AccountResponse with children
func (h *AccountHandler) toAccountResponseWithChildren(account *entity.Account) dto.AccountResponse {
	response := h.toAccountResponse(account)

	if len(account.Children) > 0 {
		response.Children = make([]dto.AccountResponse, len(account.Children))
		for i, child := range account.Children {
			response.Children[i] = h.toAccountResponseWithChildren(child)
		}
	}

	return response
}
