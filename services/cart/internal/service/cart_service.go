package service

import (
	"errors"

	"gorm.io/gorm"

	"github.com/Muhammadpurwanto/ecommerce-baju/services/cart/internal/dto"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/cart/internal/model"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/cart/internal/repository"
)

type CartService interface {
	GetCart(userID string) (*model.Cart, error)
	AddItem(userID string, req *dto.AddItemRequest) (*model.Cart, error)
	UpdateItem(userID string, itemID uint, req *dto.UpdateItemRequest) (*model.Cart, error)
	RemoveItem(userID string, itemID uint) (*model.Cart, error)
	ClearCart(userID string) error
}

type cartService struct {
	repo   repository.CartRepository
	cache  CartCacheService
}

func NewCartService(repo repository.CartRepository, cache CartCacheService) CartService {
	return &cartService{repo: repo, cache: cache}
}

// Internal helper to get or create cart
func (s *cartService) getOrCreateCart(userID string) (*model.Cart, error) {
	cart, err := s.repo.GetCartByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newCart := &model.Cart{UserID: userID}
			if err := s.repo.CreateCart(newCart); err != nil {
				return nil, err
			}
			return newCart, nil
		}
		return nil, err
	}
	return cart, nil
}

func (s *cartService) GetCart(userID string) (*model.Cart, error) {
	// 1. Coba ambil dari cache dulu
	cached, err := s.cache.GetCart(userID)
	if err == nil {
		return cached, nil
	}
	// 2. Jika tidak ada di cache, ambil dari DB
	cart, err := s.getOrCreateCart(userID)
	if err != nil {
		return nil, err
	}
	// 3. Simpan ke cache
	_ = s.cache.SetCart(userID, cart)
	return cart, nil
}

func (s *cartService) AddItem(userID string, req *dto.AddItemRequest) (*model.Cart, error) {
	cart, err := s.getOrCreateCart(userID)
	if err != nil {
		return nil, err
	}

	// Cek jika item sudah ada, maka update quantity
	item, err := s.repo.FindItem(cart.ID, req.ProductID)
	if err == nil {
		item.Quantity += req.Quantity
		if err := s.repo.UpdateItem(item); err != nil {
			return nil, err
		}
	} else {
		newItem := &model.CartItem{
			CartID:    cart.ID,
			ProductID: req.ProductID,
			Quantity:  req.Quantity,
		}
		if err := s.repo.AddItem(newItem); err != nil {
			return nil, err
		}
	}

	_ = s.cache.DeleteCart(userID)
	// Kembalikan state cart terbaru
	return s.repo.GetCartByUserID(userID)
}

func (s *cartService) UpdateItem(userID string, itemID uint, req *dto.UpdateItemRequest) (*model.Cart, error) {
	// Untuk keamanan, pastikan cart milik user
	_, err := s.getOrCreateCart(userID)
	if err != nil {
		return nil, err
	}

	// Anggap itemID ini valid punya user, idealnya join check
	item := &model.CartItem{ID: itemID, Quantity: req.Quantity}
	if err := s.repo.UpdateItem(item); err != nil {
		return nil, err
	}

	_ = s.cache.DeleteCart(userID)
	return s.repo.GetCartByUserID(userID)
}

func (s *cartService) RemoveItem(userID string, itemID uint) (*model.Cart, error) {
	_, err := s.getOrCreateCart(userID)
	if err != nil {
		return nil, err
	}

	if err := s.repo.RemoveItem(itemID); err != nil {
		return nil, err
	}

	_ = s.cache.DeleteCart(userID)
	return s.repo.GetCartByUserID(userID)
}

func (s *cartService) ClearCart(userID string) error {
	cart, err := s.getOrCreateCart(userID)
	if err != nil {
		return err
	}

	_ = s.cache.DeleteCart(userID)
	return s.repo.ClearCart(cart.ID)
}
