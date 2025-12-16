package controllers

import (
	"book_boy/api/internal/service"
	"book_boy/api/internal/domain"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AudiobookController struct {
	Service         service.AudiobookService
	ProgressService service.ProgressService
}

func NewAudiobookController(service service.AudiobookService, pgService service.ProgressService) *AudiobookController {
	return &AudiobookController{Service: service, ProgressService: pgService}
}

func (ac *AudiobookController) RegisterRoutes(r gin.IRouter) {
	audiobooks := r.Group("/audiobooks")
	{
		audiobooks.GET("", ac.GetAll)
		audiobooks.GET("/:id", ac.GetByID)
		audiobooks.POST("", ac.Create)
		audiobooks.PUT("/:id", ac.Update)
		audiobooks.DELETE("/:id", ac.Delete)
		audiobooks.GET("/search", ac.GetSimilarTitles)
	}
}

func (ac *AudiobookController) GetAll(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	result, err := ac.Service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	uid := userID.(int)
	filtered := []domain.Audiobook{}
	for _, a := range result {
		if a.UserID == uid {
			filtered = append(filtered, a)
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": filtered})
}

func (ac *AudiobookController) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}
	result, err := ac.Service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if result == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}

func (ac *AudiobookController) Create(c *gin.Context) {
	var audiobook domain.Audiobook
	var progress domain.Progress

	if err := c.ShouldBindJSON(&audiobook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	audiobook.UserID = userID.(int)

	id, err := ac.Service.Create(&audiobook)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	pgIDStr := c.Query("pgId")
	if pgIDStr != "" {
		pgID, err := strconv.Atoi(pgIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid audiobook id"})
			return
		}

		if err := ac.ProgressService.SetAudiobook(pgID, id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		progress = domain.Progress{
			UserID:      userID.(int),
			AudiobookID: &id,
		}

		_, err = ac.ProgressService.Create(&progress)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	audiobook.ID = id
	c.JSON(http.StatusCreated, gin.H{"data": audiobook})
}

func (ac *AudiobookController) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var audiobook domain.Audiobook
	if err := c.ShouldBindJSON(&audiobook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	audiobook.ID = id
	audiobook.UserID = userID.(int)

	if err := ac.Service.Update(&audiobook); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": audiobook})
}

func (ac *AudiobookController) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}
	if err := ac.Service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (ac *AudiobookController) GetSimilarTitles(c *gin.Context) {
	title := c.Query("title")
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing title query parameter"})
		return
	}

	audiobooks, err := ac.Service.GetSimilarTitles(title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if audiobooks == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "books not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": audiobooks})
}
