package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Muhammadpurwanto/ecommerce-baju/services/user/config"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/user/internal/dto"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/user/internal/service"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type OAuthHandler struct {
	cfg         *config.Config
	authService service.AuthService
	oauthConfig *oauth2.Config
}

func NewOAuthHandler(cfg *config.Config, authService service.AuthService) *OAuthHandler {
	oauthConfig := &oauth2.Config{
		ClientID:     cfg.GoogleClientID,
		ClientSecret: cfg.GoogleClientSecret,
		RedirectURL:  cfg.GoogleRedirectURL,
		Scopes:       []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint:     google.Endpoint,
	}

	return &OAuthHandler{
		cfg:         cfg,
		authService: authService,
		oauthConfig: oauthConfig,
	}
}

func (h *OAuthHandler) HandleLogin(c *fiber.Ctx) error {
	url := h.oauthConfig.AuthCodeURL("state-token")
	return c.Redirect(url)
}

func (h *OAuthHandler) HandleCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	if code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Code is required",
		})
	}

	token, err := h.oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to exchange token: " + err.Error(),
		})
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to get user info: " + err.Error(),
		})
	}
	defer resp.Body.Close()

	var userInfo dto.GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse user info: " + err.Error(),
		})
	}

	tokens, err := h.authService.LoginOrRegisterOAuth(&userInfo)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to authenticate: " + err.Error(),
		})
	}

	redirectURL := fmt.Sprintf("%s?access_token=%s&refresh_token=%s",
		h.cfg.FrontendCallbackURL,
		url.QueryEscape(tokens.AccessToken),
		url.QueryEscape(tokens.RefreshToken),
	)

	return c.Redirect(redirectURL)
}
