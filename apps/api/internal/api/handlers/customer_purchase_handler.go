package handlers

import (
	"github.com/gilabs/crm-healthcare/api/internal/domain/customer_purchase"
	customerpurchaseservice "github.com/gilabs/crm-healthcare/api/internal/service/customer_purchase"
	"github.com/gilabs/crm-healthcare/api/pkg/errors"
	"github.com/gilabs/crm-healthcare/api/pkg/response"
	"github.com/gin-gonic/gin"
)

type CustomerPurchaseHandler struct {
	service *customerpurchaseservice.Service
}

func NewCustomerPurchaseHandler(service *customerpurchaseservice.Service) *CustomerPurchaseHandler {
	return &CustomerPurchaseHandler{service: service}
}

func (h *CustomerPurchaseHandler) GetByAccount(c *gin.Context) {
	accountID := c.Param("id")

	histories, err := h.service.GetByAccountID(accountID)
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	responses := make([]customer_purchase.CustomerPurchaseResponse, len(histories))
	for i, p := range histories {
		responses[i] = *p.ToResponse()
	}

	response.SuccessResponse(c, responses, nil)
}

func (h *CustomerPurchaseHandler) GetAnalytics(c *gin.Context) {
	accountID := c.Param("id")

	// This is the old generic analytics, keeping for compatibility but frontend uses GetProductAnalytics now
	analytics, err := h.service.GetAnalytics(accountID)
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, analytics, nil)
}

func (h *CustomerPurchaseHandler) GetProductAnalytics(c *gin.Context) {
	accountID := c.Param("id")

	analytics, err := h.service.GetProductAnalytics(accountID)
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, analytics, nil)
}

func (h *CustomerPurchaseHandler) GetSummary(c *gin.Context) {
	accountID := c.Param("id")

	summary, err := h.service.GetSummary(accountID)
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, summary, nil)
}
