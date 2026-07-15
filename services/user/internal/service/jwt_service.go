package service

import (
	"fmt"
	"time"

	"github.com/Muhammadpurwanto/ecommerce-baju/services/user/config"
	"github.com/golang-jwt/jwt/v5"
)

type JWTService interface {
	GenerateTokens(userID string, role string) (string, string, int64, error)
	ValidateToken(tokenStr string) (*jwt.Token, error)
}

type jwtService struct {
	cfg *config.Config
}

func NewJWTService(cfg *config.Config) JWTService {
	return &jwtService{cfg: cfg}
}

func (s *jwtService) GenerateTokens(userID string, role string) (string, string, int64, error) {
	expirationTime := time.Now().Add(time.Duration(s.cfg.JWTExpiration) * time.Hour)
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"type":    "access",
		"exp":     expirationTime.Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return "", "", 0, err
	}

	refreshExpirationTime := time.Now().Add(time.Duration(s.cfg.JWTRefreshExpiration) * time.Hour)
	refreshClaims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"type":    "refresh",
		"exp":     refreshExpirationTime.Unix(),
		"iat":     time.Now().Unix(),
	}

	refreshTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err := refreshTokenObj.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return "", "", 0, err
	}

	return accessToken, refreshToken, expirationTime.Unix(), nil
}

func (s *jwtService) ValidateToken(tokenStr string) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.cfg.JWTSecret), nil
	})
}
