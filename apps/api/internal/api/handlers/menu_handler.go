package handlers

import (
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	"github.com/gilabs/crm-healthcare/api/pkg/errors"
	"github.com/gilabs/crm-healthcare/api/pkg/response"
	"github.com/gin-gonic/gin"
)

type MenuHandler struct {
	menuRepo interfaces.MenuRepository
}

func NewMenuHandler(menuRepo interfaces.MenuRepository) *MenuHandler {
	return &MenuHandler{
		menuRepo: menuRepo,
	}
}

// List handles list menus request
func (h *MenuHandler) List(c *gin.Context) {
	menus, err := h.menuRepo.List()
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	// Convert to response DTOs
	var menuResponses []interface{}
	for _, m := range menus {
		menuResponses = append(menuResponses, m.ToMenuResponse())
	}

	response.SuccessResponse(c, menuResponses, nil)
}

// GetByID handles get menu by ID request
func (h *MenuHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	menu, err := h.menuRepo.FindByID(id)
	if err != nil {
		errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
			"resource": "menu",
			"menu_id":  id,
		}, nil)
		return
	}

	response.SuccessResponse(c, menu.ToMenuResponse(), nil)
}
