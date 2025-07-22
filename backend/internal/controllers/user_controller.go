package controllers

import (
	"book_boy/backend/internal/bl"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	Service bl.UserService
}

func NewUserController(service bl.UserService) *UserController {
	return &UserController{Service: service}
}

func (uc *UserController) GetAll(c *gin.Context) {
	users, err := uc.Service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": users})
}
