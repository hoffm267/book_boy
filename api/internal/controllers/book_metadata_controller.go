package controllers

import (
	"book_boy/api/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BookMetadataController struct {
	Service *service.BookMetadataService
}

func NewBookMetadataController(service *service.BookMetadataService) *BookMetadataController {
	return &BookMetadataController{Service: service}
}

func (bmc *BookMetadataController) RegisterRoutes(r gin.IRouter) {
	metadata := r.Group("/metadata")
	{
		metadata.GET("/isbn/:isbn", bmc.GetBookByISBN)
	}
}

func (bmc *BookMetadataController) GetBookByISBN(c *gin.Context) {
	isbn := c.Param("isbn")
	if isbn == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "isbn is required"})
		return
	}

	book, err := bmc.Service.GetBookByISBN(isbn)
	if err != nil {
		if err.Error() == "book not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, book)
}
