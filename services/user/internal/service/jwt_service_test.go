package service

import (
	"testing"
	"time"

	"github.com/Muhammadpurwanto/ecommerce-baju/services/user/config"
	"github.com/golang-jwt/jwt/v5"
)

func TestJWTService(t *testing.T) {
	// Inisialisasi dummy konfigurasi
	cfg := &config.Config{
		JWTSecret:            "rahasia-negara-paling-aman-yang-sangat-panjang-sekali",
		JWTExpiration:        2,
		JWTRefreshExpiration: 24,
	}

	jwtSrv := NewJWTService(cfg)

	var accessToken, refreshToken string
	var exp int64
	var err error

	t.Run("Generate Tokens Successfully", func(t *testing.T) {
		userID := "user-123"
		role := "customer"

		accessToken, refreshToken, exp, err = jwtSrv.GenerateTokens(userID, role)
		if err != nil {
			t.Fatalf("Gagal membuat token: %v", err)
		}

		if accessToken == "" {
			t.Error("Diharapkan access token tidak kosong")
		}

		if refreshToken == "" {
			t.Error("Diharapkan refresh token tidak kosong")
		}

		// Waktu kedaluwarsa harus di masa depan
		if exp <= time.Now().Unix() {
			t.Errorf("Diharapkan expiration %d berada di masa depan", exp)
		}
	})

	t.Run("Validate Correct Access Token", func(t *testing.T) {
		token, err := jwtSrv.ValidateToken(accessToken)
		if err != nil {
			t.Fatalf("Gagal validasi access token: %v", err)
		}

		if !token.Valid {
			t.Error("Diharapkan token bernilai valid")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			t.Fatal("Diharapkan claims bertipe jwt.MapClaims")
		}

		if claims["user_id"] != "user-123" {
			t.Errorf("Diharapkan user_id = 'user-123', namun mendapat: %v", claims["user_id"])
		}

		if claims["role"] != "customer" {
			t.Errorf("Diharapkan role = 'customer', namun mendapat: %v", claims["role"])
		}

		if claims["type"] != "access" {
			t.Errorf("Diharapkan type = 'access', namun mendapat: %v", claims["type"])
		}
	})

	t.Run("Validate Correct Refresh Token", func(t *testing.T) {
		token, err := jwtSrv.ValidateToken(refreshToken)
		if err != nil {
			t.Fatalf("Gagal validasi refresh token: %v", err)
		}

		if !token.Valid {
			t.Error("Diharapkan refresh token bernilai valid")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			t.Fatal("Diharapkan claims bertipe jwt.MapClaims")
		}

		if claims["type"] != "refresh" {
			t.Errorf("Diharapkan type = 'refresh', namun mendapat: %v", claims["type"])
		}
	})

	t.Run("Reject Invalid Token Signature", func(t *testing.T) {
		invalidToken := accessToken + "corrupted-signature"
		_, err := jwtSrv.ValidateToken(invalidToken)
		if err == nil {
			t.Error("Diharapkan terjadi error saat memvalidasi token yang rusak, namun mendapat nil")
		}
	})
}
