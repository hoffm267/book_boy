package bl

import (
	"book_boy/backend/internal/dl"
	"book_boy/backend/internal/models"
)

type UserService interface {
	GetAll() ([]models.User, error)
}

type userService struct {
	repo dl.UserRepo
}

func NewUserService(repo dl.UserRepo) UserService {
	return &userService{repo: repo}
}

func (s *userService) GetAll() ([]models.User, error) {
	return s.repo.GetAll()
}
