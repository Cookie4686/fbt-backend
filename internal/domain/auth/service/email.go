package service

import (
	"fmt"

	"github.com/wneessen/go-mail"
)

func (s *service) SendVerificationMail(email string, otp string) error {
	message := mail.NewMsg()
	if err := message.From(s.CFG.MAIL.ADDRESS); err != nil {
		return err
	} else if err := message.To(email); err != nil {
		return err
	}
	message.Subject("Email Verification Code")
	message.SetBodyString(mail.TypeTextPlain, fmt.Sprintf("Verification Code: %s\n The code will expires in 1 hour", otp))
	if err := s.Mail.DialAndSend(message); err != nil {
		return err
	}

	return nil
}
