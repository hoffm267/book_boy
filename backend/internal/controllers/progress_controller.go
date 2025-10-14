package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"book_boy/backend/internal/bl"
	"book_boy/backend/internal/models"
)

type ProgressController struct {
	Service bl.ProgressService
}

type updatePageReq struct {
	Page int `json:"page" binding:"required,min=1"`
}

type updateTimeReq struct {
	AudiobookTime models.CustomDuration `json:"audiobook_time" binding:"required"`
}

func NewProgressController(Service bl.ProgressService) *ProgressController {
	return &ProgressController{Service: Service}
}

func (pc *ProgressController) RegisterRoutes(r gin.IRouter) {
	progress := r.Group("/progress")
	progress.GET("", pc.GetAll)
	progress.GET("/:id", pc.GetByID)
	progress.POST("", pc.Create)
	progress.PUT("/:id", pc.Update)
	progress.DELETE("/:id", pc.Delete)
	progress.PATCH("/:id/page", pc.UpdateByPage) // update by page
	progress.PATCH("/:id/time", pc.UpdateByTime) // update by timestamp
	progress.GET("/filter", pc.FilterProgress)
}

func (pc *ProgressController) GetAll(c *gin.Context) {
	progress, err := pc.Service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": progress})
}

func (pc *ProgressController) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	progress, err := pc.Service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if progress == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "progress not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": progress})
}

func (pc *ProgressController) Create(c *gin.Context) {
	// Get authenticated user ID from token (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	// Bind request body
	var progress models.Progress
	if err := c.ShouldBindJSON(&progress); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Override user_id with authenticated user (ignore any user_id from request body)
	progress.UserID = userID.(int)

	id, err := pc.Service.Create(&progress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	progress.ID = id
	c.JSON(http.StatusCreated, progress)
}

func (pc *ProgressController) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	existing, err := pc.Service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if existing == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "progress not found"})
		return
	}

	if existing.UserID != userID.(int) {
		c.JSON(http.StatusForbidden, gin.H{"error": "you can only update your own progress"})
		return
	}

	var progress models.Progress
	if err := c.ShouldBindJSON(&progress); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO make sure overriding is correct
	progress.ID = id
	progress.UserID = existing.UserID

	if err := pc.Service.Update(&progress); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, progress)
}

func (pc *ProgressController) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	existing, err := pc.Service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if existing == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "progress not found"})
		return
	}

	if existing.UserID != userID.(int) {
		c.JSON(http.StatusForbidden, gin.H{"error": "you can only delete your own progress"})
		return
	}

	if err := pc.Service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (pc *ProgressController) UpdateByPage(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req updatePageReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := pc.Service.UpdateProgressPage(id, req.Page); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (pc *ProgressController) UpdateByTime(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req updateTimeReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := pc.Service.UpdateProgressTime(id, &req.AudiobookTime); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (pc *ProgressController) FilterProgress(c *gin.Context) {
	var filter models.ProgressFilter
	if idStr := c.Query("id"); idStr != "" {
		if id, err := strconv.Atoi(idStr); err == nil {
			filter.ID = &id
		}
	}
	if userIdStr := c.Query("user_id"); userIdStr != "" {
		if userId, err := strconv.Atoi(userIdStr); err == nil {
			filter.UserID = &userId
		}
	}
	if bookIdStr := c.Query("book_id"); bookIdStr != "" {
		if bookId, err := strconv.Atoi(bookIdStr); err == nil {
			filter.BookID = &bookId
		}
	}
	if audiobookIdStr := c.Query("audiobook_id"); audiobookIdStr != "" {
		if audiobookId, err := strconv.Atoi(audiobookIdStr); err == nil {
			filter.AudiobookID = &audiobookId
		}
	}

	progresses, err := pc.Service.FilterProgress(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch progresss"})
		return
	}

	c.JSON(http.StatusOK, progresses)
}
