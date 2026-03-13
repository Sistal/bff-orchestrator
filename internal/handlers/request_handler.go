package handlers

import (
	"net/http"

	"github.com/Sistal/bff-orchestrator/internal/models"
	"github.com/Sistal/bff-orchestrator/internal/services"
	"github.com/gin-gonic/gin"
)

type RequestHandler struct {
	service services.RequestService
}

func NewRequestHandler(s services.RequestService) *RequestHandler {
	return &RequestHandler{service: s}
}

// GetRequests godoc
// @Summary      Get requests
// @Description  Get list of requests for the authenticated employee
// @Tags         requests
// @Security     BearerAuth
// @Produce      json
// @Param        tipo     query     string  false  "Request type"
// @Param        periodo  query     string  false  "Period filter"
// @Param        estado   query     string  false  "Status filter"
// @Success      200      {object}  map[string]interface{}
// @Failure      401      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /solicitudes [get]
func (h *RequestHandler) GetRequests(c *gin.Context) {
	employeeID, ok := requireEmployeeID(c)
	if !ok {
		return
	}
	params := map[string]string{
		"tipo":    c.Query("tipo"),
		"periodo": c.Query("periodo"),
		"estado":  c.Query("estado"),
	}
	resp, err := h.service.GetRequests(employeeID, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// GetRequestByID godoc
// @Summary      Get request by ID
// @Description  Get a specific request by its ID for the authenticated employee
// @Tags         requests
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "Request ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /solicitudes/{id} [get]
func (h *RequestHandler) GetRequestByID(c *gin.Context) {
	employeeID, ok := requireEmployeeID(c)
	if !ok {
		return
	}
	id := c.Param("id")
	resp, err := h.service.GetRequestByID(employeeID, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found", "message": "Request not found"})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// CreateReplenishmentRequest godoc
// @Summary      Create replenishment request
// @Description  Create a new replenishment request for the authenticated employee
// @Tags         requests
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body      models.CreateReplenishmentRequest  true  "Replenishment request"
// @Success      201      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]string
// @Failure      401      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /solicitudes/reposicion [post]
func (h *RequestHandler) CreateReplenishmentRequest(c *gin.Context) {
	employeeID, ok := requireEmployeeID(c)
	if !ok {
		return
	}
	var req models.CreateReplenishmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": err.Error()})
		return
	}
	resp, err := h.service.CreateReplenishmentRequest(employeeID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, resp)
}

// CreateGarmentChangeRequest godoc
// @Summary      Create garment change request
// @Description  Create a new garment change request for the authenticated employee
// @Tags         requests
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body      models.CreateGarmentChangeRequest  true  "Garment change request"
// @Success      201      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]string
// @Failure      401      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /solicitudes/cambio-prenda [post]
func (h *RequestHandler) CreateGarmentChangeRequest(c *gin.Context) {
	employeeID, ok := requireEmployeeID(c)
	if !ok {
		return
	}
	var req models.CreateGarmentChangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": err.Error()})
		return
	}
	resp, err := h.service.CreateGarmentChangeRequest(employeeID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, resp)
}

// UploadFile godoc
// @Summary      Upload file
// @Description  Upload a file attachment for the authenticated employee
// @Tags         requests
// @Security     BearerAuth
// @Accept       multipart/form-data
// @Produce      json
// @Param        file  formData  file  true  "File to upload"
// @Success      201   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /archivos/upload [post]
func (h *RequestHandler) UploadFile(c *gin.Context) {
	employeeID, ok := requireEmployeeID(c)
	if !ok {
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "File is required"})
		return
	}
	resp, err := h.service.UploadFile(employeeID, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, resp)
}

// GetRecentRequests godoc
// @Summary      Get recent requests
// @Description  Get the most recent requests for the authenticated employee (dashboard)
// @Tags         requests
// @Security     BearerAuth
// @Produce      json
// @Success      200  {array}   models.RequestSummary
// @Failure      401  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /solicitudes/recent [get]
func (h *RequestHandler) GetRecentRequests(c *gin.Context) {
	employeeID, ok := requireEmployeeID(c)
	if !ok {
		return
	}
	resp, err := h.service.GetRecentRequests(employeeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}
