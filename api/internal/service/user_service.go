package service

import (
	"book_boy/api/internal/domain"
	"book_boy/api/internal/repository"
)

type UserService interface {
	GetAll() ([]domain.User, error)
	GetByID(id int) (*domain.User, error)
	Create(user *domain.User) (int, error)
	Update(user *domain.User) error
	Delete(id int) error
}

type userService struct {
	repo repository.UserRepo
}

func NewUserService(repo repository.UserRepo) UserService {
	return &userService{repo: repo}
}

func (s *userService) GetAll() ([]domain.User, error) {
	return s.repo.GetAll()
}

func (s *userService) GetByID(id int) (*domain.User, error) {
	return s.repo.GetByID(id)
}

func (s *userService) Create(user *domain.User) (int, error) {
	return s.repo.Create(user)
}

func (s *userService) Update(user *domain.User) error {
	return s.repo.Update(user)
}

func (s *userService) Delete(id int) error {
	return s.repo.Delete(id)
}
