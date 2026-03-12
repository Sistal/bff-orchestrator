package handlers

import (
	"net/http"

	"github.com/Sistal/bff-orchestrator/internal/models"
	"github.com/Sistal/bff-orchestrator/internal/services"
	"github.com/gin-gonic/gin"
)

type EmployeeHandler struct {
	service services.EmployeeService
}

func NewEmployeeHandler(s services.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{service: s}
}

// GetProfile godoc
// @Summary      Get employee profile
// @Description  Get the current employee's profile
// @Tags         employees
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  models.EmployeeProfile
// @Failure      401  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /api/v1/funcionarios/me [get]
func (h *EmployeeHandler) GetProfile(c *gin.Context) {
	userID := c.GetString("userID")
	resp, err := h.service.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// UpdateContact godoc
// @Summary      Update employee contact
// @Description  Update the current employee's contact information
// @Tags         employees
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body      models.UpdateContactRequest  true  "Contact information"
// @Success      200      {object}  models.EmployeeProfile
// @Failure      400      {object}  models.ErrorResponse
// @Failure      401      {object}  models.ErrorResponse
// @Failure      500      {object}  models.ErrorResponse
// @Router       /api/v1/funcionarios/me [put]
func (h *EmployeeHandler) UpdateContact(c *gin.Context) {
	userID := c.GetString("userID")
	var req models.UpdateContactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": err.Error()})
		return
	}
	resp, err := h.service.UpdateContact(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// UpdatePreferences godoc
// @Summary      Update employee preferences
// @Description  Update the current employee's preferences (mock temporal, sin persistencia en DDL)
// @Tags         employees
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body      models.UpdatePreferencesRequest  true  "Preferences"
// @Success      200      {object}  models.UpdatePreferencesRequest
// @Failure      400      {object}  models.ErrorResponse
// @Failure      401      {object}  models.ErrorResponse
// @Failure      500      {object}  models.ErrorResponse
// @Router       /api/v1/funcionarios/me/preferencias [put]
func (h *EmployeeHandler) UpdatePreferences(c *gin.Context) {
	userID := c.GetString("userID")
	var req models.UpdatePreferencesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": err.Error()})
		return
	}
	resp, err := h.service.UpdatePreferences(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Error", "message": err.Error()})
		return
	}
	// Retornar el objeto actualizado como espera el frontend
	c.JSON(http.StatusOK, resp)
}

// UpdateSecurity godoc
// @Summary      Update employee security settings
// @Description  Update the current employee's security settings (recoveryEmail)
// @Tags         employees
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body      models.UpdateSecurityRequest  true  "Security settings"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  models.ErrorResponse
// @Failure      401      {object}  models.ErrorResponse
// @Failure      500      {object}  models.ErrorResponse
// @Router       /api/v1/funcionarios/me/seguridad [put]
func (h *EmployeeHandler) UpdateSecurity(c *gin.Context) {
	userID := c.GetString("userID")
	var req models.UpdateSecurityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": err.Error()})
		return
	}
	if err := h.service.UpdateSecurity(userID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Configuración de seguridad actualizada"})
}

// GetStats godoc
// @Summary      Get employee statistics
// @Description  Get statistics for the current employee (dashboard)
// @Tags         employees
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  models.HomeStats
// @Failure      401  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /api/v1/funcionarios/me/stats [get]
func (h *EmployeeHandler) GetStats(c *gin.Context) {
	userID := c.GetString("userID")
	resp, err := h.service.GetStats(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// GetActivity godoc
// @Summary      Get employee activity
// @Description  Get activity history for the current employee
// @Tags         employees
// @Security     BearerAuth
// @Produce      json
// @Success      200  {array}   models.ActivityLog
// @Failure      401  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /api/v1/funcionarios/me/actividad [get]
func (h *EmployeeHandler) GetActivity(c *gin.Context) {
	userID := c.GetString("userID")
	resp, err := h.service.GetActivity(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// GetMeasurements godoc
// @Summary      Get employee measurements
// @Description  Get body measurements for a specific employee
// @Tags         employees
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "Employee ID"
// @Success      200  {object}  models.BodyMeasurements
// @Failure      401  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /api/v1/funcionarios/{id}/medidas [get]
func (h *EmployeeHandler) GetMeasurements(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.service.GetMeasurements(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// RegisterMeasurements godoc
// @Summary      Register employee measurements
// @Description  Register body measurements for a specific employee
// @Tags         employees
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path      string                   true  "Employee ID"
// @Param        request  body      models.BodyMeasurements  true  "Body measurements"
// @Success      201      {object}  models.BodyMeasurements
// @Failure      400      {object}  models.ErrorResponse
// @Failure      401      {object}  models.ErrorResponse
// @Failure      500      {object}  models.ErrorResponse
// @Router       /api/v1/funcionarios/{id}/medidas [post]
func (h *EmployeeHandler) RegisterMeasurements(c *gin.Context) {
	id := c.Param("id")
	var req models.BodyMeasurements
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": err.Error()})
		return
	}
	resp, err := h.service.RegisterMeasurements(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, resp)
}

// UpdateMeasurements godoc
// @Summary      Update employee measurements
// @Description  Update body measurements for a specific employee
// @Tags         employees
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path      string                   true  "Employee ID"
// @Param        request  body      models.BodyMeasurements  true  "Body measurements"
// @Success      200      {object}  models.BodyMeasurements
// @Failure      400      {object}  models.ErrorResponse
// @Failure      401      {object}  models.ErrorResponse
// @Failure      500      {object}  models.ErrorResponse
// @Router       /api/v1/funcionarios/{id}/medidas [put]
func (h *EmployeeHandler) UpdateMeasurements(c *gin.Context) {
	id := c.Param("id")
	var req models.BodyMeasurements
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": err.Error()})
		return
	}
	resp, err := h.service.UpdateMeasurements(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// GetMeasurementsHistory godoc
// @Summary      Get employee measurements history
// @Description  Get body measurements history for a specific employee
// @Tags         employees
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "Employee ID"
// @Success      200  {array}   models.BodyMeasurements
// @Failure      401  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /api/v1/funcionarios/{id}/medidas/historial [get]
func (h *EmployeeHandler) GetMeasurementsHistory(c *gin.Context) {
	id := c.Param("id")
	// La relación es 1:1 en el DDL actual (Funcionario.id_medidas → Medidas Funcionario).
	// Se retorna la medida actual en un array. Para historial real se requiere cambio de DDL.
	m, err := h.service.GetMeasurements(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Error", "message": err.Error()})
		return
	}
	if m == nil {
		c.JSON(http.StatusOK, []models.BodyMeasurements{})
		return
	}
	c.JSON(http.StatusOK, []models.BodyMeasurements{*m})
}

// ─── Admin stubs — 501 Not Implemented ─────────────────────────────────────────

func notImplemented(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"error":   "Not Implemented",
		"message": "Este endpoint está pendiente de implementación en el microservicio",
	})
}

func (h *EmployeeHandler) ListEmployees(c *gin.Context)      { notImplemented(c) }
func (h *EmployeeHandler) CreateEmployee(c *gin.Context)     { notImplemented(c) }
func (h *EmployeeHandler) GetEmployeeByID(c *gin.Context)    { notImplemented(c) }
func (h *EmployeeHandler) UpdateEmployee(c *gin.Context)     { notImplemented(c) }
func (h *EmployeeHandler) DeleteEmployee(c *gin.Context)     { notImplemented(c) }
func (h *EmployeeHandler) FilterEmployees(c *gin.Context)    { notImplemented(c) }
func (h *EmployeeHandler) ActivateEmployee(c *gin.Context)   { notImplemented(c) }
func (h *EmployeeHandler) DeactivateEmployee(c *gin.Context) { notImplemented(c) }
