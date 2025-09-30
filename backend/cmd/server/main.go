package main

import (
	"fmt"
	"os"

	"book_boy/backend/internal/bl"
	"book_boy/backend/internal/controllers"
	"book_boy/backend/internal/db"
	"book_boy/backend/internal/dl"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Server started...")

	_ = godotenv.Load()

	if os.Getenv("DB_HOST") == "" || os.Getenv("DB_PORT") == "" {
		panic("Missing DB_HOST or DB_PORT env vars")
	}
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

	userRepo := dl.NewUserRepo(database)
	userService := bl.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	audiobookRepo := dl.NewAudiobookRepo(database)
	audiobookService := bl.NewAudiobookService(audiobookRepo)

	progressRepo := dl.NewProgressRepo(database)
	progressService := bl.NewProgressService(progressRepo)

	r := gin.Default()

	bookController := controllers.NewBookController(bookService, progressService)
	audiobookController := controllers.NewAudiobookController(audiobookService, progressService)
	progressController := controllers.NewProgressController(progressService)
	bookController.RegisterRoutes(r)

	audiobookController.RegisterRoutes(r)
	userController.RegisterRoutes(r)
	progressController.RegisterRoutes(r)

	r.GET("/ping", controllers.GetTest)
	r.Run(":8080")
}
