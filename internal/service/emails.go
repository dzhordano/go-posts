package service

import (
	"fmt"

	"github.com/dzhordano/go-posts/pkg/email"
)

type EmailsService struct {
	sender email.Sender
}

func NewEmailsService(sender email.Sender) *EmailsService {
	return &EmailsService{
		sender,
	}
}

func (s *EmailsService) SendUserVerificationEmail(input VerificationEmailInput) error {
	subject := input.Name
	link := fmt.Sprintf("http://localhost:8081/api/v1/users/verify/%s", input.VerificationCode)

	body := fmt.Sprintf("Hello, %s. <a href='%s'>Click to verify your account.</a>", input.Name, link)

	sendInput := email.SendEmailInput{
		To:      input.Email,
		Subject: subject,
		Body:    body,
	}

	return s.sender.Send(sendInput)
}
