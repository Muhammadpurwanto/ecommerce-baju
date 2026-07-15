package service

import (
	"errors"
	"math"

	"gorm.io/gorm"

	"github.com/Muhammadpurwanto/ecommerce-baju/services/user/internal/dto"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/user/internal/model"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/user/internal/repository"
)

type UserService interface {
	GetProfile(userID string) (*dto.UserResponse, error)
	UpdateProfile(userID string, req *dto.UpdateUserRequest) (*dto.UserResponse, error)
	GetAllUsers(page, perPage int) ([]dto.UserResponse, *dto.MetaData, error)
	GetUserByID(id string) (*dto.UserResponse, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) GetProfile(userID string) (*dto.UserResponse, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return toUserResponse(user), nil
}

func (s *userService) UpdateProfile(userID string, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Phone != nil {
		user.Phone = req.Phone
	}
	if req.AvatarURL != nil {
		user.AvatarURL = req.AvatarURL
	}

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	return toUserResponse(user), nil
}

func (s *userService) GetAllUsers(page, perPage int) ([]dto.UserResponse, *dto.MetaData, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	users, total, err := s.repo.FindAll(page, perPage)
	if err != nil {
		return nil, nil, err
	}

	var responses []dto.UserResponse
	for _, u := range users {
		responses = append(responses, *toUserResponse(&u))
	}

	meta := &dto.MetaData{
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: int(math.Ceil(float64(total) / float64(perPage))),
	}

	return responses, meta, nil
}

func (s *userService) GetUserByID(id string) (*dto.UserResponse, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return toUserResponse(user), nil
}

func toUserResponse(user *model.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Phone:     user.Phone,
		AvatarURL: user.AvatarURL,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
