package handlers

import (
	"net/http"

	"github.com/Sistal/bff-orchestrator/internal/models"
	"github.com/Sistal/bff-orchestrator/internal/services"
	"github.com/gin-gonic/gin"
)

type BranchHandler struct {
	service services.BranchService
}

func NewBranchHandler(s services.BranchService) *BranchHandler {
	return &BranchHandler{service: s}
}

// requireEmployeeID movido a common.go

// GetAllBranches godoc
// @Summary      Get all branches
// @Description  Get list of all branches
// @Tags         branches
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]string
// @Router       /sucursales [get]
func (h *BranchHandler) GetAllBranches(c *gin.Context) {
	resp, err := h.service.GetAllBranches()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// GetChangeHistory godoc
// @Summary      Get branch change history
// @Description  Get history of branch change requests for the authenticated employee
// @Tags         branches
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /solicitudes/cambio-sucursal/historial [get]
func (h *BranchHandler) GetChangeHistory(c *gin.Context) {
	employeeID, ok := requireEmployeeID(c)
	if !ok {
		return
	}
	resp, err := h.service.GetChangeHistory(employeeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// CreateChangeRequest godoc
// @Summary      Create branch change request
// @Description  Create a new branch change request for the authenticated employee
// @Tags         branches
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body      models.CreateBranchChangeRequest  true  "Branch change request"
// @Success      201      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]string
// @Failure      401      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /solicitudes/cambio-sucursal [post]
func (h *BranchHandler) CreateChangeRequest(c *gin.Context) {
	employeeID, ok := requireEmployeeID(c)
	if !ok {
		return
	}
	var req models.CreateBranchChangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": err.Error()})
		return
	}

	resp, err := h.service.CreateChangeRequest(employeeID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, resp)
}
