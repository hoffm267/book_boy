package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"book_boy/api/internal/service"
	"book_boy/api/internal/domain"
	"book_boy/api/internal/repository"
)

type ProgressController struct {
	Service          service.ProgressService
	BookService      service.BookService
	AudiobookService service.AudiobookService
}

type updatePageReq struct {
	Page int `json:"page" binding:"required,min=1"`
}

type updateTimeReq struct {
	AudiobookTime domain.CustomDuration `json:"audiobook_time" binding:"required"`
}

func NewProgressController(Service service.ProgressService, BookService service.BookService, AudiobookService service.AudiobookService) *ProgressController {
	return &ProgressController{
		Service:          Service,
		BookService:      BookService,
		AudiobookService: AudiobookService,
	}
}

func (pc *ProgressController) RegisterRoutes(r gin.IRouter) {
	progress := r.Group("/progress")
	progress.GET("", pc.GetAll)
	progress.GET("/:id", pc.GetByID)
	progress.POST("", pc.Create)
	progress.PUT("/:id", pc.Update)
	progress.DELETE("/:id", pc.Delete)
	progress.PATCH("/:id/page", pc.UpdateByPage)
	progress.PATCH("/:id/time", pc.UpdateByTime)
	progress.GET("/filter", pc.FilterProgress)
	progress.GET("/enriched", pc.GetEnrichedByUser)
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
	progress, err := pc.Service.GetByIDWithCompletion(id)
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
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var progress domain.Progress
	if err := c.ShouldBindJSON(&progress); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	progress.UserID = userID.(int)

	uid := userID.(int)
	var filter repository.ProgressFilter
	filter.UserID = &uid

	if progress.BookID != nil {
		filter.BookID = progress.BookID
		existingProgress, err := pc.Service.FilterProgress(filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if len(existingProgress) > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "progress already exists for this book"})
			return
		}
	}

	if progress.AudiobookID != nil {
		filter.AudiobookID = progress.AudiobookID
		filter.BookID = nil
		existingProgress, err := pc.Service.FilterProgress(filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if len(existingProgress) > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "progress already exists for this audiobook"})
			return
		}
	}

	id, err := pc.Service.Create(&progress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	progress.ID = id
	c.JSON(http.StatusCreated, progress)
}

type updateProgressReq struct {
	BookPage      *int                   `json:"book_page"`
	AudiobookTime *domain.CustomDuration `json:"audiobook_time"`
	BookID        *int                   `json:"book_id"`
	AudiobookID   *int                   `json:"audiobook_id"`
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

	var req updateProgressReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.BookID != nil {
		if err := pc.Service.SetBook(id, *req.BookID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if req.AudiobookID != nil {
		if err := pc.Service.SetAudiobook(id, *req.AudiobookID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if req.BookPage != nil {
		if err := pc.Service.UpdateProgressPage(id, *req.BookPage); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if req.AudiobookTime != nil {
		if err := pc.Service.UpdateProgressTime(id, req.AudiobookTime); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	updated, err := pc.Service.GetByIDWithCompletion(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updated)
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

	bookID := existing.BookID
	audiobookID := existing.AudiobookID

	if err := pc.Service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if bookID != nil {
		var filter repository.ProgressFilter
		filter.BookID = bookID
		remainingProgress, err := pc.Service.FilterProgress(filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if len(remainingProgress) == 0 {
			if err := pc.BookService.Delete(*bookID); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
	}

	if audiobookID != nil {
		var filter repository.ProgressFilter
		filter.AudiobookID = audiobookID
		remainingProgress, err := pc.Service.FilterProgress(filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if len(remainingProgress) == 0 {
			if err := pc.AudiobookService.Delete(*audiobookID); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
	}

	c.Status(http.StatusNoContent)
}

func (pc *ProgressController) UpdateByPage(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req updatePageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := pc.Service.UpdateProgressPage(id, req.Page); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (pc *ProgressController) UpdateByTime(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req updateTimeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := pc.Service.UpdateProgressTime(id, &req.AudiobookTime); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (pc *ProgressController) FilterProgress(c *gin.Context) {
	var filter repository.ProgressFilter
	if idStr := c.Query("id"); idStr != "" {
		if id, err := strconv.Atoi(idStr); err == nil {
			filter.ID = &id
		}
	}
	if userIDStr := c.Query("user_id"); userIDStr != "" {
		if userID, err := strconv.Atoi(userIDStr); err == nil {
			filter.UserID = &userID
		}
	}
	if bookIDStr := c.Query("book_id"); bookIDStr != "" {
		if bookID, err := strconv.Atoi(bookIDStr); err == nil {
			filter.BookID = &bookID
		}
	}
	if audiobookIDStr := c.Query("audiobook_id"); audiobookIDStr != "" {
		if audiobookID, err := strconv.Atoi(audiobookIDStr); err == nil {
			filter.AudiobookID = &audiobookID
		}
	}
	if status := c.Query("status"); status != "" {
		progressStatus := domain.ProgressStatus(status)
		filter.Status = &progressStatus
	}

	progresses, err := pc.Service.FilterProgress(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch progresss"})
		return
	}

	c.JSON(http.StatusOK, progresses)
}

func (pc *ProgressController) GetEnrichedByUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	enriched, err := pc.Service.GetAllEnrichedByUser(userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch enriched progress"})
		return
	}

	c.JSON(http.StatusOK, enriched)
}
