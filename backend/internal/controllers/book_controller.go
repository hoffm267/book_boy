package controllers

import (
	"book_boy/backend/internal/bl"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BookController struct {
	Service bl.BookService
}

func NewBookController(service bl.BookService) *BookController {
	return &BookController{Service: service}
}

func (bc *BookController) GetAll(c *gin.Context) {
	books, err := bc.Service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": books})
}

func (bc *BookController) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid book ID"})
		return
	}

	book, err := bc.Service.GetByID(id)
	if err != nil {
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}
	if book == nil {
		c.JSON(404, gin.H{"error": "book not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": book})
}
