package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/thegeorgenikhil/armur-backend-assignment/internal/database"
	"github.com/thegeorgenikhil/armur-backend-assignment/internal/models"
	"github.com/thegeorgenikhil/armur-backend-assignment/internal/utils"
	"github.com/thegeorgenikhil/armur-backend-assignment/pkg/jwt"
)

// VerifyUser verifies the user by checking the user verification token that is sent to the user's email. It then returns a response saying that the user has been verified, also removes the user verification token from the database and updates the user's verification status to verified
func VerifyUser(w http.ResponseWriter, r *http.Request) {
	confirmationCode := r.URL.Query().Get("token")

	claims, err := jwt.ValidateToken(confirmationCode)

	if err != nil {
		http.Error(w, "Verification token has expired", http.StatusBadRequest)
		return
	}

	user, err := database.GetUserWithEmail(database.DB, claims.Email)

	if err != nil {
		http.Error(w, "Verification link invalid", http.StatusBadRequest)
		return
	}

	// check for user already verified
	if user.IsVerified != 0 {
		http.Error(w, "User already verified", http.StatusBadRequest)
		return
	}

	// check for user activation token validity
	if user.UserVerificationToken.String != confirmationCode {
		http.Error(w, "Invalid verification code", http.StatusBadRequest)
		return
	}

	// update the user to verified
	err = database.VerifyUserInDB(database.DB, user.UserId)

	if err != nil {
		http.Error(w, "Not able to verify user. Try again!", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Verification successful. You can close this tab and continue to the home page"))
}

// ResendUserVerificationEmail generates a new user verification token and sends it to the user's email. It also updates the user verification token in the database. It also checks if the user has already verified their account or if the user has requested for a new verification email in the last 5 minutes (300 seconds)
func ResendUserVerificationEmail(w http.ResponseWriter, r *http.Request) {
	email := r.Context().Value(models.CtxKeyForUserEmail{}).(string)

	user, err := database.GetUserWithEmail(database.DB, email)

	if err != nil {
		utils.RepondWithError(w, err, "Some issues with the sever. Please try again", http.StatusInternalServerError)
		return
	}

	if user.IsVerified != 0 {
		utils.RepondWithError(w, err, "User already verified", http.StatusBadRequest)
		return
	}

	claims, err := jwt.ValidateToken(user.UserVerificationToken.String)

	// check if the user has already requested for a new verification email in the last 5 minutes
	secondsPassed := time.Now().Unix()-claims.ExpiresAt > 300

	if err == nil && !secondsPassed {
		utils.RepondWithError(w, err, "Too many requests. Please wait for 5 minutes before requesting for a new verification email", http.StatusTooManyRequests)
		return
	}

	userVerificationToken := jwt.GenerateToken(user.Email, jwt.UserActivationTokenDuration)
	err = utils.SendUserVerificationEmail(user.Email, userVerificationToken)
	if err != nil {
		log.Println(err)
	}

	err = database.UpdateVerificationTokenInDB(database.DB, user.Email, userVerificationToken)

	if err != nil {
		utils.RepondWithError(w, err, "Some issues with the sever. Please try again", http.StatusInternalServerError)
		return
	}

	response := models.Response{
		Msg: "Resent verification email. Please check your registered email",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// CheckUserVerificationStatus checks if the user is verified or not and returns a response with the verification status
func CheckUserVerificationStatus(w http.ResponseWriter, r *http.Request) {
	email := r.Context().Value(models.CtxKeyForUserEmail{}).(string)

	user, err := database.GetUserWithEmail(database.DB, email)

	if err != nil {
		utils.RepondWithError(w, err, "Some issues with the sever. Please try again", http.StatusInternalServerError)
		return
	}

	response := models.UserVerificationStatusResponse{
		IsUserVerified: user.IsVerified != 0,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
