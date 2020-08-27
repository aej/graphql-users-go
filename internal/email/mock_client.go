package email

import (
	"fmt"
	"log"
)

type MockEmailClient struct{}

func (m MockEmailClient) SendConfirmationEmail(email string, confirmationUrl string) {
	body := fmt.Sprintf("Hi,\n\nPlease verify your account by following the link below:\n%s", confirmationUrl)
	log.Print(body)
}

func (m MockEmailClient) SendPasswordResetEmail(email string, passwordResetToken string) {
	body := fmt.Sprintf("Hi,\n\nA password reset has been requested on your account. Use the following token\n%s", passwordResetToken)
	log.Print(body)
}
