package email

type EmailClient interface {
	SendConfirmationEmail(email string, confirmationUrl string)
	SendPasswordResetEmail(email string, passwordResetToken string)
}

var emailClient EmailClient

func init() {
	emailClient = MockEmailClient{}
}

func SendConfirmationEmail(email string, confirmationUrl string) {
	emailClient.SendConfirmationEmail(email, confirmationUrl)
}

func SendPasswordResetEmail(email string, passwordResetToken string) {
	emailClient.SendPasswordResetEmail(email, passwordResetToken)
}
