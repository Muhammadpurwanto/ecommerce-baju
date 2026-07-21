package service

import (
	"testing"

	"gorm.io/gorm"

	"github.com/Muhammadpurwanto/ecommerce-baju/services/product/internal/dto"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/product/internal/model"
)

// Mock CategoryRepository
type mockCategoryRepository struct {
	FindAllFunc    func() ([]model.Category, error)
	FindByIDFunc   func(id uint) (*model.Category, error)
	CreateFunc     func(category *model.Category) error
	UpdateFunc     func(category *model.Category) error
	DeleteFunc     func(id uint) error
}

func (m *mockCategoryRepository) FindAll() ([]model.Category, error) {
	return m.FindAllFunc()
}
func (m *mockCategoryRepository) FindByID(id uint) (*model.Category, error) {
	return m.FindByIDFunc(id)
}
func (m *mockCategoryRepository) Create(category *model.Category) error {
	return m.CreateFunc(category)
}
func (m *mockCategoryRepository) Update(category *model.Category) error {
	return m.UpdateFunc(category)
}
func (m *mockCategoryRepository) Delete(id uint) error {
	return m.DeleteFunc(id)
}

func TestCategoryService_GetByID(t *testing.T) {
	t.Run("Success Get Category By ID", func(t *testing.T) {
		mockRepo := &mockCategoryRepository{
			FindByIDFunc: func(id uint) (*model.Category, error) {
				return &model.Category{ID: 1, Name: "Shirts"}, nil
			},
		}

		srv := NewCategoryService(mockRepo)
		resp, err := srv.GetByID(1)

		if err != nil {
			t.Fatalf("diharapkan tidak ada error, mendapat: %v", err)
		}
		if resp.Name != "Shirts" {
			t.Errorf("diharapkan nama 'Shirts', mendapat: %s", resp.Name)
		}
	})

	t.Run("Category Not Found", func(t *testing.T) {
		mockRepo := &mockCategoryRepository{
			FindByIDFunc: func(id uint) (*model.Category, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		srv := NewCategoryService(mockRepo)
		resp, err := srv.GetByID(99)

		if err == nil {
			t.Fatal("diharapkan error, mendapat nil")
		}
		if err.Error() != "category not found" {
			t.Errorf("diharapkan error 'category not found', mendapat: %s", err.Error())
		}
		if resp != nil {
			t.Errorf("diharapkan response nil, mendapat: %+v", resp)
		}
	})
}

func TestCategoryService_Create(t *testing.T) {
	t.Run("Success Create Category", func(t *testing.T) {
		mockRepo := &mockCategoryRepository{
			CreateFunc: func(category *model.Category) error {
				return nil
			},
		}

		srv := NewCategoryService(mockRepo)
		req := &dto.CategoryRequest{Name: "Men Clothes"}

		resp, err := srv.Create(req)
		if err != nil {
			t.Fatalf("diharapkan sukses, mendapat error: %v", err)
		}
		if resp.Slug != "men-clothes" {
			t.Errorf("diharapkan slug 'men-clothes', mendapat: %s", resp.Slug)
		}
	})
}
