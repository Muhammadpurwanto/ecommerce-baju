package service

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"

	"github.com/Muhammadpurwanto/ecommerce-baju/services/user/config"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/user/internal/dto"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/user/internal/model"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/user/internal/util"
)

// Mock JWTService implementation
type mockJWTService struct {
	GenerateTokensFunc func(userID string, role string) (string, string, int64, error)
	ValidateTokenFunc  func(tokenStr string) (*jwt.Token, error)
}

func (m *mockJWTService) GenerateTokens(userID string, role string) (string, string, int64, error) {
	return m.GenerateTokensFunc(userID, role)
}
func (m *mockJWTService) ValidateToken(tokenStr string) (*jwt.Token, error) {
	return m.ValidateTokenFunc(tokenStr)
}

// Mock TokenCacheService implementation
type mockTokenCacheService struct {
	BlacklistTokenFunc     func(token string, expiration time.Duration) error
	IsTokenBlacklistedFunc func(token string) bool
}

func (m *mockTokenCacheService) BlacklistToken(token string, expiration time.Duration) error {
	return m.BlacklistTokenFunc(token, expiration)
}
func (m *mockTokenCacheService) IsTokenBlacklisted(token string) bool {
	return m.IsTokenBlacklistedFunc(token)
}

// Mock EventPublisher implementation
type mockEventPublisher struct {
	PublishUserRegisteredFunc func(uuid string, email string, name string) error
}

func (m *mockEventPublisher) PublishUserRegistered(uuid string, email string, name string) error {
	if m.PublishUserRegisteredFunc != nil {
		return m.PublishUserRegisteredFunc(uuid, email, name)
	}
	return nil
}

func TestAuthService_Register(t *testing.T) {
	cfg := &config.Config{
		JWTExpiration:        2,
		JWTRefreshExpiration: 24,
	}

	t.Run("Success Register", func(t *testing.T) {
		mockRepo := &mockUserRepository{
			FindByEmailFunc: func(email string) (*model.User, error) {
				// Email tidak ditemukan (artinya aman untuk registrasi baru)
				return nil, gorm.ErrRecordNotFound
			},
			CreateFunc: func(user *model.User) error {
				if user.Email != "new@example.com" {
					t.Errorf("diharapkan email 'new@example.com', namun mendapat: %s", user.Email)
				}
				return nil
			},
		}

		mockJWT := &mockJWTService{
			GenerateTokensFunc: func(userID string, role string) (string, string, int64, error) {
				return "access-token-123", "refresh-token-123", 3600, nil
			},
		}

		mockPublisher := &mockEventPublisher{}

		srv := NewAuthService(mockRepo, mockJWT, nil, cfg, mockPublisher)
		req := &dto.RegisterRequest{
			Email:    "new@example.com",
			Name:     "New User",
			Password: "securepassword",
		}

		resp, err := srv.Register(req)

		if err != nil {
			t.Fatalf("diharapkan tidak ada error, namun mendapat: %v", err)
		}
		if resp.AccessToken != "access-token-123" {
			t.Errorf("diharapkan AccessToken = 'access-token-123', mendapat: %s", resp.AccessToken)
		}
	})

	t.Run("Register Duplicate Email", func(t *testing.T) {
		mockRepo := &mockUserRepository{
			FindByEmailFunc: func(email string) (*model.User, error) {
				// Email sudah terdaftar
				return &model.User{Email: "existing@example.com"}, nil
			},
		}

		srv := NewAuthService(mockRepo, nil, nil, cfg, nil)
		req := &dto.RegisterRequest{
			Email:    "existing@example.com",
			Name:     "Existing User",
			Password: "password",
		}

		resp, err := srv.Register(req)

		if err == nil {
			t.Fatal("diharapkan terjadi error email duplikat, namun mendapat nil")
		}
		if err.Error() != "email already exists" {
			t.Errorf("diharapkan pesan error 'email already exists', namun mendapat '%s'", err.Error())
		}
		if resp != nil {
			t.Errorf("diharapkan response bernilai nil, namun mendapat %+v", resp)
		}
	})
}

