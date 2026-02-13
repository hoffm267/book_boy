package controllers

import (
	"book_boy/api/internal/domain"
	"book_boy/api/internal/repository"
	"book_boy/api/internal/service"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type BookController struct {
	Service         service.BookService
	ProgressService service.ProgressService
}

func NewBookController(service service.BookService, pgService service.ProgressService) *BookController {
	return &BookController{
		Service:         service,
		ProgressService: pgService,
	}
}

func (bc *BookController) RegisterRoutes(r gin.IRouter) {
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
	var book domain.Book
	var progress domain.Progress

	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book.ISBN = strings.ReplaceAll(strings.ReplaceAll(book.ISBN, "-", ""), " ", "")

	if book.ISBN == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "isbn is required"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	filter := repository.BookFilter{ISBN: &book.ISBN}
	existingBooks, err := bc.Service.FilterBooks(filter)

	var id int
	if err == nil && len(existingBooks) > 0 {
		id = existingBooks[0].ID
	} else {
		id, err = bc.Service.Create(&book)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	uid := userID.(int)

	skipProgress := c.Query("skipProgress") == "true"
	if skipProgress {
		savedBook, err := bc.Service.GetByID(id)
		if err != nil || savedBook == nil {
			book.ID = id
			c.JSON(http.StatusCreated, gin.H{"data": book})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"data": savedBook})
		return
	}

	pgIDStr := c.Query("pgId")
	if pgIDStr != "" {
		pgID, err := strconv.Atoi(pgIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid book id"})
			return
		}

		if err := bc.ProgressService.SetBook(pgID, id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		progFilter := repository.ProgressFilter{
			UserID: &uid,
			BookID: &id,
		}
		existingProgress, err := bc.ProgressService.FilterProgress(progFilter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if len(existingProgress) > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "progress already exists for this book"})
			return
		}

		page := 1
		progress = domain.Progress{
			UserID:   uid,
			BookID:   &id,
			BookPage: &page,
		}

		_, err = bc.ProgressService.Create(&progress)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	savedBook, err := bc.Service.GetByID(id)
	if err != nil || savedBook == nil {
		book.ID = id
		c.JSON(http.StatusCreated, gin.H{"data": book})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": savedBook})
}

func (bc *BookController) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid book ID"})
		return
	}

	var book domain.Book
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
	var filter repository.BookFilter

	if idStr := c.Query("id"); idStr != "" {
		if id, err := strconv.Atoi(idStr); err == nil {
			filter.ID = &id
		}
	}
	if isbn := c.Query("isbn"); isbn != "" {
		cleaned := strings.ReplaceAll(strings.ReplaceAll(isbn, "-", ""), " ", "")
		filter.ISBN = &cleaned
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
