package service

import (
	"errors"
	"testing"

	"gorm.io/gorm"

	"github.com/Muhammadpurwanto/ecommerce-baju/services/cart/internal/dto"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/cart/internal/model"
)

// Mock CartRepository
type mockCartRepository struct {
	GetCartByUserIDFunc func(userID string) (*model.Cart, error)
	CreateCartFunc      func(cart *model.Cart) error
	FindItemFunc        func(cartID, productID uint) (*model.CartItem, error)
	AddItemFunc         func(item *model.CartItem) error
	UpdateItemFunc      func(item *model.CartItem) error
	RemoveItemFunc      func(itemID uint) error
	ClearCartFunc       func(cartID uint) error
}

func (m *mockCartRepository) GetCartByUserID(userID string) (*model.Cart, error) {
	return m.GetCartByUserIDFunc(userID)
}
func (m *mockCartRepository) CreateCart(cart *model.Cart) error {
	return m.CreateCartFunc(cart)
}
func (m *mockCartRepository) FindItem(cartID, productID uint) (*model.CartItem, error) {
	return m.FindItemFunc(cartID, productID)
}
func (m *mockCartRepository) AddItem(item *model.CartItem) error {
	return m.AddItemFunc(item)
}
func (m *mockCartRepository) UpdateItem(item *model.CartItem) error {
	return m.UpdateItemFunc(item)
}
func (m *mockCartRepository) RemoveItem(itemID uint) error {
	return m.RemoveItemFunc(itemID)
}
func (m *mockCartRepository) ClearCart(cartID uint) error {
	return m.ClearCartFunc(cartID)
}

// Mock CartCacheService
type mockCartCacheService struct {
	GetCartFunc    func(userID string) (*model.Cart, error)
	SetCartFunc    func(userID string, cart *model.Cart) error
	DeleteCartFunc func(userID string) error
}

func (m *mockCartCacheService) GetCart(userID string) (*model.Cart, error) {
	return m.GetCartFunc(userID)
}
func (m *mockCartCacheService) SetCart(userID string, cart *model.Cart) error {
	return m.SetCartFunc(userID, cart)
}
func (m *mockCartCacheService) DeleteCart(userID string) error {
	return m.DeleteCartFunc(userID)
}

func TestCartService_GetCart(t *testing.T) {
	t.Run("Get Cart From Cache Success", func(t *testing.T) {
		mockCache := &mockCartCacheService{
			GetCartFunc: func(userID string) (*model.Cart, error) {
				return &model.Cart{ID: 1, UserID: "user-1"}, nil
			},
		}

		srv := NewCartService(nil, mockCache)
		resp, err := srv.GetCart("user-1")

		if err != nil {
			t.Fatalf("diharapkan tidak ada error, mendapat: %v", err)
		}
		if resp.UserID != "user-1" {
			t.Errorf("diharapkan UserID 'user-1', mendapat: %s", resp.UserID)
		}
	})

	t.Run("Get Cart From DB and Save to Cache", func(t *testing.T) {
		cacheSaved := false
		mockCache := &mockCartCacheService{
			GetCartFunc: func(userID string) (*model.Cart, error) {
				return nil, errors.New("cache miss")
			},
			SetCartFunc: func(userID string, cart *model.Cart) error {
				cacheSaved = true
				return nil
			},
		}

		mockRepo := &mockCartRepository{
			GetCartByUserIDFunc: func(userID string) (*model.Cart, error) {
				return &model.Cart{ID: 1, UserID: "user-1"}, nil
			},
		}

		srv := NewCartService(mockRepo, mockCache)
		resp, err := srv.GetCart("user-1")

		if err != nil {
			t.Fatalf("diharapkan tidak ada error, mendapat: %v", err)
		}
		if resp.ID != 1 {
			t.Errorf("diharapkan ID 1, mendapat: %d", resp.ID)
		}
		if !cacheSaved {
			t.Error("diharapkan data disimpan ke cache saat terjadi cache miss")
		}
	})
}

func TestCartService_AddItem(t *testing.T) {
	t.Run("Add New Item to Cart", func(t *testing.T) {
		mockCache := &mockCartCacheService{
			DeleteCartFunc: func(userID string) error {
				return nil
			},
		}

		itemAdded := false
		mockRepo := &mockCartRepository{
			GetCartByUserIDFunc: func(userID string) (*model.Cart, error) {
				return &model.Cart{ID: 1, UserID: "user-1"}, nil
			},
			FindItemFunc: func(cartID, productID uint) (*model.CartItem, error) {
				return nil, gorm.ErrRecordNotFound
			},
			AddItemFunc: func(item *model.CartItem) error {
				if item.ProductID != 101 {
					t.Errorf("diharapkan ProductID 101, mendapat: %d", item.ProductID)
				}
				itemAdded = true
				return nil
			},
		}

		srv := NewCartService(mockRepo, mockCache)
		req := &dto.AddItemRequest{
			ProductID: 101,
			Quantity:  2,
		}

		_, err := srv.AddItem("user-1", req)
		if err != nil {
			t.Fatalf("diharapkan sukses, mendapat error: %v", err)
		}
		if !itemAdded {
			t.Error("diharapkan AddItem terpanggil")
		}
	})
}
