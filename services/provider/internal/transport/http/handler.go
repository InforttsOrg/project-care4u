package http

import (
	"net/http"
	"strconv"

	"github.com/care4u/services/provider/internal/domain"
	"github.com/gin-gonic/gin"
)

type ProviderHandler struct {
	usecase domain.ProviderUsecase
}

func NewProviderHandler(usecase domain.ProviderUsecase) *ProviderHandler {
	return &ProviderHandler{usecase: usecase}
}

func (h *ProviderHandler) RegisterRoutes(r *gin.Engine) {
	providers := r.Group("/providers")
	{
		providers.POST("/", h.CreateProvider)
		providers.GET("/search", h.SearchProviders)
		providers.GET("/:id", h.GetProvider)
		providers.PUT("/:id", h.UpdateProvider)
		providers.PUT("/:id/verify", h.VerifyProvider) // Admin only in production
	}

	// Get provider by user ID
	r.GET("/users/:user_id/provider", h.GetProviderByUserID)
}

func (h *ProviderHandler) CreateProvider(c *gin.Context) {
	var req domain.CreateProviderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	provider, err := h.usecase.CreateProvider(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, provider)
}

func (h *ProviderHandler) GetProvider(c *gin.Context) {
	id := c.Param("id")
	provider, err := h.usecase.GetProvider(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, provider)
}

func (h *ProviderHandler) GetProviderByUserID(c *gin.Context) {
	userID := c.Param("user_id")
	provider, err := h.usecase.GetProviderByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, provider)
}

func (h *ProviderHandler) UpdateProvider(c *gin.Context) {
	id := c.Param("id")
	var req domain.UpdateProviderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.usecase.UpdateProvider(c.Request.Context(), id, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "provider updated"})
}

func (h *ProviderHandler) SearchProviders(c *gin.Context) {
	specialty := c.Query("specialty")
	limit := 20
	offset := 0

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}
	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	providers, err := h.usecase.SearchProviders(c.Request.Context(), specialty, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"providers": providers, "count": len(providers)})
}

func (h *ProviderHandler) VerifyProvider(c *gin.Context) {
	id := c.Param("id")
	// In production, check admin role from JWT

	if err := h.usecase.VerifyProvider(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "provider verified"})
}
