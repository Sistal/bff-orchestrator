package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/Sistal/bff-orchestrator/internal/models"
)

type OpsClient struct {
	BaseURL string
	Client  *http.Client
}

func NewOpsClient() *OpsClient {
	baseURL := os.Getenv("MS_OPERATIONS_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8083"
	}
	return &OpsClient{
		BaseURL: baseURL,
		Client:  &http.Client{},
	}
}

// ─── Request Operations ────────────────────────────────────────────────────────

func (c *OpsClient) GetRequests(employeeID string, params map[string]string) ([]models.RequestSummary, error) {
	u, _ := url.Parse(fmt.Sprintf("%s/api/v1/peticiones", c.BaseURL))
	q := u.Query()
	q.Set("id_funcionario", employeeID)
	for k, v := range params {
		if v != "" {
			q.Set(k, v)
		}
	}
	u.RawQuery = q.Encode()

	resp, err := c.Client.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get requests: %d", resp.StatusCode)
	}
	var requests []models.RequestSummary
	if err := json.NewDecoder(resp.Body).Decode(&requests); err != nil {
		return nil, err
	}
	return requests, nil
}

func (c *OpsClient) GetRequestByID(employeeID, id string) (*models.RequestSummary, error) {
	reqHTTP, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/peticiones/%s", c.BaseURL, id), nil)
	reqHTTP.Header.Set("X-User-ID", employeeID)
	resp, err := c.Client.Do(reqHTTP)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("request not found: %s", id)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get request: %d", resp.StatusCode)
	}
	var request models.RequestSummary
	if err := json.NewDecoder(resp.Body).Decode(&request); err != nil {
		return nil, err
	}
	return &request, nil
}

func (c *OpsClient) GetRecentRequests(employeeID string) ([]models.RequestSummary, error) {
	reqHTTP, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/peticiones/recientes?id_funcionario=%s", c.BaseURL, employeeID), nil)
	resp, err := c.Client.Do(reqHTTP)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get recent requests: %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// 1. Intentar decodificar como Envelope
	var envelope struct {
		Success bool                    `json:"success"`
		Data    []models.RequestSummary `json:"data"`
	}
	if err := json.Unmarshal(bodyBytes, &envelope); err == nil {
		return envelope.Data, nil
	}

	// 2. Intentar decodificar como Array directo
	var requests []models.RequestSummary
	if err := json.Unmarshal(bodyBytes, &requests); err == nil {
		return requests, nil
	}

	return nil, fmt.Errorf("failed to decode response as envelope or array")
}

func (c *OpsClient) CreateReplenishmentRequest(employeeID string, req models.CreateReplenishmentRequest) (*models.RequestSummary, error) {
	body, _ := json.Marshal(req)
	reqHTTP, _ := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/peticiones/reposicion", c.BaseURL), bytes.NewBuffer(body))
	reqHTTP.Header.Set("Content-Type", "application/json")
	reqHTTP.Header.Set("X-User-ID", employeeID)
	resp, err := c.Client.Do(reqHTTP)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to create replenishment: %d", resp.StatusCode)
	}
	var created models.RequestSummary
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		return nil, err
	}
	return &created, nil
}

func (c *OpsClient) CreateGarmentChangeRequest(employeeID string, req models.CreateGarmentChangeRequest) (*models.RequestSummary, error) {
	body, _ := json.Marshal(req)
	reqHTTP, _ := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/peticiones/cambio-prenda", c.BaseURL), bytes.NewBuffer(body))
	reqHTTP.Header.Set("Content-Type", "application/json")
	reqHTTP.Header.Set("X-User-ID", employeeID)
	resp, err := c.Client.Do(reqHTTP)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to create garment change: %d", resp.StatusCode)
	}
	var created models.RequestSummary
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		return nil, err
	}
	return &created, nil
}

func (c *OpsClient) CreateUniformRequest(employeeID string, req models.CreateUniformRequest) (*models.RequestSummary, error) {
	// Convertir employeeID a int
	id, err := strconv.Atoi(employeeID)
	if err != nil {
		return nil, fmt.Errorf("invalid employee ID (must be int): %v", err)
	}

	// Payload enriquecido con el id_funcionario
	payload := struct {
		models.CreateUniformRequest
		IDFuncionario int `json:"id_funcionario"`
	}{
		CreateUniformRequest: req,
		IDFuncionario:        id,
	}

	body, _ := json.Marshal(payload)
	reqHTTP, _ := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/peticiones/uniforme", c.BaseURL), bytes.NewBuffer(body))
	reqHTTP.Header.Set("Content-Type", "application/json")
	// Se envía también en header por compatibilidad con otros endpoints,
	// pero el body lleva el dato persitible.
	reqHTTP.Header.Set("X-User-ID", employeeID)

	resp, err := c.Client.Do(reqHTTP)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to create uniform request: %d", resp.StatusCode)
	}

	var created models.RequestSummary
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		return nil, err
	}
	return &created, nil
}

