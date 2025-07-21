package controllers

import (
	"book_boy/backend/internal/bl"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BookController struct {
	Service bl.BookService
}

func NewBookController(service bl.BookService) *BookController {
	return &BookController{Service: service}
}

func (bc *BookController) GetAll(c *gin.Context) {
	books, err := bc.Service.GetAllBooks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": books})
}
