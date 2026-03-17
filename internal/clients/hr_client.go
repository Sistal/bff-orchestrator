package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/Sistal/bff-orchestrator/internal/logger"
	"github.com/Sistal/bff-orchestrator/internal/models"
	"go.uber.org/zap"
)

type HRClient struct {
	BaseURL string
	Client  *http.Client
}

func NewHRClient() *HRClient {
	baseURL := os.Getenv("MS_FUNCIONARIO_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8082" // Default port for HR Service
	}
	return &HRClient{
		BaseURL: baseURL,
		Client:  &http.Client{},
	}
}

// GetEmployeeByUserID resuelve id_usuario → id_funcionario consultando
// GET /api/v1/funcionarios/by-usuario/{userID} en el ms-funcionario.
// Es llamado por el BearerAuthMiddleware en cada request autenticado.
// DEUDA TÉCNICA: considerar caché con TTL para reducir latencia.
func (c *HRClient) GetEmployeeByUserID(userID string) (int, error) {
	resp, err := c.GetFullEmployeeByUserID(userID)
	if err != nil {
		return 0, err
	}
	if resp == nil {
		return 0, nil
	}
	return resp.IDFuncionario, nil
}

// GetFullEmployeeByUserID devuelve el registro completo del funcionario usando su id usuario
// consultando el mismo endpoint. Útil para validar requerimientos de completitud de datos.
func (c *HRClient) GetFullEmployeeByUserID(userID string) (*models.EmployeeIDResponse, error) {
	resp, err := c.Client.Get(fmt.Sprintf("%s/api/v1/funcionarios/by-usuario/%s", c.BaseURL, userID))
	if err != nil {
		return nil, fmt.Errorf("hr_client: error al contactar ms-funcionario: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil // usuario sin funcionario asociado (admin puro)
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("hr_client: GetEmployeeByUserID status %d: %s", resp.StatusCode, string(body))
	}

	// Estructura para manejar el envelope { "success": true, "data": { "id_funcionario": ... } }
	var envelope struct {
		Success bool                      `json:"success"`
		Data    models.EmployeeIDResponse `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		return nil, fmt.Errorf("hr_client: error al decodificar respuesta: %w", err)
	}
	return &envelope.Data, nil
}

func (c *HRClient) GetEmployeeProfile(id string) (*models.EmployeeProfile, error) {
	log := logger.Get()
	url := fmt.Sprintf("%s/api/v1/funcionarios/%s", c.BaseURL, id)
	log.Debug("HRClient.GetEmployeeProfile: Llamando a ms-funcionario", zap.String("url", url))

	resp, err := c.Client.Get(url)
	if err != nil {
		log.Error("HRClient.GetEmployeeProfile: Error en request HTTP", zap.String("id", id), zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Warn("HRClient.GetEmployeeProfile: Respuesta no exitosa", zap.String("id", id), zap.Int("status_code", resp.StatusCode))
		return nil, fmt.Errorf("failed to get employee profile: %d", resp.StatusCode)
	}

	var profile models.EmployeeProfile
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		log.Error("HRClient.GetEmployeeProfile: Error decodificando respuesta", zap.String("id", id), zap.Error(err))
		return nil, err
	}

	log.Debug("HRClient.GetEmployeeProfile: Perfil obtenido exitosamente", zap.String("id", id))
	return &profile, nil
}

func (c *HRClient) UpdateContact(id string, req models.UpdateContactRequest) error {
	body, _ := json.Marshal(req)
	reqHTTP, _ := http.NewRequest("PUT", fmt.Sprintf("%s/api/v1/funcionarios/%s/contacto", c.BaseURL, id), bytes.NewBuffer(body))
	reqHTTP.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(reqHTTP)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update contact: %d", resp.StatusCode)
	}
	return nil
}

func (c *HRClient) GetMeasurements(id string) (*models.BodyMeasurements, error) {
	resp, err := c.Client.Get(fmt.Sprintf("%s/api/v1/funcionarios/%s/medidas", c.BaseURL, id))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get measurements: %d", resp.StatusCode)
	}

	var m models.BodyMeasurements
	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return nil, err
	}
	return &m, nil
}

func (c *HRClient) RegisterMeasurements(id string, req models.BodyMeasurements) (*models.BodyMeasurements, error) {
	body, _ := json.Marshal(req)
	reqHTTP, _ := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/funcionarios/%s/medidas", c.BaseURL, id), bytes.NewBuffer(body))
	reqHTTP.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(reqHTTP)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("failed to register measurements: %d", resp.StatusCode)
	}

	var m models.BodyMeasurements
	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return nil, err
	}
	return &m, nil
}

func (c *HRClient) UpdateMeasurements(id string, req models.BodyMeasurements) (*models.BodyMeasurements, error) {
	body, _ := json.Marshal(req)
	reqHTTP, _ := http.NewRequest("PUT", fmt.Sprintf("%s/api/v1/funcionarios/%s/medidas", c.BaseURL, id), bytes.NewBuffer(body))
	reqHTTP.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(reqHTTP)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to update measurements: %d", resp.StatusCode)
	}

	var m models.BodyMeasurements
	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return nil, err
	}
	return &m, nil
}

func (c *HRClient) UpdateSecurity(id string, req models.UpdateSecurityRequest) error {
	body, _ := json.Marshal(req)
	reqHTTP, _ := http.NewRequest("PUT", fmt.Sprintf("%s/api/v1/funcionarios/%s/seguridad", c.BaseURL, id), bytes.NewBuffer(body))
	reqHTTP.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(reqHTTP)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update security: %d", resp.StatusCode)
	}
	return nil
}

func (c *HRClient) GetBranches() ([]models.Branch, error) {
	resp, err := c.Client.Get(c.BaseURL + "/api/v1/sucursales")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get branches: %d", resp.StatusCode)
	}

	var branches []models.Branch
	if err := json.NewDecoder(resp.Body).Decode(&branches); err != nil {
		return nil, err
	}
	return branches, nil
}

func (c *HRClient) GetChangeHistory(employeeID string) ([]models.BranchChangeRequestHistory, error) {
	resp, err := c.Client.Get(fmt.Sprintf("%s/api/v1/transferencias?id_funcionario=%s", c.BaseURL, employeeID))
	if err != nil {
		return nil, fmt.Errorf("failed to get change history: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get change history: %d", resp.StatusCode)
	}

	var history []models.BranchChangeRequestHistory
	if err := json.NewDecoder(resp.Body).Decode(&history); err != nil {
		return nil, err
	}
	return history, nil
}

func (c *HRClient) CreateBranchChangeRequest(employeeID string, req models.CreateBranchChangeRequest) (*models.BranchChangeRequestHistory, error) {
	// Enriquecer el body con el id_funcionario resuelto por el middleware
	type branchChangePayload struct {
		models.CreateBranchChangeRequest
		IDFuncionario string `json:"id_funcionario"`
	}
	payload := branchChangePayload{
		CreateBranchChangeRequest: req,
		IDFuncionario:             employeeID,
	}
	body, _ := json.Marshal(payload)
	reqHTTP, _ := http.NewRequest("POST", c.BaseURL+"/api/v1/transferencias", bytes.NewBuffer(body))
	reqHTTP.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(reqHTTP)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("failed to create branch change request: %d", resp.StatusCode)
	}

	var history models.BranchChangeRequestHistory
	if err := json.NewDecoder(resp.Body).Decode(&history); err != nil {
		return nil, err
	}
	return &history, nil
}

func (c *HRClient) CreateEmployee(req models.CreateEmployeeRequest) (*models.EmployeeProfile, error) {
	// Se asume que userid ya viene en el req, o se podría pasar como argumento separado si se prefiere.
	// La implementación actual del modelo incluye UserID en CreateEmployeeRequest.

	body, _ := json.Marshal(req)
	// Asumimos que el endpoint en ms-funcionario es POST /api/v1/funcionarios
	reqHTTP, _ := http.NewRequest("POST", c.BaseURL+"/api/v1/funcionarios", bytes.NewBuffer(body))
	reqHTTP.Header.Set("Content-Type", "application/json")
	// Opcional: Si el ms-funcionario requiere el User-ID en header para auditoría
	reqHTTP.Header.Set("X-User-ID", fmt.Sprintf("%d", req.UserID))

	resp, err := c.Client.Do(reqHTTP)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to create employee: %d", resp.StatusCode)
	}

	var newEmployee models.EmployeeProfile
	if err := json.NewDecoder(resp.Body).Decode(&newEmployee); err != nil {
		return nil, err
	}
	return &newEmployee, nil
}

func (c *HRClient) CreateEmployeeFromLogin(req models.CreateEmployeeRequestFromLogin) (*models.EmployeeProfile, error) {
	body, _ := json.Marshal(req)
	// Asumimos que el endpoint en ms-funcionario para registrar desde login es POST /api/v1/funcionarios/register
	reqHTTP, _ := http.NewRequest("POST", c.BaseURL+"/api/v1/funcionarios/register", bytes.NewBuffer(body))
	reqHTTP.Header.Set("Content-Type", "application/json")
	// Opcional: Si el ms-funcionario requiere el User-ID en header para auditoría
	reqHTTP.Header.Set("X-User-ID", fmt.Sprintf("%d", req.UserID))

	resp, err := c.Client.Do(reqHTTP)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to create employee from login: %d", resp.StatusCode)
	}

	var newEmployee models.EmployeeProfile
	if err := json.NewDecoder(resp.Body).Decode(&newEmployee); err != nil {
		return nil, err
	}
	return &newEmployee, nil
}
