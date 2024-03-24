package smtp

import (
	"errors"
	"log/slog"

	emailverifier "github.com/AfterShip/email-verifier"
	"github.com/dzhordano/go-posts/pkg/email"
	gomail "gopkg.in/gomail.v2"
)

type SMTPSender struct {
	from string
	pass string
	host string
	port int
}

func NewSMTPSender(from, pass, host string, port int) (*SMTPSender, error) {
	if !emailverifier.IsAddressValid(from) {
		return nil, errors.New("invalid from email")
	}

	return &SMTPSender{
		from,
		pass,
		host,
		port,
	}, nil
}

func (s *SMTPSender) Send(input email.SendEmailInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	// forming message
	msg := gomail.NewMessage()
	msg.SetHeader("From", s.from)
	msg.SetHeader("To", input.To)
	msg.SetHeader("Subject", input.Subject)
	// TODO: нужно ли прям верифицировать почту? или можно просто проверить валидность почты при регистрации?
	msg.SetBody("text/html", input.Body)

	dialer := gomail.NewDialer(s.host, s.port, s.from, s.pass)
	if err := dialer.DialAndSend(msg); err != nil {
		slog.Info(err.Error())
		return errors.New("failed to send email via smtp")
	}

	return nil
}
