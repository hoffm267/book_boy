package controllers

import (
	"book_boy/api/internal/service"
	"book_boy/api/internal/domain"
	"net/http"

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
