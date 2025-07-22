package controllers

import (
	"book_boy/backend/internal/bl"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProgressController struct {
	Service bl.ProgressService
}

func NewProgressController(service bl.ProgressService) *ProgressController {
	return &ProgressController{Service: service}
}

func (pc *ProgressController) GetAll(c *gin.Context) {
	progress, err := pc.Service.GetAllProgress()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": progress})
}
