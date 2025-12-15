package controllers

import (
	"book_boy/internal/bl"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BookMetadataController struct {
	Service *bl.BookMetadataService
}

func NewBookMetadataController(service *bl.BookMetadataService) *BookMetadataController {
	return &BookMetadataController{Service: service}
}

func (bmc *BookMetadataController) RegisterRoutes(r gin.IRouter) {
	metadata := r.Group("/metadata")
	{
		metadata.GET("/search", bmc.SearchBooks)
		metadata.GET("/isbn/:isbn", bmc.GetBookByISBN)
	}
}

func (bmc *BookMetadataController) SearchBooks(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'q' is required"})
		return
	}

	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil {
			limit = parsedLimit
		}
	}

	results, err := bmc.Service.SearchBooks(query, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
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
