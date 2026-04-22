package handlers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gilabs/crm-healthcare/api/internal/config"
	googlecalendartokenservice "github.com/gilabs/crm-healthcare/api/internal/service/google_calendar_token"
	"github.com/gilabs/crm-healthcare/api/pkg/errors"
	"github.com/gilabs/crm-healthcare/api/pkg/response"
	"github.com/gin-gonic/gin"
)

// OAuthState represents the state parameter for OAuth flow
type OAuthState struct {
	UserID   string `json:"user_id"`
	Platform string `json:"platform"`
}

// encodeState encodes OAuthState to base64 string
func encodeState(state OAuthState) string {
	data, _ := json.Marshal(state)
	return base64.URLEncoding.EncodeToString(data)
}

// decodeState decodes base64 string to OAuthState
func decodeState(stateStr string) (OAuthState, error) {
	var state OAuthState
	data, err := base64.URLEncoding.DecodeString(stateStr)
	if err != nil {
		return state, err
	}
	err = json.Unmarshal(data, &state)
	return state, err
}

type GoogleCalendarAuthHandler struct {
	tokenService *googlecalendartokenservice.Service
	config       *config.GoogleCalendarConfig
}

func NewGoogleCalendarAuthHandler(tokenService *googlecalendartokenservice.Service, cfg *config.GoogleCalendarConfig) *GoogleCalendarAuthHandler {
	return &GoogleCalendarAuthHandler{
		tokenService: tokenService,
		config:       cfg,
	}
}

// GetAuthURL returns the OAuth2 authorization URL
func (h *GoogleCalendarAuthHandler) GetAuthURL(c *gin.Context) {
	userID := ""
	if userIDVal, exists := c.Get("user_id"); exists {
		if id, ok := userIDVal.(string); ok {
			userID = id
		}
	}

	if userID == "" {
		errors.ErrorResponse(c, "UNAUTHORIZED", nil, nil)
		return
	}

	// Get platform from query param (web or mobile)
	platform := c.Query("platform")
	if platform == "" {
		platform = "web" // default to web
	}

	// Build state with userID and platform
	oauthState := OAuthState{
		UserID:   userID,
		Platform: platform,
	}
	state := encodeState(oauthState)

	// Generate auth URL with platform-specific redirect URI
	authURL := h.tokenService.GetAuthURLForPlatform(state, platform)

	response.SuccessResponse(c, map[string]interface{}{
		"auth_url": authURL,
		"state":    state,
	}, nil)
}

// HandleCallback handles OAuth2 callback and stores token
// Note: This endpoint does NOT require authentication (called by Google OAuth redirect)
// User ID is extracted from the 'state' parameter set during auth URL generation
func (h *GoogleCalendarAuthHandler) HandleCallback(c *gin.Context) {
	code := c.Query("code")
	stateStr := c.Query("state")

	// Decode state to get userID and platform
	oauthState, err := decodeState(stateStr)
	if err != nil {
		// Fallback to old behavior if state decoding fails
		oauthState = OAuthState{
			UserID:   stateStr,
			Platform: "web",
		}
	}

	if code == "" {
		redirectURL := h.buildRedirectURL(oauthState, "missing_code", "")
		c.Redirect(http.StatusFound, redirectURL)
		return
	}

	if stateStr == "" {
		redirectURL := h.buildRedirectURL(oauthState, "missing_state", "")
		c.Redirect(http.StatusFound, redirectURL)
		return
	}

	// For mobile: Exchange code and store token immediately, then forward to mobile app
	// For web: Just redirect to frontend (frontend will handle the rest)
	if oauthState.Platform == "mobile" {
		ctx := context.Background()
		token, err := h.tokenService.HandleOAuth2Callback(ctx, code)
		if err != nil {
			redirectURL := h.buildRedirectURL(oauthState, "exchange_failed", "")
			c.Redirect(http.StatusFound, redirectURL)
			return
		}

		if err := h.tokenService.StoreToken(oauthState.UserID, token); err != nil {
			redirectURL := h.buildRedirectURL(oauthState, "failed_to_store_token", "")
			c.Redirect(http.StatusFound, redirectURL)
			return
		}

		// Redirect to mobile app with success
		redirectURL := h.buildRedirectURL(oauthState, "", "success")
		c.Redirect(http.StatusFound, redirectURL)
		return
	}

	// Web flow: Redirect to frontend with code, frontend will exchange
	frontendURL := getFrontendURL(c)
	redirectURL := frontendURL + "/google-calendar/callback?code=" + code + "&state=" + stateStr
	c.Redirect(http.StatusFound, redirectURL)
}

