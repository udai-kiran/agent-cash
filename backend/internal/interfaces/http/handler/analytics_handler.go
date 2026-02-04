package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/udai-kiran/agentic-cash/internal/application/dto"
	"github.com/udai-kiran/agentic-cash/internal/application/service"
)

// AnalyticsHandler handles analytics-related HTTP requests
type AnalyticsHandler struct {
	analyticsService *service.AnalyticsService
}

// NewAnalyticsHandler creates a new analytics handler
func NewAnalyticsHandler(analyticsService *service.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsService: analyticsService,
	}
}

// GetIncomeExpense returns income vs expense analytics
func (h *AnalyticsHandler) GetIncomeExpense(c *gin.Context) {
	// Parse date range
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	var startDate, endDate time.Time
	var err error

	if startDateStr == "" {
		// Default to last 12 months
		startDate = time.Now().AddDate(-1, 0, 0)
	} else {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Bad Request",
				Message: "Invalid start_date format. Use YYYY-MM-DD",
				Code:    http.StatusBadRequest,
			})
			return
		}
	}

	if endDateStr == "" {
		endDate = time.Now()
	} else {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Bad Request",
				Message: "Invalid end_date format. Use YYYY-MM-DD",
				Code:    http.StatusBadRequest,
			})
			return
		}
	}

	response, err := h.analyticsService.GetIncomeExpense(c.Request.Context(), startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Failed to calculate income/expense",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetCategoryBreakdown returns spending breakdown by category
func (h *AnalyticsHandler) GetCategoryBreakdown(c *gin.Context) {
	// Parse date range
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	var startDate, endDate time.Time
	var err error

	if startDateStr == "" {
		// Default to last 6 months
		startDate = time.Now().AddDate(0, -6, 0)
	} else {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Bad Request",
				Message: "Invalid start_date format. Use YYYY-MM-DD",
				Code:    http.StatusBadRequest,
			})
			return
		}
	}

	if endDateStr == "" {
		endDate = time.Now()
	} else {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Bad Request",
				Message: "Invalid end_date format. Use YYYY-MM-DD",
				Code:    http.StatusBadRequest,
			})
			return
		}
	}

	response, err := h.analyticsService.GetCategoryBreakdown(c.Request.Context(), startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Failed to calculate category breakdown",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetNetWorth returns current net worth
func (h *AnalyticsHandler) GetNetWorth(c *gin.Context) {
	response, err := h.analyticsService.GetNetWorth(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Failed to calculate net worth",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, response)
}
