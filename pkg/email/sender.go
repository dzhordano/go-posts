package email

import (
	"errors"

	emailverifier "github.com/AfterShip/email-verifier"
)

type SendEmailInput struct {
	To      string
	Subject string
	Body    string
}

type Sender interface {
	Send(input SendEmailInput) error
}

func (i *SendEmailInput) Validate() error {
	if i.To == "" {
		return errors.New("empty To")
	}

	if i.Subject == "" || i.Body == "" {
		return errors.New("empty subject / body")
	}

	if !emailverifier.IsAddressValid(i.To) {
		return errors.New("invalid To address")
	}

	return nil
}
