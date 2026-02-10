package controllers

import (
	"book_boy/api/internal/domain"
	"book_boy/api/internal/service"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	Service service.AuthService
}

func NewAuthController(service service.AuthService) *AuthController {
	return &AuthController{Service: service}
}

func (ac *AuthController) RegisterRoutes(r gin.IRouter) {
	auth := r.Group("/auth")
	{
		auth.POST("/register", ac.Register)
		auth.POST("/login", ac.Login)
		auth.POST("/demo", ac.DemoLogin)
	}
}

func (ac *AuthController) Register(c *gin.Context) {
	var req domain.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := ac.Service.Register(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	loginReq := &domain.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}
	token, _, err := ac.Service.Login(loginReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "registration successful but login failed"})
		return
	}

	response := domain.AuthResponse{
		Token: token,
		User:  *user,
	}

	c.JSON(http.StatusCreated, response)
}

func (ac *AuthController) Login(c *gin.Context) {
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, user, err := ac.Service.Login(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	response := domain.AuthResponse{
		Token: token,
		User:  *user,
	}

	c.JSON(http.StatusOK, response)
}

func (ac *AuthController) DemoLogin(c *gin.Context) {
	demoEmail := os.Getenv("DEMO_USER_EMAIL")
	demoPassword := os.Getenv("DEMO_USER_PASSWORD")
	if demoEmail == "" || demoPassword == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "demo account not available"})
		return
	}

	token, user, err := ac.Service.Login(&domain.LoginRequest{
		Email:    demoEmail,
		Password: demoPassword,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "demo login failed"})
		return
	}

	c.JSON(http.StatusOK, domain.AuthResponse{
		Token: token,
		User:  *user,
	})
}