// buildRedirectURL builds the redirect URL based on platform
func (h *GoogleCalendarAuthHandler) buildRedirectURL(state OAuthState, errorCode string, success string) string {
	if state.Platform == "mobile" {
		// Mobile deep link
		if errorCode != "" {
			return "crmhealth://google-calendar/callback?error=" + errorCode + "&state=" + state.UserID
		}
		return "crmhealth://google-calendar/callback?success=true&state=" + state.UserID
	}

	// Web frontend URL — dari GOOGLE_CALENDAR_FRONTEND_URL env var (via config)
	frontendURL := h.config.FrontendURL
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}
	if errorCode != "" {
		return frontendURL + "/google-calendar/callback?error=" + errorCode + "&state=" + state.UserID
	}
	return frontendURL + "/google-calendar/callback?success=true&state=" + state.UserID
}

// GetConnectionStatus checks if user has connected Google Calendar
func (h *GoogleCalendarAuthHandler) GetConnectionStatus(c *gin.Context) {
	userID := ""
	if userIDVal, exists := c.Get("user_id"); exists {
		if id, ok := userIDVal.(string); ok {
			userID = id
		}
	}

	if userID == "" {
		errors.ErrorResponse(c, "UNAUTHORIZED", nil, nil)
		return
	}

	_, err := h.tokenService.GetToken(userID)
	isConnected := err == nil

	response.SuccessResponse(c, map[string]interface{}{
		"connected": isConnected,
	}, nil)
}

// Disconnect removes Google Calendar token
func (h *GoogleCalendarAuthHandler) Disconnect(c *gin.Context) {
	// Get user ID from context
	userID := ""
	if userIDVal, exists := c.Get("user_id"); exists {
		if id, ok := userIDVal.(string); ok {
			userID = id
		}
	}

	if userID == "" {
		errors.ErrorResponse(c, "UNAUTHORIZED", nil, nil)
		return
	}

	if err := h.tokenService.DeleteToken(userID); err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, map[string]interface{}{
		"message": "Google Calendar disconnected successfully",
	}, nil)
}

// getFrontendURL extracts frontend URL from request headers
func getFrontendURL(c *gin.Context) string {
	// Allow explicit override via env var for reliability in OAuth redirect
	if env := os.Getenv("FRONTEND_URL"); env != "" {
		return env
	}

	// Prefer Origin header when present (usually the safest)
	frontendURL := c.GetHeader("Origin")
	if frontendURL != "" {
		return frontendURL
	}

	// Fall back to Referer, but ignore referers originating from Google accounts
	// because the OAuth callback request comes from Google and will set their domain.
	frontendURL = c.GetHeader("Referer")
	if frontendURL != "" {
		parts := strings.Split(frontendURL, "/")
		if len(parts) >= 3 {
			candidate := strings.Join(parts[:3], "/")
			// If referer is a Google domain, ignore and fallback to default
			lower := strings.ToLower(candidate)
			if strings.Contains(lower, "accounts.google.com") || strings.Contains(lower, "googleusercontent.com") {
				// ignore
			} else {
				return candidate
			}
		}
	}

	// Final default for local development
	return "http://localhost:3000"
}

// ExchangeCode exchanges authorization code for token (Mobile Option 2)
// This endpoint is used by mobile apps to exchange the authorization code received from Google
func (h *GoogleCalendarAuthHandler) ExchangeCode(c *gin.Context) {
	// Get user ID from context
	userID := ""
	if userIDVal, exists := c.Get("user_id"); exists {
		if id, ok := userIDVal.(string); ok {
			userID = id
		}
	}

	if userID == "" {
		errors.ErrorResponse(c, "UNAUTHORIZED", nil, nil)
		return
	}

	var req struct {
		Code  string `json:"code" binding:"required"`
		State string `json:"state" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		errors.ErrorResponse(c, "INVALID_REQUEST", map[string]interface{}{"error": err.Error()}, nil)
		return
	}

	// Validate state parameter
	oauthState, err := decodeState(req.State)
	if err != nil {
		errors.ErrorResponse(c, "INVALID_STATE", map[string]interface{}{"message": "Invalid state parameter"}, nil)
		return
	}

	// Verify user ID matches
	if oauthState.UserID != userID {
		errors.UnauthorizedResponse(c, "USER_MISMATCH")
		return
	}

	// Exchange code for token
	ctx := context.Background()
	token, err := h.tokenService.ExchangeCode(ctx, req.Code, oauthState.Platform)
	if err != nil {
		log.Printf("[Google Calendar Handler] ExchangeCode - ERROR: %v", err)
		errors.InternalServerErrorResponse(c, "Failed to exchange authorization code")
		return
	}

	// Store token
	if err := h.tokenService.StoreToken(userID, token); err != nil {
		log.Printf("[Google Calendar Handler] ExchangeCode - ERROR storing token: %v", err)
		errors.InternalServerErrorResponse(c, "Failed to store token")
		return
	}

	response.SuccessResponse(c, map[string]interface{}{
		"connected": true,
		"message":   "Google Calendar connected successfully",
	}, nil)
}
