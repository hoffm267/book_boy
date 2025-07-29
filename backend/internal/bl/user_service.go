package bl

import (
	"book_boy/backend/internal/dl"
	"book_boy/backend/internal/models"
)

type UserService interface {
	GetAll() ([]models.User, error)
	GetByID(id int) (*models.User, error)
	Create(user *models.User) (int, error)
	Update(user *models.User) error
	Delete(id int) error
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

func (s *userService) GetByID(id int) (*models.User, error) {
	return s.repo.GetByID(id)
}

func (s *userService) Create(user *models.User) (int, error) {
	return s.repo.Create(user)
}

func (s *userService) Update(user *models.User) error {
	return s.repo.Update(user)
}

func (s *userService) Delete(id int) error {
	return s.repo.Delete(id)
}
