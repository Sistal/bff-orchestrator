package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/Sistal/bff-orchestrator/internal/models"
)

type IdentityClient struct {
	BaseURL string
	Client  *http.Client
}

func NewIdentityClient() *IdentityClient {
	baseURL := os.Getenv("MS_AUTHENTICATION_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8081"
	}
	return &IdentityClient{
		BaseURL: baseURL,
		Client:  &http.Client{},
	}
}

// ─── helpers ─────────────────────────────────────────────────────────────────

// doJSON ejecuta la petición y decodifica el body en dest. Retorna el código HTTP.
func (c *IdentityClient) doJSON(req *http.Request, dest interface{}) (int, error) {
	resp, err := c.Client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, err
	}
	if dest != nil && len(body) > 0 {
		if err := json.Unmarshal(body, dest); err != nil {
			return resp.StatusCode, err
		}
	}
	return resp.StatusCode, nil
}

func (c *IdentityClient) newJSONReq(method, path string, body interface{}) (*http.Request, error) {
	var buf *bytes.Buffer
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		buf = bytes.NewBuffer(b)
	} else {
		buf = bytes.NewBuffer(nil)
	}
	req, err := http.NewRequest(method, c.BaseURL+path, buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

// ─── Auth endpoints ───────────────────────────────────────────────────────────

// Login llama POST /api/v1/auth/login y retorna el body completo de APIResponse.
func (c *IdentityClient) Login(payload map[string]string) (int, *models.APIResponse, error) {
	req, err := c.newJSONReq("POST", "/api/v1/auth/login", payload)
	if err != nil {
		return 0, nil, err
	}
	var result models.APIResponse
	code, err := c.doJSON(req, &result)
	return code, &result, err
}

// Register llama POST /api/v1/auth/register.
func (c *IdentityClient) Register(payload interface{}) (int, *models.APIResponse, error) {
	req, err := c.newJSONReq("POST", "/api/v1/auth/register", payload)
	if err != nil {
		return 0, nil, err
	}
	var result models.APIResponse
	code, err := c.doJSON(req, &result)
	return code, &result, err
}

// ValidateToken llama GET /api/v1/auth/validate con el Bearer token.
// Sigue siendo usado por el middleware de autenticación.
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

	if resp.StatusCode == http.StatusUnauthorized {
		return &models.AuthValidateResponse{Valid: false}, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("validation failed with status: %d", resp.StatusCode)
	}

	// El ms-authentication devuelve APIResponse con data = ValidateTokenData.
	// Mapeamos al AuthValidateResponse legacy que usa el middleware.
	var envelope models.APIResponse
	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, err
	}
	if !envelope.Success {
		return &models.AuthValidateResponse{Valid: false}, nil
	}

	// Extraer data como ValidateTokenData
	dataBytes, _ := json.Marshal(envelope.Data)
	var vd models.ValidateTokenData
	if err := json.Unmarshal(dataBytes, &vd); err != nil {
		return nil, err
	}

	return &models.AuthValidateResponse{
		Valid:          true,
		UserID:         vd.IDUsuario,
		Username:       vd.NombreUsuario,
		NombreCompleto: vd.NombreCompleto,
		RUT:            vd.RUT,
		Role:           vd.IDRol,
		RoleName:       vd.NombreRol,
		EstadoID:       vd.IDEstadoUsuario,
		EstadoName:     vd.NombreEstado,
		IssuedAt:       vd.Iat,
		ExpiresAt:      vd.Exp,
	}, nil
}

// ValidateTokenRaw llama GET /api/v1/auth/validate y retorna el APIResponse completo.
func (c *IdentityClient) ValidateTokenRaw(token string) (int, *models.APIResponse, error) {
	req, err := http.NewRequest("GET", c.BaseURL+"/api/v1/auth/validate", nil)
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	var result models.APIResponse
	code, err := c.doJSON(req, &result)
	return code, &result, err
}

// Refresh llama POST /api/v1/auth/refresh.
func (c *IdentityClient) Refresh(payload interface{}) (int, *models.APIResponse, error) {
	req, err := c.newJSONReq("POST", "/api/v1/auth/refresh", payload)
	if err != nil {
		return 0, nil, err
	}
	var result models.APIResponse
	code, err := c.doJSON(req, &result)
	return code, &result, err
}

// Logout llama POST /api/v1/auth/logout con el Bearer token.
func (c *IdentityClient) Logout(token string) (int, *models.APIResponse, error) {
	req, err := http.NewRequest("POST", c.BaseURL+"/api/v1/auth/logout", nil)
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	var result models.APIResponse
	code, err := c.doJSON(req, &result)
	return code, &result, err
}

// ChangePassword llama PUT /api/v1/auth/change-password.
func (c *IdentityClient) ChangePassword(token string, payload interface{}) (int, *models.APIResponse, error) {
	req, err := c.newJSONReq("PUT", "/api/v1/auth/change-password", payload)
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	var result models.APIResponse
	code, err := c.doJSON(req, &result)
	return code, &result, err
}

// GetMe llama GET /api/v1/auth/me.
func (c *IdentityClient) GetMe(token string) (int, *models.APIResponse, error) {
	req, err := http.NewRequest("GET", c.BaseURL+"/api/v1/auth/me", nil)
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	var result models.APIResponse
	code, err := c.doJSON(req, &result)
	return code, &result, err
}

// GetRoles llama GET /api/v1/auth/roles.
func (c *IdentityClient) GetRoles(token string, activosSolo *bool) (int, *models.APIResponse, error) {
	u, _ := url.Parse(c.BaseURL + "/api/v1/auth/roles")
	if activosSolo != nil {
		q := u.Query()
		q.Set("activos_solo", strconv.FormatBool(*activosSolo))
		u.RawQuery = q.Encode()
	}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	var result models.APIResponse
	code, err := c.doJSON(req, &result)
	return code, &result, err
}

// ─── Users admin endpoints ────────────────────────────────────────────────────

// CreateUser llama POST /api/v1/auth/users (admin).
func (c *IdentityClient) CreateUser(token string, payload interface{}) (int, *models.APIResponse, error) {
	req, err := c.newJSONReq("POST", "/api/v1/auth/users", payload)
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	var result models.APIResponse
	code, err := c.doJSON(req, &result)
	return code, &result, err
}

// ListUsers llama GET /api/v1/auth/users con query params opcionales.
func (c *IdentityClient) ListUsers(token string, params url.Values) (int, *models.APIResponse, error) {
	u, _ := url.Parse(c.BaseURL + "/api/v1/auth/users")
	u.RawQuery = params.Encode()
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	var result models.APIResponse
	code, err := c.doJSON(req, &result)
	return code, &result, err
}

// GetUserByID llama GET /api/v1/auth/users/:id.
func (c *IdentityClient) GetUserByID(token string, id int) (int, *models.APIResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/auth/users/%d", c.BaseURL, id), nil)
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	var result models.APIResponse
	code, err := c.doJSON(req, &result)
	return code, &result, err
}

// UpdateUser llama PUT /api/v1/auth/users/:id.
func (c *IdentityClient) UpdateUser(token string, id int, payload interface{}) (int, *models.APIResponse, error) {
	req, err := c.newJSONReq("PUT", fmt.Sprintf("/api/v1/auth/users/%d", id), payload)
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	var result models.APIResponse
	code, err := c.doJSON(req, &result)
	return code, &result, err
}
