package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/hoffm267/book_boy/backend/internal/handler"
)

func main() {
	fmt.Println("Server started.\n")

	r := gin.Default()
	r.GET("/ping", handler.GetTest)

	r.Run()
}
