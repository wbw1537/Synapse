package notification

import (
	"fmt"
	"log"
	"net/smtp"
	"strings"

	"github.com/wbw1537/synapse/internal/config"
)

type Sender interface {
	Send(subject, body string) error
}

type SMTPSender struct {
	cfg *config.Config
}

func NewSMTPSender(cfg *config.Config) *SMTPSender {
	return &SMTPSender{cfg: cfg}
}

func (s *SMTPSender) Send(subject, body string) error {
	if !s.cfg.EnableAlerts {
		return nil
	}
	if s.cfg.SMTPHost == "" || s.cfg.SMTPTo == "" {
		log.Println("SMTP not configured, skipping alert")
		return nil
	}

	auth := smtp.PlainAuth("", s.cfg.SMTPUser, s.cfg.SMTPPass, s.cfg.SMTPHost)
	to := strings.Split(s.cfg.SMTPTo, ",")

	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: [Synapse] %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/plain; charset=\"utf-8\"\r\n"+
		"\r\n"+
		"%s\r\n", s.cfg.SMTPTo, subject, body))

	addr := fmt.Sprintf("%s:%s", s.cfg.SMTPHost, s.cfg.SMTPPort)

	// Note: For MVP we assume STARTTLS or Plain based on port, but net/smtp defaults are tricky.
	// Often SendMail handles it. If using 587, it usually expects STARTTLS.
	err := smtp.SendMail(addr, auth, s.cfg.SMTPFrom, to, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Printf("Alert email sent to %s: %s", s.cfg.SMTPTo, subject)
	return nil
}
