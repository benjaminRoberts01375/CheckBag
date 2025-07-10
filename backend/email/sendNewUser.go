package email

func SendNewUser(email string, firstName string, activationToken string) {
	message := `Hello ` + firstName + `, and welcome to` + Config.ServiceName + `! Please click the link below to confirm your account:
` + Config.URL() + `user-sign-up-confirmation/` + activationToken
	sendEmail(email, "Gift Guardian Account Confirmation", message)
}
