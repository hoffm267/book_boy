package controllers

import (
	"book_boy/backend/internal/bl"
	"book_boy/backend/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AudiobookController struct {
	Service bl.AudiobookService
}

func NewAudiobookController(service bl.AudiobookService) *AudiobookController {
	return &AudiobookController{Service: service}
}

func (ac *AudiobookController) RegisterRoutes(r *gin.Engine) {
	audiobooks := r.Group("/audiobooks")
	{
		audiobooks.GET("", ac.GetAll)
		audiobooks.GET("/:id", ac.GetByID)
		audiobooks.POST("", ac.Create)
		audiobooks.PUT("/:id", ac.Update)
		audiobooks.DELETE("/:id", ac.Delete)
	}
}

func (ac *AudiobookController) GetAll(c *gin.Context) {
	result, err := ac.Service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
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
	var audiobook models.Audiobook
	if err := c.ShouldBindJSON(&audiobook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := ac.Service.Create(&audiobook)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
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
	var audiobook models.Audiobook
	if err := c.ShouldBindJSON(&audiobook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	audiobook.ID = id
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
