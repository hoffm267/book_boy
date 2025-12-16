package controllers

import (
	"book_boy/api/internal/service"
	"book_boy/api/internal/domain"
	"book_boy/api/internal/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TrackingController struct {
	Service service.TrackingService
}

func NewTrackingController(service service.TrackingService) *TrackingController {
	return &TrackingController{Service: service}
}

func (tc *TrackingController) RegisterRoutes(r gin.IRouter) {
	tracking := r.Group("/tracking")
	tracking.POST("/start", tc.StartTracking)
	tracking.GET("/current", tc.GetCurrentTracking)
}

func (tc *TrackingController) StartTracking(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var req domain.StartTrackingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	progress, err := tc.Service.StartTracking(userID.(int), &req)
	if err != nil {
		if errors.IsValidationError(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, progress)
}

func (tc *TrackingController) GetCurrentTracking(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	currentTracking, err := tc.Service.GetCurrentTracking(userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, currentTracking)
}
