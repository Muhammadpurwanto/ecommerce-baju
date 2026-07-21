package service

import (
	"testing"

	"go.uber.org/zap"

	"github.com/Muhammadpurwanto/ecommerce-baju/services/notification/config"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/notification/internal/dto"
)

func TestEmailService_SendEmail_FailurePath(t *testing.T) {
	// Dummy config dengan SMTP Host tidak valid agar gagal melakukan koneksi dial
	cfg := &config.Config{
		SMTPHost:   "127.0.0.1",
		SMTPPort:   9999, // Port yang pasti tidak aktif/tertutup
		SMTPUser:   "test@example.com",
		SMTPPass:   "password",
		SenderName: "Test Sender",
	}

	logger, _ := zap.NewDevelopment()
	srv := NewEmailService(cfg, logger)

	req := &dto.SendEmailRequest{
		To:      "recipient@example.com",
		Subject: "Test Subject",
		Body:    "<h1>Hello</h1>",
	}

	err := srv.SendEmail(req)

	// Diharapkan terjadi error koneksi karena host/port tidak valid
	if err == nil {
		t.Error("diharapkan terjadi error koneksi dial SMTP, namun mendapat nil")
	}
}
