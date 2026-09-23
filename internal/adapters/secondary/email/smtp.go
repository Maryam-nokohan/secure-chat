package email

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"net/smtp"
	"strings"
	"time"

	"github.com/maryam-nokohan/secure-chat/internal/core/ports"
	"github.com/maryam-nokohan/secure-chat/pkg"
)

type SMTPConfig struct {
	Host, Port, Username, Password, From string
}

type SMTPSender struct {
	cfg  SMTPConfig
	from *mail.Address
}

func NewSMTPSender(cfg SMTPConfig) (ports.EmailSender, error) {
	from, err := mail.ParseAddress(cfg.From)
	if err != nil {
		return nil, fmt.Errorf("invalid SMTP_FROM: %w", err)
	}
	if cfg.Port == "" {
		cfg.Port = "587"
	}
	return &SMTPSender{cfg: cfg, from: from}, nil
}

func (s *SMTPSender) SendVerificationCode(ctx context.Context, to, code string, ttl time.Duration) error {
	if strings.ContainsAny(to, "\r\n") {
		return errors.New("invalid recipient")
	}
	body := fmt.Sprintf("Your Hovar verification code is: %s\r\n\r\nIt expires in %d minutes. If you didn't request it, ignore this email.\r\n",
		code, int(ttl.Minutes()))
	msg := strings.Join([]string{
		"From: " + s.from.String(),
		"To: " + to,
		"Subject: Your Hovar verification code",
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=UTF-8",
		"",
		body,
	}, "\r\n")

	var auth smtp.Auth
	if s.cfg.Username != "" {
		auth = smtp.PlainAuth("", s.cfg.Username, s.cfg.Password, s.cfg.Host)
	}

	done := make(chan error, 1)
	go func() {
		done <- smtp.SendMail(s.cfg.Host+":"+s.cfg.Port, auth, s.from.Address, []string{to}, []byte(msg))
	}()
	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(15 * time.Second):
		return errors.New("smtp: timed out")
	}
}

type ConsoleSender struct{}

func NewConsoleSender() ports.EmailSender { return ConsoleSender{} }

func (ConsoleSender) SendVerificationCode(_ context.Context, to, code string, _ time.Duration) error {
	pkg.LogInfo(fmt.Sprintf("[DEV EMAIL] verification code for %s: %s", to, code))
	return nil
}
