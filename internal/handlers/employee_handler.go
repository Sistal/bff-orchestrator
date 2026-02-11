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
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/funcionarios/me [get]
func (h *EmployeeHandler) GetProfile(c *gin.Context) {
	// Use userID from context instead of assuming 123
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
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]string
// @Failure      500      {object}  map[string]string
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
// @Description  Update the current employee's preferences
// @Tags         employees
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body      models.UpdatePreferencesRequest  true  "Preferences"
// @Success      200
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/funcionarios/me/preferencias [put]
func (h *EmployeeHandler) UpdatePreferences(c *gin.Context) {
	userID := c.GetString("userID")
	var req models.UpdatePreferencesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": err.Error()})
		return
	}

	if err := h.service.UpdatePreferences(userID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Error", "message": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

// UpdateSecurity godoc
// @Summary      Update employee security
// @Description  Update the current employee's security settings
// @Tags         employees
// @Security     BearerAuth
// @Produce      json
// @Success      200
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/funcionarios/me/seguridad [put]
func (h *EmployeeHandler) UpdateSecurity(c *gin.Context) {
	userID := c.GetString("userID")
	if err := h.service.UpdateSecurity(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Error", "message": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

// GetStats godoc
// @Summary      Get employee statistics
// @Description  Get statistics for the current employee
// @Tags         employees
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]string
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
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]string
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
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/funcionarios/{id}/medidas [get]
func (h *EmployeeHandler) GetMeasurements(c *gin.Context) {
	// id := c.Param("id") // Removed param usage if we use context userID
	userID := c.GetString("userID")
	resp, err := h.service.GetMeasurements(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// RegisterMeasurements godoc
// @Summary      Register employee measurements
// @Description  Register or update body measurements for a specific employee
// @Tags         employees
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path      string                  true  "Employee ID"
// @Param        request  body      models.BodyMeasurements  true  "Body measurements"
// @Success      201      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /api/v1/funcionarios/{id}/medidas [post]
func (h *EmployeeHandler) RegisterMeasurements(c *gin.Context) {
	// id := c.Param("id")
	userID := c.GetString("userID")
	var req models.BodyMeasurements
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": err.Error()})
		return
	}

	resp, err := h.service.RegisterMeasurements(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, resp)
}
