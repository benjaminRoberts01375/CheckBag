package email

func SendForgotPassword(email string, activationToken string) {
	message := `A password reset was issued for your` + Config.ServiceName + ` account. To reset your password, click the link below:
https://giftguardian.benlab.us/db/reset-password/` + activationToken + `
If you did not request a password reset, someone else may have attempted to reset your password.
If you did try to reset your password and no longer wish to, simply ignore this email. The link will expire shortly.`
	sendEmail(email, Config.ServiceName+" Password Reset", message)
}
