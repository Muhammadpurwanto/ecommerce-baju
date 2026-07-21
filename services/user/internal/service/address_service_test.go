package service

import (
	"testing"

	"gorm.io/gorm"

	"github.com/Muhammadpurwanto/ecommerce-baju/services/user/internal/dto"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/user/internal/model"
)

// Mock AddressRepository implementation
type mockAddressRepository struct {
	FindByUserIDFunc func(userID string) ([]model.Address, error)
	FindByIDFunc     func(id uint, userID string) (*model.Address, error)
	CreateFunc       func(address *model.Address) error
	UpdateFunc       func(address *model.Address) error
	DeleteFunc       func(id uint, userID string) error
	ResetDefaultFunc func(userID string) error
}

func (m *mockAddressRepository) FindByUserID(userID string) ([]model.Address, error) {
	return m.FindByUserIDFunc(userID)
}
func (m *mockAddressRepository) FindByID(id uint, userID string) (*model.Address, error) {
	return m.FindByIDFunc(id, userID)
}
func (m *mockAddressRepository) Create(address *model.Address) error {
	return m.CreateFunc(address)
}
func (m *mockAddressRepository) Update(address *model.Address) error {
	return m.UpdateFunc(address)
}
func (m *mockAddressRepository) Delete(id uint, userID string) error {
	return m.DeleteFunc(id, userID)
}
func (m *mockAddressRepository) ResetDefault(userID string) error {
	return m.ResetDefaultFunc(userID)
}

func TestAddressService_GetAddresses(t *testing.T) {
	t.Run("Success Get Addresses", func(t *testing.T) {
		mockRepo := &mockAddressRepository{
			FindByUserIDFunc: func(userID string) ([]model.Address, error) {
				return []model.Address{
					{ID: 1, UserID: "user-1", Label: "Home"},
					{ID: 2, UserID: "user-1", Label: "Office"},
				}, nil
			},
		}

		srv := NewAddressService(mockRepo)
		resp, err := srv.GetAddresses("user-1")

		if err != nil {
			t.Fatalf("diharapkan tidak ada error, namun mendapat: %v", err)
		}
		if len(resp) != 2 {
			t.Errorf("diharapkan mendapat 2 alamat, namun mendapat: %d", len(resp))
		}
	})
}

func TestAddressService_CreateAddress(t *testing.T) {
	t.Run("Success Create Default Address", func(t *testing.T) {
		resetCalled := false
		mockRepo := &mockAddressRepository{
			ResetDefaultFunc: func(userID string) error {
				resetCalled = true
				return nil
			},
			CreateFunc: func(address *model.Address) error {
				return nil
			},
		}

		srv := NewAddressService(mockRepo)
		req := &dto.CreateAddressRequest{
			Label:     "Home",
			Recipient: "John Doe",
			IsDefault: true,
		}

		resp, err := srv.CreateAddress("user-1", req)

		if err != nil {
			t.Fatalf("diharapkan tidak ada error, mendapat: %v", err)
		}
		if !resetCalled {
			t.Error("diharapkan ResetDefault terpanggil karena alamat diset default")
		}
		if resp.Label != "Home" {
			t.Errorf("diharapkan label 'Home', mendapat: %s", resp.Label)
		}
	})
}

func TestAddressService_GetAddressByID(t *testing.T) {
	t.Run("Success Get Address By ID", func(t *testing.T) {
		mockRepo := &mockAddressRepository{
			FindByIDFunc: func(id uint, userID string) (*model.Address, error) {
				return &model.Address{ID: 1, UserID: "user-1", Label: "Home"}, nil
			},
		}

		srv := NewAddressService(mockRepo)
		resp, err := srv.GetAddressByID(1, "user-1")

		if err != nil {
			t.Fatalf("diharapkan tidak ada error, mendapat: %v", err)
		}
		if resp.ID != 1 {
			t.Errorf("diharapkan ID 1, mendapat: %d", resp.ID)
		}
	})

	t.Run("Address Not Found", func(t *testing.T) {
		mockRepo := &mockAddressRepository{
			FindByIDFunc: func(id uint, userID string) (*model.Address, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		srv := NewAddressService(mockRepo)
		resp, err := srv.GetAddressByID(99, "user-1")

		if err == nil {
			t.Fatal("diharapkan terjadi error, mendapat nil")
		}
		if err.Error() != "address not found" {
			t.Errorf("diharapkan error 'address not found', mendapat: %s", err.Error())
		}
		if resp != nil {
			t.Errorf("diharapkan response nil, mendapat: %+v", resp)
		}
	})
}
