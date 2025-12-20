package main

import (
	"fmt"
	"log"
	"os"

	"book_boy/api/external/book_metadata"
	"book_boy/api/internal/controllers"
	"book_boy/api/internal/db"
	"book_boy/api/internal/infra"
	"book_boy/api/internal/middleware"
	"book_boy/api/internal/repository"
	"book_boy/api/internal/service"
	"book_boy/api/internal/workers"

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

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "localhost:6379"
	}
	cache := infra.NewCache(redisURL)
	fmt.Println("Connected to Redis:", redisURL)

	rabbitmqURL := os.Getenv("RABBITMQ_URL")
	if rabbitmqURL == "" {
		rabbitmqURL = "amqp://guest:guest@localhost:5672/"
	}

	rabbitConn, err := infra.ConnectRabbitMQ(rabbitmqURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitConn.Close()
	fmt.Println("Connected to RabbitMQ:", rabbitmqURL)

	publisher, err := infra.NewEventPublisher(rabbitConn, "book_events")
	if err != nil {
		log.Fatalf("Failed to create event publisher: %v", err)
	}

	sseManager := infra.NewSSEManager()

	bookRepo := repository.NewBookRepo(database)
	bookService := service.NewBookService(bookRepo, cache, publisher)

	userRepo := repository.NewUserRepo(database)
	userService := service.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	authService := service.NewAuthService(userRepo)
	authController := controllers.NewAuthController(authService)

	audiobookRepo := repository.NewAudiobookRepo(database)
	audiobookService := service.NewAudiobookService(audiobookRepo, cache)

	progressRepo := repository.NewProgressRepo(database)
	progressService := service.NewProgressService(progressRepo)

	trackingService := service.NewTrackingService(bookRepo, audiobookRepo, progressRepo)

	bookMetadataServiceURL := os.Getenv("BOOK_METADATA_SERVICE_URL")
	if bookMetadataServiceURL == "" {
		bookMetadataServiceURL = "http://localhost:8000"
	}
	bookMetadataClient := book_metadata.NewClient(bookMetadataServiceURL)
	bookMetadataService := service.NewBookMetadataService(bookMetadataClient, cache)

	metadataConsumer := workers.NewMetadataEventConsumer(rabbitConn, bookService, sseManager)
	if err := metadataConsumer.Start(); err != nil {
		log.Fatalf("Failed to start metadata event consumer: %v", err)
	}
	fmt.Println("Started metadata event consumer")

	bookController := controllers.NewBookController(bookService, progressService, bookMetadataService)
	audiobookController := controllers.NewAudiobookController(audiobookService, progressService)
	progressController := controllers.NewProgressController(progressService)
	trackingController := controllers.NewTrackingController(trackingService)
	bookMetadataController := controllers.NewBookMetadataController(bookMetadataService)

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		allowedOrigins := map[string]bool{
			"https://bookboy.app":     true,
			"https://www.bookboy.app": true,
			"https://api.bookboy.app": true,
			"http://localhost:5173":   true,
			"http://localhost:8080":   true,
		}

		if allowedOrigins[origin] {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}

		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

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
		bookMetadataController.RegisterRoutes(protected)
		bookController.RegisterRoutes(protected)
		audiobookController.RegisterRoutes(protected)
		userController.RegisterRoutes(protected)
		progressController.RegisterRoutes(protected)
		trackingController.RegisterRoutes(protected)
	}

	r.GET("/events", func(c *gin.Context) {
		token := c.Query("token")
		if token == "" {
			c.JSON(401, gin.H{"error": "token required"})
			return
		}

		_, err := authService.GetUserFromToken(token)
		if err != nil {
			c.JSON(401, gin.H{"error": "invalid token"})
			return
		}

		sseManager.ServeHTTP(c)
	})

	r.Run(":8080")
}
