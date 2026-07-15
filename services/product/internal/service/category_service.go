package service

import (
	"errors"

	"github.com/gosimple/slug"
	"gorm.io/gorm"

	"github.com/Muhammadpurwanto/ecommerce-baju/services/product/internal/dto"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/product/internal/model"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/product/internal/repository"
)

type CategoryService interface {
	GetAll() ([]model.Category, error)
	GetByID(id uint) (*model.Category, error)
	Create(req *dto.CategoryRequest) (*model.Category, error)
	Update(id uint, req *dto.CategoryRequest) (*model.Category, error)
	Delete(id uint) error
}

type categoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) GetAll() ([]model.Category, error) {
	return s.repo.FindAll()
}

func (s *categoryService) GetByID(id uint) (*model.Category, error) {
	cat, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}
	return cat, nil
}

func (s *categoryService) Create(req *dto.CategoryRequest) (*model.Category, error) {
	cat := &model.Category{
		Name: req.Name,
		Slug: slug.Make(req.Name),
	}
	if err := s.repo.Create(cat); err != nil {
		return nil, err
	}
	return cat, nil
}

func (s *categoryService) Update(id uint, req *dto.CategoryRequest) (*model.Category, error) {
	cat, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}

	cat.Name = req.Name
	cat.Slug = slug.Make(req.Name)

	if err := s.repo.Update(cat); err != nil {
		return nil, err
	}
	return cat, nil
}

func (s *categoryService) Delete(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("category not found")
	}
	return s.repo.Delete(id)
}
