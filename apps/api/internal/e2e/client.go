package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"testing"
)

// Client is a wrapper around http.Client with specific helpers for the API
type Client struct {
	httpClient *http.Client
	baseURL    string
	UserID     string
	AccessToken string
}

// NewClient creates a new Client with a cookie jar
func NewClient(baseURL string) (*Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	return &Client{
		httpClient: &http.Client{
			Jar: jar,
		},
		baseURL: baseURL,
	}, nil
}

// Login attempts to log in and store the session cookies
func (c *Client) Login(email, password string) error {
	payload := map[string]string{
		"email":    email,
		"password": password,
	}
	
	resp, err := c.Post("/api/v1/auth/login", payload)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("login failed: status %d, body: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data struct {
			User struct {
				ID string `json:"id"`
			} `json:"user"`
			Token string `json:"token"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode login response: %v", err)
	}

	c.UserID = result.Data.User.ID
	c.AccessToken = result.Data.Token
	return nil
}

// Get performs a GET request
func (c *Client) Get(path string) (*http.Response, error) {
	return c.doRequest("GET", path, nil)
}

// Post performs a POST request
func (c *Client) Post(path string, body interface{}) (*http.Response, error) {
	return c.doRequest("POST", path, body)
}

// Put performs a PUT request
func (c *Client) Put(path string, body interface{}) (*http.Response, error) {
	return c.doRequest("PUT", path, body)
}

// Delete performs a DELETE request
func (c *Client) Delete(path string) (*http.Response, error) {
	return c.doRequest("DELETE", path, nil)
}

func (c *Client) doRequest(method, path string, body interface{}) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewBuffer(jsonBytes)
	}

	req, err := http.NewRequest(method, c.baseURL+path, bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	
	if c.AccessToken != "" {
		req.Header.Set("Authorization", "Bearer " + c.AccessToken)
	}

	// CSRF Handling
	c.addCSRFToken(req)

	return c.httpClient.Do(req)
}

// addCSRFToken looks for specific cookies and adds the X-CSRF-Token header
func (c *Client) addCSRFToken(req *http.Request) {
	u, _ := url.Parse(c.baseURL)
	cookies := c.httpClient.Jar.Cookies(u)

	for _, cookie := range cookies {
		if cookie.Name == "csrf_token" {
			req.Header.Set("X-CSRF-Token", cookie.Value)
			break
		}
	}
}

// ParseJSON helper to decode response body
func ParseJSON(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}

// GetAdminClient helper to get Admin Client
func GetAdminClient(t *testing.T) *Client {
	config := LoadConfig()
	client, err := NewClient(config.APIURL)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Try default admin credentials from ENV or Seeder defaults
	email := "superadmin@gilabs.id"
	password := "password" // Default fallback, but likely wrong if seeder used env

	// Use Config/Env if available (Config doesn't have it yet, need to load from os)
	// assuming LoadConfig called godotenv.Load
	if e := os.Getenv("DEFAULT_ADMIN_EMAIL"); e != "" {
		email = e
	}
	if p := os.Getenv("DEFAULT_ADMIN_PASSWORD"); p != "" {
		password = p
	}

	err = client.Login(email, password)
	if err != nil {
		t.Logf("Failed to login as default admin (%s): %v.", email, err)
		// Don't fail immediately, but subsequent tests might fail
	}
	return client
}
