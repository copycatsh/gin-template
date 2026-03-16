package repository

import (
	"gin-template/internal/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindAll() ([]model.User, error)
	FindByID(id string) (model.User, error)
	Create(user *model.User) error
	Save(user *model.User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindAll() ([]model.User, error) {
	var users []model.User
	err := r.db.Find(&users).Error
	return users, err
}

func (r *userRepository) FindByID(id string) (model.User, error) {
	var user model.User
	err := r.db.First(&user, "id = ?", id).Error
	return user, err
}

func (r *userRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) Save(user *model.User) error {
	return r.db.Save(user).Error
}
