package service

import (
	"gin-template/internal/model"
	"gin-template/internal/repository"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetAll() ([]model.User, error) {
	return s.repo.FindAll()
}

func (s *UserService) GetByID(id string) (model.User, error) {
	return s.repo.FindByID(id)
}

func (s *UserService) Create(user *model.User) error {
	return s.repo.Create(user)
}

func (s *UserService) Update(id string, user *model.User) (model.User, error) {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return model.User{}, err
	}

	existing.Name = user.Name
	existing.Email = user.Email

	if err := s.repo.Save(&existing); err != nil {
		return model.User{}, err
	}

	return existing, nil
}
