package mailresend

import (
	"context"
	"fmt"
	"html"
	"strings"

	"github.com/resend/resend-go/v3"
)

// Sender sends OTP messages via Resend when APIKey is set; otherwise Send* is a no-op.
type Sender struct {
	cfg    Config
	client *resend.Client
}

func NewSender(cfg Config) *Sender {
	key := strings.TrimSpace(cfg.APIKey)
	if key == "" {
		return &Sender{cfg: cfg}
	}
	return &Sender{
		cfg:    cfg,
		client: resend.NewClient(key),
	}
}

func (s *Sender) SendRegisterOTP(ctx context.Context, toEmail, code string) error {
	return s.send(ctx, toEmail, code, "Код подтверждения регистрации")
}

func (s *Sender) SendLoginOTP(ctx context.Context, toEmail, code string) error {
	return s.send(ctx, toEmail, code, "Код подтверждения входа")
}

func (s *Sender) send(ctx context.Context, toEmail, code, subject string) error {
	if s.client == nil {
		return nil
	}
	from := strings.TrimSpace(s.cfg.From)
	if from == "" {
		return fmt.Errorf("RESEND_FROM is required when RESEND_API_KEY is set")
	}
	to := strings.TrimSpace(toEmail)
	if to == "" {
		return fmt.Errorf("empty recipient email")
	}

	safe := html.EscapeString(strings.TrimSpace(code))
	body := `<!DOCTYPE html><html lang="ru"><head><meta charset="UTF-8"></head><body style="margin:0;padding:24px;font-family:system-ui,sans-serif;">` +
		`<p style="font-size:28px;letter-spacing:0.25em;font-weight:600;margin:0;">` + safe + `</p></body></html>`

	params := &resend.SendEmailRequest{
		From:    from,
		To:      []string{to},
		Subject: subject,
		Html:    body,
	}

	_, err := s.client.Emails.SendWithContext(ctx, params)
	return err
}
