package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/udai-kiran/agentic-cash/internal/application/dto"
	"github.com/udai-kiran/agentic-cash/internal/domain/repository"
)

// CommodityHandler handles commodity-related HTTP requests
type CommodityHandler struct {
	commodityRepo repository.CommodityRepository
}

// NewCommodityHandler creates a new commodity handler
func NewCommodityHandler(commodityRepo repository.CommodityRepository) *CommodityHandler {
	return &CommodityHandler{
		commodityRepo: commodityRepo,
	}
}

// GetCurrencies retrieves all currency commodities
func (h *CommodityHandler) GetCurrencies(c *gin.Context) {
	commodities, err := h.commodityRepo.FindCurrencies(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Failed to retrieve currencies",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	response := make([]dto.CommodityResponse, len(commodities))
	for i, commodity := range commodities {
		response[i] = dto.CommodityResponse{
			GUID:      commodity.GUID,
			Namespace: commodity.Namespace,
			Mnemonic:  commodity.Mnemonic,
			Fullname:  commodity.Fullname,
			Fraction:  commodity.Fraction,
		}
	}

	c.JSON(http.StatusOK, response)
}
