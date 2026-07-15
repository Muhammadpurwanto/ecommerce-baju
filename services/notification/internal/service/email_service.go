package service

import (
	"go.uber.org/zap"
	"gopkg.in/gomail.v2"

	"github.com/Muhammadpurwanto/ecommerce-baju/services/notification/config"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/notification/internal/dto"
)

type EmailService interface {
	SendEmail(req *dto.SendEmailRequest) error
}

type emailService struct {
	cfg *config.Config
	log *zap.Logger
}

func NewEmailService(cfg *config.Config, log *zap.Logger) EmailService {
	return &emailService{cfg: cfg, log: log}
}

func (s *emailService) SendEmail(req *dto.SendEmailRequest) error {
	m := gomail.NewMessage()
	
	// Format pengirim: "Nama <email@domain.com>"
	m.SetHeader("From", m.FormatAddress(s.cfg.SMTPUser, s.cfg.SenderName))
	m.SetHeader("To", req.To)
	m.SetHeader("Subject", req.Subject)
	m.SetBody("text/html", req.Body)

	d := gomail.NewDialer(s.cfg.SMTPHost, s.cfg.SMTPPort, s.cfg.SMTPUser, s.cfg.SMTPPass) // Menginisialisasi dialer untuk mengirim email

	if err := d.DialAndSend(m); err != nil {
		s.log.Error("Failed to send email", zap.Error(err), zap.String("to", req.To))
		return err
	}

	s.log.Info("Email sent successfully", zap.String("to", req.To))
	return nil
}
