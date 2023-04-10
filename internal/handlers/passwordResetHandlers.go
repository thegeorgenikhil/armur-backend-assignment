package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/thegeorgenikhil/armur-backend-assignment/internal/database"
	"github.com/thegeorgenikhil/armur-backend-assignment/internal/models"
	"github.com/thegeorgenikhil/armur-backend-assignment/internal/utils"
	"github.com/thegeorgenikhil/armur-backend-assignment/pkg/bcrypt"
	"github.com/thegeorgenikhil/armur-backend-assignment/pkg/jwt"
)

// ForgotPassword sends a password reset token and sends an email to the user with the token, on clicking the link the user will be redirected to the reset password page where he can reset his password
func ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var forgotPasswordRequest models.ForgotPasswordRequest

	err := json.NewDecoder(r.Body).Decode(&forgotPasswordRequest)

	if err != nil {
		utils.RepondWithError(w, err, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = validate.Struct(forgotPasswordRequest)

	if err != nil {
		utils.RespondWithValidationErrors(w, err)
		return
	}

	ok := database.CheckIfUserExists(database.DB, forgotPasswordRequest.Email)
	if !ok {
		utils.RepondWithError(w, err, "Invalid email given", http.StatusNotFound)
		return
	}

	// if user has already requested for a password reset token before, we will remove the previous token and generate a new one
	err = database.RemoveResetPasswordTokenFromDB(database.DB, forgotPasswordRequest.Email)

	if err != nil {
		utils.RepondWithError(w, err, "Some issues with the sever. Please try again", http.StatusInternalServerError)
		return
	}

	// we are using jwt to generate a random set of characters which will be used as the password reset token
	passwordResetToken := jwt.GenerateToken(forgotPasswordRequest.Email, jwt.PasswordResetTokenDuration)

	passwordResetToken = strings.Split(passwordResetToken, ".")[0]

	hashedResetToken, err := bcrypt.HashPassword(passwordResetToken)

	// password reset token will be valid for 10 minutes
	tokenDuration := time.Now().Add(jwt.PasswordResetTokenDuration).Unix()

	if err != nil {
		utils.RepondWithError(w, err, "Some issues with the sever. Please try again", http.StatusInternalServerError)
		return
	}

	err = database.InsertResetPasswordTokenInDB(database.DB, forgotPasswordRequest.Email, hashedResetToken, tokenDuration)

	if err != nil {
		utils.RepondWithError(w, err, "Some issues with the sever. Please try again", http.StatusInternalServerError)
		return
	}

	// sending the password reset token to the user's email with the token we just generated
	err = utils.SendPasswordResetEmail(forgotPasswordRequest.Email, hashedResetToken)

	if err != nil {
		log.Println(err)
	}

	response := models.Response{
		Msg: "Reset password email sent. Please check your email",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// ResetPassword resets the password of the user. It takes the password user email, reset token and the new password as input and resets the password of the user. It gets the reset token from the email link we just sent in the above handler. It then checks if the reset token is same as the one in the database. If yes then it resets the password of the user and removes the reset token from the database
func ResetPassword(w http.ResponseWriter, r *http.Request) {
	var resetPasswordRequest models.ResetPasswordRequest

	err := json.NewDecoder(r.Body).Decode(&resetPasswordRequest)

	if err != nil {
		utils.RepondWithError(w, err, "Invalid request body", http.StatusBadRequest)
		return
	}

	// checking for validation errors
	err = validate.Struct(resetPasswordRequest)

	if err != nil {
		utils.RespondWithValidationErrors(w, err)
		return
	}

	ok := database.CheckIfUserExists(database.DB, resetPasswordRequest.Email)

	if !ok {
		utils.RepondWithError(w, err, "Invalid email given", http.StatusNotFound)
		return
	}

	resetToken, resetTokenDuration, err := database.GetResetPasswordTokenFromDB(database.DB, resetPasswordRequest.Email)

	if err != nil {
		utils.RepondWithError(w, err, "Some issues with the sever. Please try again", http.StatusInternalServerError)
		return
	}

	if resetToken == "" {
		utils.RepondWithError(w, err, "No reset password token found. Please request for a new one", http.StatusNotFound)
		return
	}

	if resetTokenDuration < time.Now().Unix() {
		utils.RepondWithError(w, err, "Reset password token expired. Please request for a new one", http.StatusNotFound)
		return
	}

	if resetToken != resetPasswordRequest.ResetToken {
		utils.RepondWithError(w, err, "Invalid reset password token", http.StatusNotFound)
		return
	}

	hashedPassword, err := bcrypt.HashPassword(resetPasswordRequest.Password)

	if err != nil {
		utils.RepondWithError(w, err, "Some issues with the sever. Please try again", http.StatusInternalServerError)
		return
	}

	err = database.UpdateUserPasswordInDB(database.DB, resetPasswordRequest.Email, hashedPassword)

	if err != nil {
		utils.RepondWithError(w, err, "Not able to reset password. Please try again", http.StatusInternalServerError)
		return
	}

	err = database.RemoveResetPasswordTokenFromDB(database.DB, resetPasswordRequest.Email)

	if err != nil {
		utils.RepondWithError(w, err, "Some issues with the sever. Please try again", http.StatusInternalServerError)
		return
	}

	response := models.Response{
		Msg: "Password reset successful. You can now login with your new password",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
