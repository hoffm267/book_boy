package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"

	"book_boy/backend/internal/bl"
	"book_boy/backend/internal/controllers"
	"book_boy/backend/internal/db"
	"book_boy/backend/internal/dl"
)

func main() {
	fmt.Println("Server started...")

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	database := db.InitDB(connStr)

	bookRepo := dl.NewBookRepo(database)
	bookService := bl.NewBookService(bookRepo)
	bookController := controllers.NewBookController(bookService)

	userRepo := dl.NewUserRepo(database)
	userService := bl.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	audiobookRepo := dl.NewAudiobookRepo(database)
	audiobookService := bl.NewAudiobookService(audiobookRepo)
	audiobookController := controllers.NewAudiobookController(audiobookService)

	progressRepo := dl.NewProgressRepo(database)
	progressService := bl.NewProgressService(progressRepo)
	progressController := controllers.NewProgressController(progressService)

	r := gin.Default()
	r.GET("/ping", controllers.GetTest)
	r.GET("/books", bookController.GetAll)
	r.GET("/books/:id", bookController.GetByID)
	r.GET("/users", userController.GetAll)
	r.GET("/audiobooks", audiobookController.GetAll)
	r.GET("/progress", progressController.GetAll)
	r.Run(":8080")
}