func TestAuthService_Login(t *testing.T) {
	cfg := &config.Config{
		JWTExpiration:        2,
		JWTRefreshExpiration: 24,
	}

	hashedPassword, _ := util.HashPassword("correctpassword")

	t.Run("Success Login", func(t *testing.T) {
		mockRepo := &mockUserRepository{
			FindByEmailFunc: func(email string) (*model.User, error) {
				return &model.User{
					ID:       "user-123",
					Email:    "user@example.com",
					Password: hashedPassword,
					Role:     "customer",
				}, nil
			},
		}

		mockJWT := &mockJWTService{
			GenerateTokensFunc: func(userID string, role string) (string, string, int64, error) {
				return "access-token", "refresh-token", 3600, nil
			},
		}

		srv := NewAuthService(mockRepo, mockJWT, nil, cfg, nil)
		req := &dto.LoginRequest{
			Email:    "user@example.com",
			Password: "correctpassword",
		}

		resp, err := srv.Login(req)

		if err != nil {
			t.Fatalf("diharapkan login sukses tanpa error, namun mendapat: %v", err)
		}
		if resp.AccessToken != "access-token" {
			t.Errorf("diharapkan AccessToken = 'access-token', mendapat: %s", resp.AccessToken)
		}
	})

	t.Run("Login Incorrect Password", func(t *testing.T) {
		mockRepo := &mockUserRepository{
			FindByEmailFunc: func(email string) (*model.User, error) {
				return &model.User{
					Email:    "user@example.com",
					Password: hashedPassword,
				}, nil
			},
		}

		srv := NewAuthService(mockRepo, nil, nil, cfg, nil)
		req := &dto.LoginRequest{
			Email:    "user@example.com",
			Password: "wrongpassword",
		}

		resp, err := srv.Login(req)

		if err == nil {
			t.Fatal("diharapkan login gagal dengan password salah, mendapat nil")
		}
		if err.Error() != "invalid email or password" {
			t.Errorf("diharapkan pesan error 'invalid email or password', mendapat '%s'", err.Error())
		}
		if resp != nil {
			t.Errorf("diharapkan response bernilai nil, mendapat %+v", resp)
		}
	})
}

func TestAuthService_Logout(t *testing.T) {
	cfg := &config.Config{
		JWTExpiration:        2,
		JWTRefreshExpiration: 24,
	}

	t.Run("Success Logout", func(t *testing.T) {
		blacklistedCount := 0
		mockCache := &mockTokenCacheService{
			BlacklistTokenFunc: func(token string, expiration time.Duration) error {
				blacklistedCount++
				return nil
			},
		}

		srv := NewAuthService(nil, nil, mockCache, cfg, nil)
		err := srv.Logout("access", "refresh")

		if err != nil {
			t.Fatalf("diharapkan logout sukses, namun mendapat: %v", err)
		}
		if blacklistedCount != 2 {
			t.Errorf("diharapkan mem-blacklist 2 token (access & refresh), namun memproses: %d", blacklistedCount)
		}
	})
}

func TestAuthService_RefreshToken(t *testing.T) {
	cfg := &config.Config{
		JWTExpiration:        2,
		JWTRefreshExpiration: 24,
	}

	t.Run("Success Refresh Token", func(t *testing.T) {
		mockJWT := &mockJWTService{
			ValidateTokenFunc: func(tokenStr string) (*jwt.Token, error) {
				// Mengembalikan token valid dengan klaim refresh
				return &jwt.Token{
					Valid: true,
					Claims: jwt.MapClaims{
						"type":    "refresh",
						"user_id": "user-123",
						"role":    "customer",
					},
				}, nil
			},
			GenerateTokensFunc: func(userID string, role string) (string, string, int64, error) {
				return "new-access", "new-refresh", 3600, nil
			},
		}

		mockCache := &mockTokenCacheService{
			IsTokenBlacklistedFunc: func(token string) bool {
				return false
			},
			BlacklistTokenFunc: func(token string, expiration time.Duration) error {
				return nil
			},
		}

		srv := NewAuthService(nil, mockJWT, mockCache, cfg, nil)
		req := &dto.RefreshTokenRequest{
			RefreshToken: "old-refresh-token",
		}

		resp, err := srv.RefreshToken(req)

		if err != nil {
			t.Fatalf("diharapkan sukses melakukan refresh token, mendapat error: %v", err)
		}
		if resp.AccessToken != "new-access" {
			t.Errorf("diharapkan access token baru 'new-access', mendapat: %s", resp.AccessToken)
		}
	})

	t.Run("Refresh Blacklisted Token", func(t *testing.T) {
		mockJWT := &mockJWTService{
			ValidateTokenFunc: func(tokenStr string) (*jwt.Token, error) {
				return &jwt.Token{
					Valid: true,
					Claims: jwt.MapClaims{
						"type":    "refresh",
						"user_id": "user-123",
						"role":    "customer",
					},
				}, nil
			},
		}

		mockCache := &mockTokenCacheService{
			IsTokenBlacklistedFunc: func(token string) bool {
				return true // Token sudah di-blacklist
			},
		}

		srv := NewAuthService(nil, mockJWT, mockCache, cfg, nil)
		req := &dto.RefreshTokenRequest{
			RefreshToken: "blacklisted-token",
		}

		resp, err := srv.RefreshToken(req)

		if err == nil {
			t.Fatal("diharapkan gagal karena token diblacklist, mendapat nil")
		}
		if err.Error() != "refresh token is blacklisted" {
			t.Errorf("diharapkan error 'refresh token is blacklisted', mendapat '%s'", err.Error())
		}
		if resp != nil {
			t.Errorf("diharapkan response bernilai nil, mendapat %+v", resp)
		}
	})
}
