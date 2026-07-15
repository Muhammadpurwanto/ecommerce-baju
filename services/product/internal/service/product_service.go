package service

import (
	"errors"

	"github.com/gosimple/slug"

	"github.com/Muhammadpurwanto/ecommerce-baju/services/product/internal/dto"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/product/internal/model"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/product/internal/repository"
)

type ProductService interface {
	GetAll(page, perPage int) ([]model.Product, error)
	GetBySlug(slug string) (*model.Product, error)
	GetByID(id uint) (*model.Product, error)
	Create(req *dto.ProductRequest) (*model.Product, error)
	Update(id uint, req *dto.ProductRequest) (*model.Product, error)
	Delete(id uint) error
	ReserveStock(items []dto.StockItem) error
	ReleaseStock(items []dto.StockItem) error
}

type productService struct {
	repo repository.ProductRepository
	cache ProductCacheService
}

func NewProductService(repo repository.ProductRepository, cache ProductCacheService) ProductService {
	return &productService{repo: repo, cache: cache}
}

func (s *productService) GetAll(page, perPage int) ([]model.Product, error) {
	cached, err := s.cache.GetProducts(page, perPage)
	if err == nil {
		return cached, nil
	}
	products, _, err := s.repo.FindAll(page, perPage)
	if err != nil {
		return nil, err
	}
	_ = s.cache.SetProducts(page, perPage, products)
	return products, nil
}

func (s *productService) GetBySlug(slugStr string) (*model.Product, error) {
	cached, err := s.cache.GetProduct(slugStr)
	if err == nil {
		return cached, nil
	}
	product, err := s.repo.FindBySlug(slugStr)
	if err != nil {
		return nil, err
	}
	_ = s.cache.SetProduct(slugStr, product)
	return product, nil
}

func (s *productService) GetByID(id uint) (*model.Product, error) {
	product, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("product not found")
	}
	return product, nil
}

func (s *productService) Create(req *dto.ProductRequest) (*model.Product, error) {
	productSlug := slug.Make(req.Name)

	product := &model.Product{
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Slug:        productSlug,
		Description: req.Description,
		Brand:       req.Brand,
		Gender:      req.Gender,
		BasePrice:   req.BasePrice,
		Weight:      req.Weight,
		Stock:       req.Stock,
		ImageURL:    req.ImageURL,
		IsActive:    req.IsActive,
	}

	if err := s.repo.Create(product); err != nil {
		return nil, err
	}

	_ = s.cache.InvalidateAll()
	return product, nil
}

func (s *productService) Update(id uint, req *dto.ProductRequest) (*model.Product, error) {
	product, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("product not found")
	}
	product.CategoryID = req.CategoryID
	product.Name = req.Name
	product.Slug = slug.Make(req.Name)
	product.Description = req.Description
	product.Brand = req.Brand
	product.Gender = req.Gender
	product.BasePrice = req.BasePrice
	product.Weight = req.Weight
	product.Stock = req.Stock
	product.ImageURL = req.ImageURL
	product.IsActive = req.IsActive

	if err := s.repo.Update(product); err != nil {
		return nil, err
	}
	
	_ = s.cache.InvalidateAll()
	return product, nil
}

func (s *productService) Delete(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("product not found")
	}

	_ = s.cache.InvalidateAll()
	return s.repo.Delete(id)
}

func (s *productService) ReserveStock(items []dto.StockItem) error {
	if err := s.repo.ReserveStock(items); err != nil {
		return err
	}
	_ = s.cache.InvalidateAll()
	return nil
}

func (s *productService) ReleaseStock(items []dto.StockItem) error {
	if err := s.repo.ReleaseStock(items); err != nil {
		return err
	}
	_ = s.cache.InvalidateAll()
	return nil
}