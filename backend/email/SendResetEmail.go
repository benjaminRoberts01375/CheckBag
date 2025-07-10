package email

func SendResetEmail(updatedEmail string, activationToken string) {
	emailMessage := `Please verify your new ` + Config.ServiceName + ` email address by clicking the link below:
` + Config.URL() + `user-reset-email-confirmation/` + activationToken
	sendEmail(updatedEmail, Config.ServiceName+" Email Reset", emailMessage)
}
