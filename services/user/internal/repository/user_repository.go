package repository

import (
	"github.com/Muhammadpurwanto/ecommerce-baju/services/user/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByID(id string) (*model.User, error)
	FindByEmail(email string) (*model.User, error)
	FindAll(page, perPage int) ([]model.User, int64, error)
	Update(user *model.User) error
	Create(user *model.User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByID(id string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindAll(page, perPage int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	r.db.Model(&model.User{}).Count(&total)

	offset := (page - 1) * perPage
	if err := r.db.Offset(offset).Limit(perPage).Order("created_at DESC").Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}
