package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/Sistal/bff-orchestrator/internal/models"
)

type IdentityClient struct {
	BaseURL string
	Client  *http.Client
}

func NewIdentityClient() *IdentityClient {
	baseURL := os.Getenv("MS_AUTHENTICATION_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8081" // Default port for Identity Service
	}
	return &IdentityClient{
		BaseURL: baseURL,
		Client:  &http.Client{},
	}
}

func (c *IdentityClient) Login(username, password string) (string, error) {
	reqBody, _ := json.Marshal(map[string]string{
		"username": username,
		"password": password,
	})

	resp, err := c.Client.Post(c.BaseURL+"/api/v1/auth/login", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("login failed with status: %d", resp.StatusCode)
	}

	var result map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result["token"], nil
}

func (c *IdentityClient) ValidateToken(token string) (*models.AuthValidateResponse, error) {
	req, err := http.NewRequest("GET", c.BaseURL+"/api/v1/auth/validate", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// If 401, it's just invalid
		if resp.StatusCode == http.StatusUnauthorized {
			return &models.AuthValidateResponse{Valid: false}, nil
		}
		return nil, fmt.Errorf("validation failed with status: %d", resp.StatusCode)
	}

	var result models.AuthValidateResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
