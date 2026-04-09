package http

import (
	"net/http"
	"strconv"

	"github.com/care4u/services/booking/internal/domain"
	"github.com/gin-gonic/gin"
)

type BookingHandler struct {
	usecase domain.BookingUsecase
}

func NewBookingHandler(usecase domain.BookingUsecase) *BookingHandler {
	return &BookingHandler{usecase: usecase}
}

func (h *BookingHandler) RegisterRoutes(r *gin.Engine) {
	bookings := r.Group("/bookings")
	{
		bookings.POST("/", h.CreateBooking)
		bookings.GET("/:id", h.GetBooking)
		bookings.PUT("/:id/cancel", h.CancelBooking)
		bookings.PUT("/:id/confirm", h.ConfirmBooking)
	}

	// Patient routes
	r.GET("/patients/:patient_id/bookings", h.GetPatientBookings)

	// Provider routes
	r.GET("/providers/:provider_id/bookings", h.GetProviderBookings)
}

func (h *BookingHandler) CreateBooking(c *gin.Context) {
	var req domain.CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	booking, err := h.usecase.CreateBooking(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, booking)
}

func (h *BookingHandler) GetBooking(c *gin.Context) {
	id := c.Param("id")
	booking, err := h.usecase.GetBooking(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, booking)
}

func (h *BookingHandler) GetPatientBookings(c *gin.Context) {
	patientID := c.Param("patient_id")
	limit, offset := h.getPagination(c)

	bookings, err := h.usecase.GetPatientBookings(c.Request.Context(), patientID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"bookings": bookings, "count": len(bookings)})
}

func (h *BookingHandler) GetProviderBookings(c *gin.Context) {
	providerID := c.Param("provider_id")
	limit, offset := h.getPagination(c)

	bookings, err := h.usecase.GetProviderBookings(c.Request.Context(), providerID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"bookings": bookings, "count": len(bookings)})
}

func (h *BookingHandler) CancelBooking(c *gin.Context) {
	id := c.Param("id")
	// In a real app, get userID from JWT token
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID required"})
		return
	}

	if err := h.usecase.CancelBooking(c.Request.Context(), id, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "booking cancelled"})
}

func (h *BookingHandler) ConfirmBooking(c *gin.Context) {
	id := c.Param("id")

	if err := h.usecase.ConfirmBooking(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "booking confirmed"})
}

func (h *BookingHandler) getPagination(c *gin.Context) (int, int) {
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
	return limit, offset
}
