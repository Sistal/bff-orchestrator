package services

import (
	"net/url"

	"github.com/Sistal/bff-orchestrator/internal/clients"
	"github.com/Sistal/bff-orchestrator/internal/models"
)

// AuthService define las operaciones expuestas por el BFF de autenticación.
type AuthService interface {
	// Públicos
	Login(nombreUsuario, password string) (int, *models.APIResponse, error)
	Register(req *models.RegisterRequest) (int, *models.APIResponse, error)
	Validate(token string) (int, *models.APIResponse, error)
	Refresh(refreshToken string) (int, *models.APIResponse, error)

	// Protegidos (requieren JWT)
	GetMe(token string) (int, *models.APIResponse, error)
	Logout(token string) (int, *models.APIResponse, error)
	ChangePassword(token string, req *models.ChangePasswordRequest) (int, *models.APIResponse, error)
	GetRoles(token string, activosSolo *bool) (int, *models.APIResponse, error)

	// Administración
	CreateUser(token string, req *models.RegisterRequest) (int, *models.APIResponse, error)
	ListUsers(token string, params url.Values) (int, *models.APIResponse, error)
	GetUserByID(token string, id int) (int, *models.APIResponse, error)
	UpdateUser(token string, id int, req *models.UpdateUsuarioRequest) (int, *models.APIResponse, error)
}

// HTTPAuthService delega todas las operaciones al ms-authentication.
type HTTPAuthService struct {
	client *clients.IdentityClient
}

func NewHTTPAuthService(client *clients.IdentityClient) AuthService {
	return &HTTPAuthService{client: client}
}

func (s *HTTPAuthService) Login(nombreUsuario, password string) (int, *models.APIResponse, error) {
	return s.client.Login(map[string]string{
		"nombre_usuario": nombreUsuario,
		"password":       password,
	})
}

func (s *HTTPAuthService) Register(req *models.RegisterRequest) (int, *models.APIResponse, error) {
	return s.client.Register(req)
}

func (s *HTTPAuthService) Validate(token string) (int, *models.APIResponse, error) {
	return s.client.ValidateTokenRaw(token)
}

func (s *HTTPAuthService) Refresh(refreshToken string) (int, *models.APIResponse, error) {
	return s.client.Refresh(map[string]string{
		"refresh_token": refreshToken,
	})
}

func (s *HTTPAuthService) GetMe(token string) (int, *models.APIResponse, error) {
	return s.client.GetMe(token)
}

func (s *HTTPAuthService) Logout(token string) (int, *models.APIResponse, error) {
	return s.client.Logout(token)
}

func (s *HTTPAuthService) ChangePassword(token string, req *models.ChangePasswordRequest) (int, *models.APIResponse, error) {
	return s.client.ChangePassword(token, req)
}

func (s *HTTPAuthService) GetRoles(token string, activosSolo *bool) (int, *models.APIResponse, error) {
	return s.client.GetRoles(token, activosSolo)
}

func (s *HTTPAuthService) CreateUser(token string, req *models.RegisterRequest) (int, *models.APIResponse, error) {
	return s.client.CreateUser(token, req)
}

func (s *HTTPAuthService) ListUsers(token string, params url.Values) (int, *models.APIResponse, error) {
	return s.client.ListUsers(token, params)
}

func (s *HTTPAuthService) GetUserByID(token string, id int) (int, *models.APIResponse, error) {
	return s.client.GetUserByID(token, id)
}

func (s *HTTPAuthService) UpdateUser(token string, id int, req *models.UpdateUsuarioRequest) (int, *models.APIResponse, error) {
	return s.client.UpdateUser(token, id, req)
}
