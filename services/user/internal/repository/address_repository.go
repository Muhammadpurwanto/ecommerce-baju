package repository

import (
	"github.com/Muhammadpurwanto/ecommerce-baju/services/user/internal/model"
	"gorm.io/gorm"
)

type AddressRepository interface {
	FindByUserID(userID string) ([]model.Address, error)
	FindByID(id uint, userID string) (*model.Address, error)
	Create(address *model.Address) error
	Update(address *model.Address) error
	Delete(id uint, userID string) error
	ResetDefault(userID string) error
}

type addressRepository struct {
	db *gorm.DB
}

func NewAddressRepository(db *gorm.DB) AddressRepository {
	return &addressRepository{db: db}
}

func (r *addressRepository) FindByUserID(userID string) ([]model.Address, error) {
	var addresses []model.Address
	if err := r.db.Where("user_id = ?", userID).Order("is_default DESC, created_at DESC").Find(&addresses).Error; err != nil {
		return nil, err
	}
	return addresses, nil
}

func (r *addressRepository) FindByID(id uint, userID string) (*model.Address, error) {
	var address model.Address
	if err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&address).Error; err != nil {
		return nil, err
	}
	return &address, nil
}

func (r *addressRepository) Create(address *model.Address) error {
	return r.db.Create(address).Error
}

func (r *addressRepository) Update(address *model.Address) error {
	return r.db.Save(address).Error
}

func (r *addressRepository) Delete(id uint, userID string) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Address{}).Error
}

func (r *addressRepository) ResetDefault(userID string) error {
	return r.db.Model(&model.Address{}).Where("user_id = ?", userID).Update("is_default", false).Error
}
