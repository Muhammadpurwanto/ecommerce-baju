package service

import (
	"testing"

	"gorm.io/gorm"

	"github.com/Muhammadpurwanto/ecommerce-baju/services/user/internal/dto"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/user/internal/model"
)

// Mock UserRepository implementation
type mockUserRepository struct {
	FindByIDFunc    func(id string) (*model.User, error)
	FindByEmailFunc func(email string) (*model.User, error)
	FindAllFunc     func(page, perPage int) ([]model.User, int64, error)
	UpdateFunc      func(user *model.User) error
	CreateFunc      func(user *model.User) error
}

func (m *mockUserRepository) FindByID(id string) (*model.User, error) {
	return m.FindByIDFunc(id)
}
func (m *mockUserRepository) FindByEmail(email string) (*model.User, error) {
	return m.FindByEmailFunc(email)
}
func (m *mockUserRepository) FindAll(page, perPage int) ([]model.User, int64, error) {
	return m.FindAllFunc(page, perPage)
}
func (m *mockUserRepository) Update(user *model.User) error {
	return m.UpdateFunc(user)
}
func (m *mockUserRepository) Create(user *model.User) error {
	return m.CreateFunc(user)
}

func TestUserService_GetProfile(t *testing.T) {
	t.Run("Success Get Profile", func(t *testing.T) {
		mockRepo := &mockUserRepository{
			FindByIDFunc: func(id string) (*model.User, error) {
				return &model.User{
					ID:       "user-1",
					Email:    "test@example.com",
					Name:     "Test User",
					Role:     "customer",
					IsActive: true,
				}, nil
			},
		}

		srv := NewUserService(mockRepo)
		resp, err := srv.GetProfile("user-1")

		if err != nil {
			t.Fatalf("diharapkan tidak ada error, namun mendapat: %v", err)
		}
		if resp.ID != "user-1" {
			t.Errorf("diharapkan ID 'user-1', namun mendapat: %s", resp.ID)
		}
		if resp.Email != "test@example.com" {
			t.Errorf("diharapkan Email 'test@example.com', namun mendapat: %s", resp.Email)
		}
	})

	t.Run("User Not Found", func(t *testing.T) {
		mockRepo := &mockUserRepository{
			FindByIDFunc: func(id string) (*model.User, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		srv := NewUserService(mockRepo)
		resp, err := srv.GetProfile("non-existent")

		if err == nil {
			t.Fatal("diharapkan terjadi error, namun mendapat nil")
		}
		if err.Error() != "user not found" {
			t.Errorf("diharapkan error message 'user not found', namun mendapat: '%s'", err.Error())
		}
		if resp != nil {
			t.Errorf("diharapkan response bernilai nil, namun mendapat: %+v", resp)
		}
	})
}

func TestUserService_UpdateProfile(t *testing.T) {
	t.Run("Success Update Profile", func(t *testing.T) {
		updatedName := "New Name"
		updatedPhone := "0812345678"
		updatedAvatar := "http://avatar.url"

		mockRepo := &mockUserRepository{
			FindByIDFunc: func(id string) (*model.User, error) {
				return &model.User{
					ID:    "user-1",
					Email: "test@example.com",
					Name:  "Old Name",
				}, nil
			},
			UpdateFunc: func(user *model.User) error {
				if user.Name != updatedName {
					t.Errorf("diharapkan Name di-update menjadi '%s', namun mendapat '%s'", updatedName, user.Name)
				}
				return nil
			},
		}

		srv := NewUserService(mockRepo)
		req := &dto.UpdateUserRequest{
			Name:      updatedName,
			Phone:     &updatedPhone,
			AvatarURL: &updatedAvatar,
		}

		resp, err := srv.UpdateProfile("user-1", req)

		if err != nil {
			t.Fatalf("diharapkan tidak ada error, namun mendapat: %v", err)
		}
		if resp.Name != updatedName {
			t.Errorf("diharapkan Name di respon bernilai '%s', namun mendapat '%s'", updatedName, resp.Name)
		}
	})

	t.Run("User Not Found On Update", func(t *testing.T) {
		mockRepo := &mockUserRepository{
			FindByIDFunc: func(id string) (*model.User, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}

		srv := NewUserService(mockRepo)
		resp, err := srv.UpdateProfile("non-existent", &dto.UpdateUserRequest{Name: "New Name"})

		if err == nil {
			t.Fatal("diharapkan terjadi error, namun mendapat nil")
		}
		if err.Error() != "user not found" {
			t.Errorf("diharapkan error 'user not found', namun mendapat '%s'", err.Error())
		}
		if resp != nil {
			t.Errorf("diharapkan response bernilai nil, namun mendapat %+v", resp)
		}
	})
}

func TestUserService_GetAllUsers(t *testing.T) {
	t.Run("Success Get All Users", func(t *testing.T) {
		mockRepo := &mockUserRepository{
			FindAllFunc: func(page, perPage int) ([]model.User, int64, error) {
				return []model.User{
					{ID: "user-1", Name: "User One"},
					{ID: "user-2", Name: "User Two"},
				}, 2, nil
			},
		}

		srv := NewUserService(mockRepo)
		resp, meta, err := srv.GetAllUsers(1, 10)

		if err != nil {
			t.Fatalf("diharapkan tidak ada error, namun mendapat: %v", err)
		}
		if len(resp) != 2 {
			t.Errorf("diharapkan mendapat 2 user, namun mendapat: %d", len(resp))
		}
		if meta.Total != 2 {
			t.Errorf("diharapkan total data bernilai 2, namun mendapat: %d", meta.Total)
		}
		if meta.TotalPages != 1 {
			t.Errorf("diharapkan total halaman bernilai 1, namun mendapat: %d", meta.TotalPages)
		}
	})
}
