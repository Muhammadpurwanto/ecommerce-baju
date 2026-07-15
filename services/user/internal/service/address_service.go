package service

import (
	"errors"

	"gorm.io/gorm"

	"github.com/Muhammadpurwanto/ecommerce-baju/services/user/internal/dto"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/user/internal/model"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/user/internal/repository"
)

type AddressService interface {
	GetAddresses(userID string) ([]dto.AddressResponse, error)
	GetAddressByID(id uint, userID string) (*dto.AddressResponse, error)
	CreateAddress(userID string, req *dto.CreateAddressRequest) (*dto.AddressResponse, error)
	UpdateAddress(id uint, userID string, req *dto.UpdateAddressRequest) (*dto.AddressResponse, error)
	DeleteAddress(id uint, userID string) error
	SetDefaultAddress(id uint, userID string) error
}

type addressService struct {
	repo repository.AddressRepository
}

func NewAddressService(repo repository.AddressRepository) AddressService {
	return &addressService{repo: repo}
}

func (s *addressService) GetAddresses(userID string) ([]dto.AddressResponse, error) {
	addresses, err := s.repo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	var responses []dto.AddressResponse
	for _, a := range addresses {
		responses = append(responses, *toAddressResponse(&a))
	}

	return responses, nil
}

func (s *addressService) GetAddressByID(id uint, userID string) (*dto.AddressResponse, error) {
	address, err := s.repo.FindByID(id, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("address not found")
		}
		return nil, err
	}

	return toAddressResponse(address), nil
}

func (s *addressService) CreateAddress(userID string, req *dto.CreateAddressRequest) (*dto.AddressResponse, error) {
	if req.IsDefault {
		if err := s.repo.ResetDefault(userID); err != nil {
			return nil, err
		}
	}

	address := &model.Address{
		UserID:     userID,
		Label:      req.Label,
		Recipient:  req.Recipient,
		Phone:      req.Phone,
		Province:   req.Province,
		City:       req.City,
		District:   req.District,
		PostalCode: req.PostalCode,
		Detail:     req.Detail,
		IsDefault:  req.IsDefault,
	}

	if err := s.repo.Create(address); err != nil {
		return nil, err
	}

	return toAddressResponse(address), nil
}

func (s *addressService) UpdateAddress(id uint, userID string, req *dto.UpdateAddressRequest) (*dto.AddressResponse, error) {
	address, err := s.repo.FindByID(id, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("address not found")
		}
		return nil, err
	}

	if req.Label != "" {
		address.Label = req.Label
	}
	if req.Recipient != "" {
		address.Recipient = req.Recipient
	}
	if req.Phone != "" {
		address.Phone = req.Phone
	}
	if req.Province != "" {
		address.Province = req.Province
	}
	if req.City != "" {
		address.City = req.City
	}
	if req.District != "" {
		address.District = req.District
	}
	if req.PostalCode != "" {
		address.PostalCode = req.PostalCode
	}
	if req.Detail != "" {
		address.Detail = req.Detail
	}
	if req.IsDefault != nil {
		if *req.IsDefault {
			if err := s.repo.ResetDefault(userID); err != nil {
				return nil, err
			}
		}
		address.IsDefault = *req.IsDefault
	}

	if err := s.repo.Update(address); err != nil {
		return nil, err
	}

	return toAddressResponse(address), nil
}

func (s *addressService) DeleteAddress(id uint, userID string) error {
	_, err := s.repo.FindByID(id, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("address not found")
		}
		return err
	}

	return s.repo.Delete(id, userID)
}

func (s *addressService) SetDefaultAddress(id uint, userID string) error {
	address, err := s.repo.FindByID(id, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("address not found")
		}
		return err
	}

	if err := s.repo.ResetDefault(userID); err != nil {
		return err
	}

	address.IsDefault = true
	return s.repo.Update(address)
}

func toAddressResponse(address *model.Address) *dto.AddressResponse {
	return &dto.AddressResponse{
		ID:         address.ID,
		UserID:     address.UserID,
		Label:      address.Label,
		Recipient:  address.Recipient,
		Phone:      address.Phone,
		Province:   address.Province,
		City:       address.City,
		District:   address.District,
		PostalCode: address.PostalCode,
		Detail:     address.Detail,
		IsDefault:  address.IsDefault,
		CreatedAt:  address.CreatedAt,
	}
}
