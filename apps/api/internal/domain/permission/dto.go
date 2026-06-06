package permission

import "time"

// PermissionResponse represents permission response DTO.
type PermissionResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Resource    string    `json:"resource"`
	Action      string    `json:"action"`
	Description string    `json:"description"`
	Access      bool      `json:"access"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// PermissionShortResponse represents a concise permission response DTO.
type PermissionShortResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Code     string `json:"code"`
	Resource string `json:"resource"`
	Action   string `json:"action"`
	Access   bool   `json:"access"`
}

// PermissionSimpleResponse represents minimal permission response for list endpoint.
type PermissionSimpleResponse struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Code     string  `json:"code"`
	Resource string  `json:"resource"`
	Action   string  `json:"action"`
	MenuID   *string `json:"menu_id,omitempty"`
	Menu     *Menu   `json:"menu,omitempty"`
}

// ToPermissionResponse converts Permission to PermissionResponse.
func (p *Permission) ToPermissionResponse() *PermissionResponse {
	return &PermissionResponse{
		ID:          p.ID,
		Name:        p.Name,
		Code:        p.Code,
		Resource:    p.Resource,
		Action:      p.Action,
		Description: p.Description,
		Access:      false,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

// ToPermissionShortResponse converts Permission to PermissionShortResponse.
func (p *Permission) ToPermissionShortResponse() *PermissionShortResponse {
	return &PermissionShortResponse{
		ID:       p.ID,
		Name:     p.Name,
		Code:     p.Code,
		Resource: p.Resource,
		Action:   p.Action,
		Access:   true,
	}
}

// ToPermissionSimpleResponse converts Permission to PermissionSimpleResponse.
func (p *Permission) ToPermissionSimpleResponse() *PermissionSimpleResponse {
	return &PermissionSimpleResponse{
		ID:       p.ID,
		Name:     p.Name,
		Code:     p.Code,
		Resource: p.Resource,
		Action:   p.Action,
		MenuID:   p.MenuID,
		Menu:     p.Menu,
	}
}

// MenuResponse represents menu response DTO with nested structure.
type MenuResponse struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	Icon      string         `json:"icon"`
	URL       string         `json:"url"`
	ParentID  *string        `json:"parent_id,omitempty"`
	Children  []MenuResponse `json:"children,omitempty"`
	Order     int            `json:"order"`
	Status    string         `json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

// ToMenuResponse converts Menu to MenuResponse recursively.
func (m *Menu) ToMenuResponse() *MenuResponse {
	resp := &MenuResponse{
		ID:        m.ID,
		Name:      m.Name,
		Icon:      m.Icon,
		URL:       m.URL,
		ParentID:  m.ParentID,
		Order:     m.Order,
		Status:    m.Status,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
	if len(m.Children) > 0 {
		resp.Children = make([]MenuResponse, len(m.Children))
		for i, child := range m.Children {
			resp.Children[i] = *child.ToMenuResponse()
		}
	}
	return resp
}
