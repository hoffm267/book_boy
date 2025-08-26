package controllers

import (
	"book_boy/backend/internal/bl"
	"book_boy/backend/internal/models"
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

func (bc *BookController) RegisterRoutes(r *gin.Engine) {
	books := r.Group("/books")
	{
		books.GET("", bc.GetAll)
		books.GET("/:id", bc.GetByID)
		books.POST("", bc.Create)
		books.PUT("/:id", bc.Update)
		books.DELETE("/:id", bc.Delete)
		books.GET("/search", bc.GetSimilarTitles)
		books.GET("/filter", bc.FilterBooks)
	}
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
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid book ID"})
		return
	}

	book, err := bc.Service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if book == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": book})
}

func (bc *BookController) Create(c *gin.Context) {
	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := bc.Service.Create(&book)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	book.ID = id
	c.JSON(http.StatusCreated, gin.H{"data": book})
}

func (bc *BookController) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid book ID"})
		return
	}

	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	book.ID = id

	if err := bc.Service.Update(&book); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": book})
}

func (bc *BookController) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid book ID"})
		return
	}

	if err := bc.Service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (bc *BookController) GetSimilarTitles(c *gin.Context) {
	title := c.Query("title")
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing title query parameter"})
		return
	}

	books, err := bc.Service.GetSimilarTitles(title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if books == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "books not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": books})
}

func (bc *BookController) FilterBooks(c *gin.Context) {
	var filter models.BookFilter
	if idStr := c.Query("id"); idStr != "" {
		if id, err := strconv.Atoi(idStr); err == nil {
			filter.ID = &id
		}
	}
	if isbn := c.Query("isbn"); isbn != "" {
		filter.ISBN = &isbn
	}
	if title := c.Query("title"); title != "" {
		filter.Title = &title
	}
	if tpStr := c.Query("total_pages"); tpStr != "" {
		if tp, err := strconv.Atoi(tpStr); err == nil {
			filter.TotalPages = &tp
		}
	}

	books, err := bc.Service.FilterBooks(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch books"})
		return
	}

	c.JSON(http.StatusOK, books)
}
