package handlers

import (
	"net/http"

	"github.com/Sistal/bff-orchestrator/internal/services"
	"github.com/gin-gonic/gin"
)

type CatalogHandler struct {
	service services.CatalogService
}

func NewCatalogHandler(s services.CatalogService) *CatalogHandler {
	return &CatalogHandler{service: s}
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
