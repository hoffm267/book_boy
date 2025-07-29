package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"book_boy/backend/internal/bl"
	"book_boy/backend/internal/models"
)

type ProgressController struct {
	service bl.ProgressService
}

func NewProgressController(service bl.ProgressService) *ProgressController {
	return &ProgressController{service: service}
}

func (pc *ProgressController) RegisterRoutes(r *gin.Engine) {
	progress := r.Group("/progress")
	progress.GET("", pc.GetAll)
	progress.GET("/:id", pc.GetByID)
	progress.POST("", pc.Create)
	progress.PUT("/:id", pc.Update)
	progress.DELETE("/:id", pc.Delete)
}

func (pc *ProgressController) GetAll(c *gin.Context) {
	progress, err := pc.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, progress)
}

func (pc *ProgressController) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	progress, err := pc.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if progress == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "progress not found"})
		return
	}
	c.JSON(http.StatusOK, progress)
}

func (pc *ProgressController) Create(c *gin.Context) {
	var p models.Progress
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := pc.service.Create(&p)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	p.ID = id
	c.JSON(http.StatusCreated, p)
}

func (pc *ProgressController) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var p models.Progress
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p.ID = id
	if err := pc.service.Update(&p); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, p)
}

func (pc *ProgressController) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := pc.service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
