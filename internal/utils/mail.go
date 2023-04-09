package utils

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/thegeorgenikhil/armur-backend-assignment/pkg/mail"
)

var (
	fromEmail         string
	fromEmailPassword string
	frontendURL       string
)

// SendUserVerificationEmail is a utility function for the original SendMailWithHTML function, which sends the verification email
func SendUserVerificationEmail(to string, token string) error {
	verificationURL := fmt.Sprintf("%s/api/user-activation/verify?token=%s", frontendURL, token)
	err := mail.SendMailWithHTML(
		"Verify your email address",
		fmt.Sprintf(`<p>Please verify your email by clicking on the button below</p><a style="display: inline-block; background-color: #136dfa; color: white; padding: 10px 20px; text-align: center; text-decoration: none; border-radius: 5px;" href="%s">Verify Email</a>`, verificationURL),
		fromEmail,
		[]string{to},
		fromEmailPassword,
	)
	return err
}

// SendPasswordResetEmail is a utility function for the original SendMailWithHTML function, which sends the password reset email
func SendPasswordResetEmail(to string, token string) error {
	resetURL := fmt.Sprintf("%s/reset-password.html?token=%s&email=%s", frontendURL, token, to)
	err := mail.SendMailWithHTML(
		"Reset Password",
		fmt.Sprintf(`<p>Please reset your password by clicking on the button below</p><a style="display: inline-block; background-color: #136dfa; color: white; padding: 10px 20px; text-align: center; text-decoration: none; border-radius: 5px;" href="%s">Reset Password</a>`, resetURL),
		fromEmail,
		[]string{to},
		fromEmailPassword,
	)

	return err
}

// 
func init() {
	godotenv.Load()
	fromEmail = os.Getenv("FROM_EMAIL")
	fromEmailPassword = os.Getenv("FROM_EMAIL_PASSWORD")
	frontendURL = os.Getenv("FRONTEND_URL")
}
