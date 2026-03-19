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

type CatalogHandler struct {
	service         services.CatalogService
	employeeService services.EmployeeService
}

func NewCatalogHandler(s services.CatalogService, es services.EmployeeService) *CatalogHandler {
	return &CatalogHandler{service: s, employeeService: es}
}

// GetSizes godoc
// @Summary      Get sizes catalog
// @Description  Get all available sizes
// @Tags         catalog
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]string
// @Router       /catalogo/tallas [get]
func (h *CatalogHandler) GetSizes(c *gin.Context) {
	resp, err := h.service.GetSizes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// GetChangeReasons godoc
// @Summary      Get change reasons catalog
// @Description  Get all change reasons for garment changes
// @Tags         catalog
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]string
// @Router       /catalogo/motivos-cambio [get]
func (h *CatalogHandler) GetChangeReasons(c *gin.Context) {
	resp, err := h.service.GetChangeReasons()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// GetGarmentTypes godoc
// @Summary      Get garment types catalog
// @Description  Get all available garment types
// @Tags         catalog
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]string
// @Router       /catalogo/prenda-tipos [get]
func (h *CatalogHandler) GetGarmentTypes(c *gin.Context) {
	resp, err := h.service.GetGarmentTypes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// GetActiveCampaign godoc
// @Summary      Get active campaign
// @Description  Get the currently active campaign
// @Tags         catalog
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]string
// @Router       /campanas/activa [get]
func (h *CatalogHandler) GetActiveCampaign(c *gin.Context) {
	resp, err := h.service.GetActiveCampaign()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// GetCompanies godoc
// @Summary      Get companies
// @Description  Get all available companies
// @Tags         catalog
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  []models.Empresa
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/empresas [get]
func (h *CatalogHandler) GetCompanies(c *gin.Context) {
	log := logger.Get()
	log.Info("GetCompanies: Request recibida", zap.String("ip", c.ClientIP()))

	resp, err := h.service.GetCompanies()
	if err != nil {
		log.Error("GetCompanies: Error obteniendo empresas", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Error", "message": err.Error()})
		return
	}

	log.Info("GetCompanies: Request procesada con éxito", zap.Int("count", len(resp)))
	c.JSON(http.StatusOK, resp)
}

// GetSegments godoc
// @Summary      Get segments
// @Description  Get segments for a company
// @Tags         catalog
// @Security     BearerAuth
// @Produce      json
// @Param        idEmpresa path int true "Company ID"
// @Success      200  {object}  []models.Segmento
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/segmentos/{idEmpresa} [get]
func (h *CatalogHandler) GetSegments(c *gin.Context) {
	log := logger.Get()
	idStr := c.Param("idEmpresa")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Warn("GetSegments: ID de empresa inválido", zap.String("id_param", idStr), zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Company ID"})
		return
	}

	log.Info("GetSegments: Request recibida", zap.Int("id_empresa", id), zap.String("ip", c.ClientIP()))

	resp, err := h.service.GetSegments(id)
	if err != nil {
		log.Error("GetSegments: Error obteniendo segmentos", zap.Int("id_empresa", id), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Error", "message": err.Error()})
		return
	}

	log.Info("GetSegments: Request procesada con éxito", zap.Int("id_empresa", id), zap.Int("count", len(resp)))
	c.JSON(http.StatusOK, resp)
}

// GetBranches godoc
// @Summary      Get branches
// @Description  Get branches for a company
// @Tags         catalog
// @Security     BearerAuth
// @Produce      json
// @Param        idEmpresa path int true "Company ID"
// @Success      200  {object}  []models.Sucursal
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/sucursales/{idEmpresa} [get]
func (h *CatalogHandler) GetBranches(c *gin.Context) {
	log := logger.Get()
	idStr := c.Param("idEmpresa")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Warn("GetBranches: ID de empresa inválido", zap.String("id_param", idStr), zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Company ID"})
		return
	}

	log.Info("GetBranches: Request recibida", zap.Int("id_empresa", id), zap.String("ip", c.ClientIP()))

	resp, err := h.service.GetBranches(id)
	if err != nil {
		log.Error("GetBranches: Error obteniendo sucursales", zap.Int("id_empresa", id), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Error", "message": err.Error()})
		return
	}

	log.Info("GetBranches: Request procesada con éxito", zap.Int("id_empresa", id), zap.Int("count", len(resp)))
	req := models.CatalogResponse[models.Sucursal]{
		Success: true,
		Data:    resp,
	}
	c.JSON(http.StatusOK, req)
}

// GetUniforms godoc
// @Summary      Get uniforms by segment
// @Description  Get uniforms available for the authenticated employee's segment
// @Tags         catalog
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  models.CatalogResponse[models.Uniform]
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /catalogo/uniformes [get]
func (h *CatalogHandler) GetUniforms(c *gin.Context) {
	employeeIDStr, ok := requireEmployeeID(c)
	if !ok {
		return
	}

	// Obtener perfil para conocer el segmento
	profile, err := h.employeeService.GetProfile(employeeIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Error", "message": "Failed to get employee profile"})
		return
	}

	// Consultar ms-catalogo con el ID de segmento
	data, err := h.service.GetUniformsBySegment(profile.IDSegmento)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Error", "message": err.Error()})
		return
	}

	req := models.CatalogResponse[models.Uniform]{
		Success: true,
		Data:    data,
	}
	c.JSON(http.StatusOK, req)
}
