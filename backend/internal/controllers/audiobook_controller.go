package controllers

import (
	"book_boy/backend/internal/bl"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AudiobookController struct {
	Service bl.AudiobookService
}

func NewAudiobookController(service bl.AudiobookService) *AudiobookController {
	return &AudiobookController{Service: service}
}

func (ac *AudiobookController) GetAll(c *gin.Context) {
	audiobooks, err := ac.Service.GetAllAudiobooks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": audiobooks})
}
