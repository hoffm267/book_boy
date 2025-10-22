package controllers

import (
	"book_boy/internal/bl"
	"book_boy/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TrackingController struct {
	Service bl.TrackingService
}

func NewTrackingController(service bl.TrackingService) *TrackingController {
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

	var req models.StartTrackingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	progress, err := tc.Service.StartTracking(userID.(int), &req)
	if err != nil {
		if models.IsValidationError(err) {
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
