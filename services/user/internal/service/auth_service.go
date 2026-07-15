package service

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"

	"github.com/Muhammadpurwanto/ecommerce-baju/services/user/broker"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/user/config"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/user/internal/dto"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/user/internal/model"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/user/internal/repository"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/user/internal/util"
)

type AuthService interface {
	Register(req *dto.RegisterRequest) (*dto.TokenResponse, error)
	Login(req *dto.LoginRequest) (*dto.TokenResponse, error)
	Logout(accessToken string, refreshToken string) error
	RefreshToken(req *dto.RefreshTokenRequest) (*dto.TokenResponse, error)
	LoginOrRegisterOAuth(userInfo *dto.GoogleUserInfo) (*dto.TokenResponse, error)
}

type authService struct {
	repo       repository.UserRepository
	jwtSvc     JWTService
	cfg        *config.Config
	tokenCache TokenCacheService
	eventPub   broker.EventPublisher
}

func NewAuthService(repo repository.UserRepository, jwtSvc JWTService, tokenCache TokenCacheService, cfg *config.Config, eventPub broker.EventPublisher) AuthService {
	return &authService{
		repo:       repo,
		jwtSvc:     jwtSvc,
		tokenCache: tokenCache,
		cfg:        cfg,
		eventPub:   eventPub,
	}
}

func (s *authService) Register(req *dto.RegisterRequest) (*dto.TokenResponse, error) {
	_, err := s.repo.FindByEmail(req.Email)
	if err == nil {
		return nil, errors.New("email already exists")
	}

	generatedUserID := uuid.New().String()

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		ID:       generatedUserID,
		Email:    req.Email,
		Name:     req.Name,
		Password: hashedPassword,
		Provider: "local",
		Role:     "customer",
		IsActive: true,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	if s.eventPub != nil {
		go s.eventPub.PublishUserRegistered(generatedUserID, req.Email, req.Name)
	}

	accessToken, refreshToken, exp, err := s.jwtSvc.GenerateTokens(user.ID, "customer")
	if err != nil {
		return nil, err
	}

	return &dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    exp,
	}, nil
}

func (s *authService) Login(req *dto.LoginRequest) (*dto.TokenResponse, error) {
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid email or password")
		}
		return nil, err
	}

	if !util.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("invalid email or password")
	}

	role := user.Role
	if role == "" {
		role = "customer"
	}

	accessToken, refreshToken, exp, err := s.jwtSvc.GenerateTokens(user.ID, role)
	if err != nil {
		return nil, err
	}

	return &dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    exp,
	}, nil
}

func (s *authService) Logout(accessToken string, refreshToken string) error {
	expiration := time.Duration(s.cfg.JWTExpiration) * time.Hour
	if err := s.tokenCache.BlacklistToken(accessToken, expiration); err != nil {
		return err
	}

	refreshExpiration := time.Duration(s.cfg.JWTRefreshExpiration) * time.Hour
	return s.tokenCache.BlacklistToken(refreshToken, refreshExpiration)
}

func (s *authService) RefreshToken(req *dto.RefreshTokenRequest) (*dto.TokenResponse, error) {
	token, err := s.jwtSvc.ValidateToken(req.RefreshToken)
	if err != nil || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid refresh token claims")
	}

	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "refresh" {
		return nil, errors.New("invalid token type")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return nil, errors.New("invalid user id in token")
	}
	role, ok := claims["role"].(string)
	if !ok {
		return nil, errors.New("invalid role in token")
	}

	if s.tokenCache.IsTokenBlacklisted(req.RefreshToken) {
		return nil, errors.New("refresh token is blacklisted")
	}

	accessToken, newRefreshToken, exp, err := s.jwtSvc.GenerateTokens(userID, role)
	if err != nil {
		return nil, err
	}

	s.tokenCache.BlacklistToken(req.RefreshToken, time.Duration(s.cfg.JWTRefreshExpiration)*time.Hour)

	return &dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    exp,
	}, nil
}

func (s *authService) LoginOrRegisterOAuth(userInfo *dto.GoogleUserInfo) (*dto.TokenResponse, error) {
	user, err := s.repo.FindByEmail(userInfo.Email)
	if err == nil {
		if user.ProviderID == nil || *user.ProviderID == "" {
			providerID := userInfo.ID
			user.ProviderID = &providerID
			user.Provider = "google"
			if userInfo.Picture != "" {
				user.AvatarURL = &userInfo.Picture
			}
			s.repo.Update(user)
		}
	} else {
		generatedUserID := uuid.New().String()
		providerID := userInfo.ID
		avatarURL := userInfo.Picture

		user = &model.User{
			ID:         generatedUserID,
			Email:      userInfo.Email,
			Name:       userInfo.Name,
			Provider:   "google",
			ProviderID: &providerID,
			Role:       "customer",
			IsActive:   true,
		}
		if avatarURL != "" {
			user.AvatarURL = &avatarURL
		}

		if err := s.repo.Create(user); err != nil {
			return nil, err
		}

		if s.eventPub != nil {
			go s.eventPub.PublishUserRegistered(generatedUserID, userInfo.Email, userInfo.Name)
		}
	}

	role := user.Role
	if role == "" {
		role = "customer"
	}

	accessToken, refreshToken, exp, err := s.jwtSvc.GenerateTokens(user.ID, role)
	if err != nil {
		return nil, err
	}

	return &dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    exp,
	}, nil
}
