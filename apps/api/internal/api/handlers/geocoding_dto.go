package handlers

// GeocodeRequest represents geocoding request.
type GeocodeRequest struct {
	Address string `json:"address" binding:"required"`
}

// GeocodeResponse represents geocoding response.
type GeocodeResponse struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Address   string  `json:"address"`
}