func (c *OpsClient) UploadFile(employeeID string, file *multipart.FileHeader) (*models.FileUploadResponse, error) {
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, err := writer.CreateFormFile("file", file.Filename)
	if err != nil {
		return nil, err
	}

	fileBytes := make([]byte, file.Size)
	if _, err := src.Read(fileBytes); err != nil {
		return nil, err
	}
	part.Write(fileBytes)
	writer.Close()

	reqHTTP, _ := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/archivos/upload", c.BaseURL), &buf)
	reqHTTP.Header.Set("Content-Type", writer.FormDataContentType())
	reqHTTP.Header.Set("X-User-ID", employeeID)

	resp, err := c.Client.Do(reqHTTP)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to upload file: %d", resp.StatusCode)
	}
	var result models.FileUploadResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ─── Delivery Operations ───────────────────────────────────────────────────────

func (c *OpsClient) GetDeliveries(employeeID string) ([]models.DeliverySummary, error) {
	reqHTTP, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/despachos?id_funcionario=%s", c.BaseURL, employeeID), nil)
	resp, err := c.Client.Do(reqHTTP)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get deliveries: %d", resp.StatusCode)
	}
	var deliveries []models.DeliverySummary
	if err := json.NewDecoder(resp.Body).Decode(&deliveries); err != nil {
		return nil, err
	}
	return deliveries, nil
}

func (c *OpsClient) GetDeliveryByID(employeeID, id string) (*models.DeliverySummary, error) {
	reqHTTP, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/despachos/%s", c.BaseURL, id), nil)
	reqHTTP.Header.Set("X-User-ID", employeeID)
	resp, err := c.Client.Do(reqHTTP)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("delivery not found: %s", id)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get delivery: %d", resp.StatusCode)
	}
	var delivery models.DeliverySummary
	if err := json.NewDecoder(resp.Body).Decode(&delivery); err != nil {
		return nil, err
	}
	return &delivery, nil
}

func (c *OpsClient) ConfirmDelivery(employeeID, id string) (*models.DeliverySummary, error) {
	reqHTTP, _ := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/despachos/%s/confirmar", c.BaseURL, id), nil)
	reqHTTP.Header.Set("X-User-ID", employeeID)
	resp, err := c.Client.Do(reqHTTP)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to confirm delivery: %d", resp.StatusCode)
	}
	var delivery models.DeliverySummary
	if err := json.NewDecoder(resp.Body).Decode(&delivery); err != nil {
		return nil, err
	}
	return &delivery, nil
}

func (c *OpsClient) GetUpcomingDeliveries(employeeID string) ([]models.DeliverySummary, error) {
	reqHTTP, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/despachos/upcoming?id_funcionario=%s", c.BaseURL, employeeID), nil)
	resp, err := c.Client.Do(reqHTTP)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get upcoming deliveries: %d", resp.StatusCode)
	}
	var deliveries []models.DeliverySummary
	if err := json.NewDecoder(resp.Body).Decode(&deliveries); err != nil {
		return nil, err
	}
	return deliveries, nil
}

// ─── Notification Operations ───────────────────────────────────────────────────

func (c *OpsClient) GetNotifications(employeeID string) ([]models.Notification, error) {
	reqHTTP, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/notificaciones?id_usuario=%s", c.BaseURL, employeeID), nil)
	resp, err := c.Client.Do(reqHTTP)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get notifications: %d", resp.StatusCode)
	}
	var notifs []models.Notification
	if err := json.NewDecoder(resp.Body).Decode(&notifs); err != nil {
		return nil, err
	}
	return notifs, nil
}

func (c *OpsClient) MarkAsRead(employeeID, id string) error {
	reqHTTP, _ := http.NewRequest("PATCH", fmt.Sprintf("%s/api/v1/notificaciones/%s/leida", c.BaseURL, id), nil)
	reqHTTP.Header.Set("X-User-ID", employeeID)
	resp, err := c.Client.Do(reqHTTP)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to mark notification as read: %d", resp.StatusCode)
	}
	return nil
}

func (c *OpsClient) MarkAllAsRead(employeeID string) (int, error) {
	reqHTTP, _ := http.NewRequest("PATCH", fmt.Sprintf("%s/api/v1/notificaciones/leer-todas", c.BaseURL), nil)
	reqHTTP.Header.Set("X-User-ID", employeeID)
	resp, err := c.Client.Do(reqHTTP)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to mark all notifications as read: %d", resp.StatusCode)
	}
	var result models.MarkAllReadResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}
	return result.Count, nil
}

// ─── Dashboard ─────────────────────────────────────────────────────────────────

func (c *OpsClient) GetDashboardStats(employeeID string) (*models.HomeStats, error) {
	reqHTTP, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/operaciones/dashboard?id_funcionario=%s", c.BaseURL, employeeID), nil)
	resp, err := c.Client.Do(reqHTTP)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get dashboard stats: %d", resp.StatusCode)
	}
	var stats models.HomeStats
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		return nil, err
	}
	return &stats, nil
}
