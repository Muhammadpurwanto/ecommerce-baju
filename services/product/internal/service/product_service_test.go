package service

import (
	"errors"
	"testing"

	"github.com/Muhammadpurwanto/ecommerce-baju/services/product/internal/dto"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/product/internal/model"
)

// Mock ProductRepository
type mockProductRepository struct {
	FindAllFunc      func(page, perPage int) ([]model.Product, int64, error)
	FindByIDFunc     func(id uint) (*model.Product, error)
	FindBySlugFunc   func(slug string) (*model.Product, error)
	CreateFunc       func(product *model.Product) error
	UpdateFunc       func(product *model.Product) error
	DeleteFunc       func(id uint) error
	ReserveStockFunc func(items []dto.StockItem) error
	ReleaseStockFunc func(items []dto.StockItem) error
}

func (m *mockProductRepository) FindAll(page, perPage int) ([]model.Product, int64, error) {
	return m.FindAllFunc(page, perPage)
}
func (m *mockProductRepository) FindByID(id uint) (*model.Product, error) {
	return m.FindByIDFunc(id)
}
func (m *mockProductRepository) FindBySlug(slug string) (*model.Product, error) {
	return m.FindBySlugFunc(slug)
}
func (m *mockProductRepository) Create(product *model.Product) error {
	return m.CreateFunc(product)
}
func (m *mockProductRepository) Update(product *model.Product) error {
	return m.UpdateFunc(product)
}
func (m *mockProductRepository) Delete(id uint) error {
	return m.DeleteFunc(id)
}
func (m *mockProductRepository) ReserveStock(items []dto.StockItem) error {
	return m.ReserveStockFunc(items)
}
func (m *mockProductRepository) ReleaseStock(items []dto.StockItem) error {
	return m.ReleaseStockFunc(items)
}

// Mock ProductCacheService
type mockProductCacheService struct {
	GetProductsFunc   func(page, perPage int) ([]model.Product, error)
	SetProductsFunc   func(page, perPage int, products []model.Product) error
	GetProductFunc    func(slugStr string) (*model.Product, error)
	SetProductFunc    func(slugStr string, product *model.Product) error
	InvalidateAllFunc func() error
}

func (m *mockProductCacheService) GetProducts(page, perPage int) ([]model.Product, error) {
	return m.GetProductsFunc(page, perPage)
}
func (m *mockProductCacheService) SetProducts(page, perPage int, products []model.Product) error {
	return m.SetProductsFunc(page, perPage, products)
}
func (m *mockProductCacheService) GetProduct(slugStr string) (*model.Product, error) {
	return m.GetProductFunc(slugStr)
}
func (m *mockProductCacheService) SetProduct(slugStr string, product *model.Product) error {
	return m.SetProductFunc(slugStr, product)
}
func (m *mockProductCacheService) InvalidateAll() error {
	return m.InvalidateAllFunc()
}

func TestProductService_GetByID(t *testing.T) {
	t.Run("Success Get Product By ID", func(t *testing.T) {
		mockRepo := &mockProductRepository{
			FindByIDFunc: func(id uint) (*model.Product, error) {
				return &model.Product{ID: 1, Name: "Flannel Shirt"}, nil
			},
		}

		srv := NewProductService(mockRepo, nil)
		resp, err := srv.GetByID(1)

		if err != nil {
			t.Fatalf("diharapkan tidak ada error, mendapat: %v", err)
		}
		if resp.Name != "Flannel Shirt" {
			t.Errorf("diharapkan nama 'Flannel Shirt', mendapat: %s", resp.Name)
		}
	})

	t.Run("Product Not Found", func(t *testing.T) {
		mockRepo := &mockProductRepository{
			FindByIDFunc: func(id uint) (*model.Product, error) {
				return nil, errors.New("not found")
			},
		}

		srv := NewProductService(mockRepo, nil)
		resp, err := srv.GetByID(99)

		if err == nil {
			t.Fatal("diharapkan error, mendapat nil")
		}
		if err.Error() != "product not found" {
			t.Errorf("diharapkan pesan error 'product not found', mendapat: %s", err.Error())
		}
		if resp != nil {
			t.Errorf("diharapkan response nil, mendapat: %+v", resp)
		}
	})
}

func TestProductService_Create(t *testing.T) {
	t.Run("Success Create Product", func(t *testing.T) {
		cacheInvalidated := false
		mockCache := &mockProductCacheService{
			InvalidateAllFunc: func() error {
				cacheInvalidated = true
				return nil
			},
		}

		mockRepo := &mockProductRepository{
			CreateFunc: func(product *model.Product) error {
				return nil
			},
		}

		srv := NewProductService(mockRepo, mockCache)
		req := &dto.ProductRequest{
			CategoryID: 1,
			Name:       "New Flannel Shirt",
			Brand:      "LocalBrand",
			BasePrice:  100000,
			Stock:      10,
		}

		resp, err := srv.Create(req)
		if err != nil {
			t.Fatalf("diharapkan tidak ada error, mendapat: %v", err)
		}
		if resp.Slug != "new-flannel-shirt" {
			t.Errorf("diharapkan slug 'new-flannel-shirt', mendapat: %s", resp.Slug)
		}
		if !cacheInvalidated {
			t.Error("diharapkan cache dibersihkan saat membuat produk baru")
		}
	})
}
