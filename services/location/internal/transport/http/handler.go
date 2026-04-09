package http

import (
	"net/http"

	"github.com/care4u/services/location/internal/domain"
	"github.com/gin-gonic/gin"
)

type LocationHandler struct {
	usecase domain.LocationUsecase
}

func NewLocationHandler(usecase domain.LocationUsecase) *LocationHandler {
	return &LocationHandler{usecase: usecase}
}

func (h *LocationHandler) RegisterRoutes(r *gin.Engine) {
	location := r.Group("/location")
	{
		location.GET("/nearby", h.GetNearbyProviders)
		location.POST("/update", h.UpdateLocation)
	}
}

func (h *LocationHandler) GetNearbyProviders(c *gin.Context) {
	var req domain.NearbyRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	providers, err := h.usecase.GetNearbyProviders(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"providers": providers, "count": len(providers)})
}

func (h *LocationHandler) UpdateLocation(c *gin.Context) {
	var req domain.UpdateLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.usecase.UpdateLocation(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Location updated"})
}
