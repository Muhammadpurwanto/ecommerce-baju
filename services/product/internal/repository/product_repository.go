package repository

import (
	"fmt"

	"github.com/Muhammadpurwanto/ecommerce-baju/services/product/internal/dto"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/product/internal/model"
	"gorm.io/gorm"
)

type ProductRepository interface {
	FindAll(page, perPage int) ([]model.Product, int64, error)
	FindByID(id uint) (*model.Product, error)
	FindBySlug(slug string) (*model.Product, error)
	Create(product *model.Product) error
	Update(product *model.Product) error
	Delete(id uint) error
	ReserveStock(items []dto.StockItem) error
	ReleaseStock(items []dto.StockItem) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) FindAll(page, perPage int) ([]model.Product, int64, error) {
	var products []model.Product
	var total int64

	r.db.Model(&model.Product{}).Count(&total)
	offset := (page - 1) * perPage

	err := r.db.Preload("Category").
		Offset(offset).Limit(perPage).Order("created_at DESC").Find(&products).Error
	
	return products, total, err
}

func (r *productRepository) FindByID(id uint) (*model.Product, error) {
	var product model.Product
	err := r.db.Preload("Category").First(&product, id).Error
	return &product, err
}

func (r *productRepository) FindBySlug(slug string) (*model.Product, error) {
	var product model.Product
	err := r.db.Preload("Category").Where("slug = ?", slug).First(&product).Error
	return &product, err
}

func (r *productRepository) Create(product *model.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) Update(product *model.Product) error {
	return r.db.Omit("Category").Save(product).Error
}

func (r *productRepository) Delete(id uint) error {
	return r.db.Delete(&model.Product{}, id).Error
}

func (r *productRepository) ReserveStock(items []dto.StockItem) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, item := range items {
			result := tx.Model(&model.Product{}).
				Where("id = ? AND stock >= ?", item.ProductID, item.Quantity).
				Update("stock", gorm.Expr("stock - ?", item.Quantity))
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				var prod model.Product
				if err := tx.Select("name, stock").First(&prod, item.ProductID).Error; err == nil {
					return fmt.Errorf("insufficient stock for product %s (requested: %d, available: %d)", prod.Name, item.Quantity, prod.Stock)
				}
				return fmt.Errorf("product with ID %d not found or insufficient stock", item.ProductID)
			}
		}
		return nil
	})
}

func (r *productRepository) ReleaseStock(items []dto.StockItem) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, item := range items {
			result := tx.Model(&model.Product{}).
				Where("id = ?", item.ProductID).
				Update("stock", gorm.Expr("stock + ?", item.Quantity))
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				return fmt.Errorf("product with ID %d not found", item.ProductID)
			}
		}
		return nil
	})
}
