package services

import (
	"errors"

	"github.com/Sistal/bff-orchestrator/internal/clients"
	"github.com/Sistal/bff-orchestrator/internal/models"
)

type AuthService interface {
	Validate(token string) (*models.AuthValidateResponse, error)
	GetMe(token string) (*models.AuthMeResponse, error)
}

type MockAuthService struct{}

func NewMockAuthService() AuthService {
	return &MockAuthService{}
}

func (s *MockAuthService) Validate(token string) (*models.AuthValidateResponse, error) {
	if token == "invalid" {
		return nil, errors.New("invalid token")
	}
	return &models.AuthValidateResponse{
		Valid:    true,
		UserID:   123,
		Username: "juan.perez",
		Role:     1,
	}, nil
}

func (s *MockAuthService) GetMe(token string) (*models.AuthMeResponse, error) {
	return &models.AuthMeResponse{
		ID:       123,
		Username: "juan.perez",
		Role:     1,
	}, nil
}

type HTTPAuthService struct {
	client *clients.IdentityClient
}

func NewHTTPAuthService(client *clients.IdentityClient) AuthService {
	return &HTTPAuthService{client: client}
}

func (s *HTTPAuthService) Validate(token string) (*models.AuthValidateResponse, error) {
	return s.client.ValidateToken(token)
}

func (s *HTTPAuthService) GetMe(token string) (*models.AuthMeResponse, error) {
	// The identity service doesn't have /auth/me, but /auth/validate returns the user info
	val, err := s.client.ValidateToken(token)
	if err != nil {
		return nil, err
	}
	if !val.Valid {
		return nil, errors.New("token invalid or expired")
	}

	return &models.AuthMeResponse{
		ID:       val.UserID,
		Username: val.Username,
		Role:     val.Role,
	}, nil
}
