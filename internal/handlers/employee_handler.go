package handlers

import (
	"net/http"
	"strconv"

	"github.com/Sistal/bff-orchestrator/internal/logger"
	"github.com/Sistal/bff-orchestrator/internal/models"
	"github.com/Sistal/bff-orchestrator/internal/services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
	log := logger.Get()
	employeeID, ok := h.requireEmployeeByUserId(c)
	if !ok {
		log.Warn("GetProfile: Usuario no autorizado o sin employeeID")
		return
	}

	log.Debug("GetProfile: Solicitud recibida", zap.String("employee_id", employeeID))

	resp, err := h.service.GetProfile(employeeID)
	if err != nil {
		log.Error("GetProfile: Error al obtener perfil", zap.String("employee_id", employeeID), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Error", "message": err.Error()})
		return
	}

	log.Debug("GetProfile: Perfil retornado exitosamente", zap.String("employee_id", employeeID))
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
	employeeID, ok := h.requireEmployeeByUserId(c)
	if !ok {
		return
	}
	var req models.UpdateContactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": err.Error()})
		return
	}
	resp, err := h.service.UpdateContact(employeeID, req)
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
	employeeID, ok := h.requireEmployeeByUserId(c)
	if !ok {
		return
	}
	var req models.UpdatePreferencesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": err.Error()})
		return
	}
	resp, err := h.service.UpdatePreferences(employeeID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Error", "message": err.Error()})
		return
	}
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
	employeeID, ok := h.requireEmployeeByUserId(c)
	if !ok {
		return
	}
	var req models.UpdateSecurityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": err.Error()})
		return
	}
	if err := h.service.UpdateSecurity(employeeID, req); err != nil {
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
	employeeID, ok := h.requireEmployeeByUserId(c)
	if !ok {
		return
	}
	resp, err := h.service.GetStats(employeeID)
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
	employeeID, ok := h.requireEmployeeByUserId(c)
	if !ok {
		return
	}
	resp, err := h.service.GetActivity(employeeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// GetMeasurements godoc
// @Summary      Get employee measurements
// @Description  Get body measurements for a specific employee (admin)
// @Tags         employees
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "Employee ID"
// @Success      200  {object}  models.BodyMeasurements
// @Failure      401  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /api/v1/funcionarios/{id}/medidas [get]
func (h *EmployeeHandler) GetMeasurements(c *gin.Context) {
	// Se extrae el id_funcionario desde el userID del token, ignorando el param :id
	id, ok := h.requireEmployeeByUserId(c)
	if !ok {
		return
	}

	resp, err := h.service.GetMeasurements(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// RegisterMeasurements godoc
// @Summary      Register employee measurements
// @Description  Register body measurements for a specific employee (admin)
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
// @Description  Update body measurements for a specific employee (admin)
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
// @Description  Get body measurements history for a specific employee (admin)
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
func (h *EmployeeHandler) GetEmployeeByID(c *gin.Context)    { notImplemented(c) }
func (h *EmployeeHandler) UpdateEmployee(c *gin.Context)     { notImplemented(c) }
func (h *EmployeeHandler) DeleteEmployee(c *gin.Context)     { notImplemented(c) }
func (h *EmployeeHandler) FilterEmployees(c *gin.Context)    { notImplemented(c) }
func (h *EmployeeHandler) ActivateEmployee(c *gin.Context)   { notImplemented(c) }
func (h *EmployeeHandler) DeactivateEmployee(c *gin.Context) { notImplemented(c) }

// CreateEmployee godoc
// @Summary      Create employee
// @Description  Create a new employee profile linked to a user.
// @Tags         employees
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body      models.CreateEmployeeRequest  true  "Employee details"
// @Success      201      {object}  models.EmployeeProfile
// @Failure      400      {object}  models.ErrorResponse
// @Failure      401      {object}  models.ErrorResponse
// @Failure      500      {object}  models.ErrorResponse
// @Router       /api/v1/funcionarios [post]
func (h *EmployeeHandler) CreateEmployee(c *gin.Context) {
	// Obtener userID desde el token (contexto)
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": "User ID not found in token"})
		return
	}

	var req models.CreateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": err.Error()})
		return
	}

	// Forzar el UserID en el request con el del token para seguridad
	// (En caso de que sea un admin creando para otro, la lógica podría cambiar,
	//  pero aquí asumimos que el usuario que lo crea es el mismo si no es admin,
	//  o si es admin, habría que ver si se permite sobreescribir).
	//
	// Vamos a asumir aquí que el endpoint es para que un usuario se "auto-complete" como funcionario
	// O para que un admin cree uno. Para simplificar, usamos el userID del token.
	// Si se desea soportar creación por admins para terceros, habría que chequear roles.
	// Por ahora: Parseamos userID a int y lo asignamos.
	importStrconv, err := strconv.Atoi(userID)
	if err == nil {
		// Si en el json venía otro user_id y no somos admin, esto lo sobreescribe.
		// Si queremos permitir que un admin ponga cualquier ID, deberíamos chequear rol.
		// Dado el requerimiento "extrae el id_usuario del token", priorizamos el token.
		req.UserID = importStrconv
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "Invalid User ID in token"})
		return
	}

	resp, err := h.service.CreateEmployee(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, resp)
}

// InitialRegister godoc
// @Summary      Initial register update for employee
// @Description  Update employee data for initial registration using data from body. Resolves employee ID from token.
// @Tags         employees
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body      object  true  "Employee update data"
// @Success      200      {object}  map[string]string
// @Failure      400      {object}  models.ErrorResponse
// @Failure      401      {object}  models.ErrorResponse
// @Failure      500      {object}  models.ErrorResponse
// @Router       /api/v1/funcionarios/registro-inicial [post]
func (h *EmployeeHandler) InitialRegister(c *gin.Context) {
	employeeID, ok := h.requireEmployeeByUserId(c)
	if !ok {
		return
	}

	var req interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": err.Error()})
		return
	}

	err := h.service.InitialRegister(employeeID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Registro inicial completado exitosamente"})
}

// requireEmployeeByUserId extrae el userID del contexto y consulta al servicio para obtener el employeeID.
// Reemplaza la lógica anterior que dependía del middleware para poblar "employeeID".
func (h *EmployeeHandler) requireEmployeeByUserId(c *gin.Context) (string, bool) {
	log := logger.Get()
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "User ID no encontrado en el contexto",
		})
		return "", false
	}

	empIDInt, err := h.service.GetEmployeeByUserID(userID)
	if err != nil {
		log.Error("requireEmployeeByUserId: Error al obtener funcionario por usuario",
			zap.String("user_id", userID),
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Error", "message": "Error al consultar información del funcionario"})
		return "", false
	}

	if empIDInt == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "Usuario sin funcionario asociado",
		})
		return "", false
	}

	return strconv.Itoa(empIDInt), true
}
