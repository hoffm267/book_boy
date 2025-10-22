package main

import (
	"fmt"
	"os"

	"book_boy/internal/bl"
	"book_boy/internal/controllers"
	"book_boy/internal/db"
	"book_boy/internal/dl"
	"book_boy/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Server started...")

	_ = godotenv.Load()

	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	if os.Getenv("DB_HOST") == "" || os.Getenv("DB_PORT") == "" {
		panic("Missing DB_HOST or DB_PORT env vars")
	}
	sslmode := os.Getenv("DB_SSLMODE")
	if sslmode == "" {
		sslmode = "disable"
	}
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		sslmode,
	)

	database := db.InitDB(connStr)

	fmt.Println("Running database migrations...")
	if err := db.RunMigrations(database); err != nil {
		panic(fmt.Sprintf("Failed to run migrations: %v", err))
	}

	bookRepo := dl.NewBookRepo(database)
	bookService := bl.NewBookService(bookRepo)

	userRepo := dl.NewUserRepo(database)
	userService := bl.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	authService := bl.NewAuthService(userRepo)
	authController := controllers.NewAuthController(authService)

	audiobookRepo := dl.NewAudiobookRepo(database)
	audiobookService := bl.NewAudiobookService(audiobookRepo)

	progressRepo := dl.NewProgressRepo(database)
	progressService := bl.NewProgressService(progressRepo)

	trackingService := bl.NewTrackingService(bookRepo, audiobookRepo, progressRepo)

	bookController := controllers.NewBookController(bookService, progressService)
	audiobookController := controllers.NewAudiobookController(audiobookService, progressService)
	progressController := controllers.NewProgressController(progressService)
	trackingController := controllers.NewTrackingController(trackingService)

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	authController.RegisterRoutes(r)

	protected := r.Group("")
	protected.Use(middleware.AuthMiddleware(authService))
	{
		bookController.RegisterRoutes(protected)
		audiobookController.RegisterRoutes(protected)
		userController.RegisterRoutes(protected)
		progressController.RegisterRoutes(protected)
		trackingController.RegisterRoutes(protected)
	}

	r.Run(":8080")
}
