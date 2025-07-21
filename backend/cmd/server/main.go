package main

import (
	"fmt"

	"book_boy/backend/internal/bl"
	"book_boy/backend/internal/controllers"
	"book_boy/backend/internal/db"
	"book_boy/backend/internal/dl"

	"github.com/gin-gonic/gin"
)

//TODO change to use env variables
/*
import "fmt"
import "os"

connStr := fmt.Sprintf(

	"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
	os.Getenv("DB_HOST"),
	os.Getenv("DB_PORT"),
	os.Getenv("DB_USER"),
	os.Getenv("DB_PASSWORD"),
	os.Getenv("DB_NAME"),

)
*/
func main() {
	fmt.Println("Server started.")

	connStr := "host=db port=5432 user=postgres password=postgres dbname=postgres sslmode=disable"
	database := db.InitDB(connStr)

	bookRepo := dl.NewBookRepo(database)
	bookService := bl.NewBookService(bookRepo)
	bookController := controllers.NewBookController(bookService)

	r := gin.Default()
	r.GET("/ping", controllers.GetTest)
	r.GET("/books", bookController.GetAll)
	r.Run(":8080")
}
