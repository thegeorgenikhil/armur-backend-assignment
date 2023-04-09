package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/thegeorgenikhil/armur-backend-assignment/internal/database"
	"github.com/thegeorgenikhil/armur-backend-assignment/internal/models"
	"github.com/thegeorgenikhil/armur-backend-assignment/internal/utils"
	"github.com/thegeorgenikhil/armur-backend-assignment/pkg/bcrypt"
	"github.com/thegeorgenikhil/armur-backend-assignment/pkg/jwt"
)

var validate = validator.New()

// Register creates a new user in the database with all the necessary checks
func Register(w http.ResponseWriter, r *http.Request) {
	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		utils.RepondWithError(w, err, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	// checking for validation errors
	err = validate.Struct(user)

	if err != nil {
		utils.RespondWithValidationErrors(w, err)
		return
	}

	userExist := database.CheckIfUserExists(database.DB, user.Email)

	if userExist {
		utils.RepondWithError(w, err, "User already exists!", http.StatusBadRequest)
		return
	}

	// hashes the password before storing it in the database
	hashedPassword, err := bcrypt.HashPassword(user.Password)

	if err != nil {
		utils.RepondWithError(w, err, "Some issues with the server. Please try again", http.StatusInternalServerError)
		return
	}

	user.Password = hashedPassword

	// generate access and refresh tokens
	accessToken := jwt.GenerateToken(user.Email, jwt.AccessTokenDuration)
	refreshToken := jwt.GenerateToken(user.Email, jwt.RefreshTokenDuration)

	user.RefreshToken.Scan(refreshToken)

	// generate user activation token which will be sent to the user's email
	// on clicking the link in the email, the user will be verified
	// this is done by checking if the token is present at the time of registration is the same as the one in the user's email
	userActivationToken := jwt.GenerateToken(user.Email, jwt.UserActivationTokenDuration)

	user.UserVerificationToken.Scan(userActivationToken)

	// sqlite does not have a boolean type so we use an integer
	user.IsVerified = 0

	_ = database.InsertNewUser(database.DB, user)

	err = utils.SendUserVerificationEmail(user.Email, user.UserVerificationToken.String)

	if err != nil {
		log.Println(err)
	}

	// convert the user verified as a boolean value
	userVerifiedFlag := user.IsVerified != 0

	userCreatedReponse := &models.AuthenticatedUserResponse{
		AccessToken:    accessToken,
		RefreshToken:   user.RefreshToken.String,
		IsUserVerified: userVerifiedFlag,
		Msg:            "Signed up successfully. Please verify you email using the confirmation mail",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(userCreatedReponse)
}

// Login logs in the user and returns the access token and refresh token along with the user verified flag
func Login(w http.ResponseWriter, r *http.Request) {
	var loginUserRequest models.LoginUserRequest
	var user *models.User

	err := json.NewDecoder(r.Body).Decode(&loginUserRequest)

	if err != nil {
		utils.RepondWithError(w, err, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = validate.Struct(loginUserRequest)

	if err != nil {
		utils.RespondWithValidationErrors(w, err)
		return
	}

	// checking if user exists in the database
	ok := database.CheckIfUserExists(database.DB, loginUserRequest.Email)

	if !ok {
		utils.RepondWithError(w, err, "No user with the given email", http.StatusBadRequest)
		return
	}

	// fetching the user's details from the database
	user, err = database.GetUserWithEmail(database.DB, loginUserRequest.Email)

	if err != nil {
		utils.RepondWithError(w, err, "Not able to login. Try again!", http.StatusInternalServerError)
		return
	}

	// using bcrypt to verify the password hashes
	ok = bcrypt.VerifyPassword(user.Password, loginUserRequest.Password)

	// if the passwords do not match, return false
	if !ok {
		utils.RepondWithError(w, err, "Incorrect password", http.StatusUnauthorized)
		return
	}

	// we will generate new access and refresh tokens every time the user logs in
	accessToken := jwt.GenerateToken(user.Email, jwt.AccessTokenDuration)
	refreshToken := jwt.GenerateToken(user.Email, jwt.RefreshTokenDuration)

	// updating the refresh token in the database
	err = database.UpdateRefreshTokenInDB(database.DB, user.Email, refreshToken)

	if err != nil {
		utils.RepondWithError(w, err, "Not able to login. Try again!", http.StatusInternalServerError)
		return
	}

	userVerifiedFlag := user.IsVerified != 0

	loginUserResponse := &models.AuthenticatedUserResponse{
		Msg:            "Logged In Successfully",
		AccessToken:    accessToken,
		RefreshToken:   refreshToken,
		IsUserVerified: userVerifiedFlag,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(loginUserResponse)
}

// Logout removes the refresh token from the database thus making sure that the user cannot use the refresh token to generate a new access token, though the access token will still be valid till it expires
func Logout(w http.ResponseWriter, r *http.Request) {
	email := r.Context().Value(models.CtxKeyForUserEmail{}).(string)

	err := database.RemoveRefreshTokenFromDB(database.DB, email)

	if err != nil {
		utils.RepondWithError(w, err, "Not able to logout. Try again!", http.StatusInternalServerError)
		return
	}

	logoutUserResponse := &models.LogoutUserResponse{
		Msg: "Successfully logged out",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(logoutUserResponse)
}

// RefreshToken generates a new access token using the refresh token, this is done to avoid the user from having to login in each time the access token expires
func RefreshToken(w http.ResponseWriter, r *http.Request) {
	var refreshTokenRequest models.RefreshTokenRequest

	err := json.NewDecoder(r.Body).Decode(&refreshTokenRequest)

	if err != nil {
		utils.RepondWithError(w, err, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = validate.Struct(refreshTokenRequest)

	if err != nil {
		utils.RespondWithValidationErrors(w, err)
		return
	}

	claims, err := jwt.ValidateToken(refreshTokenRequest.RefreshToken)

	if err != nil {
		utils.RepondWithError(w, err, "Invalid refresh token. Please signin again", http.StatusBadRequest)
		return
	}

	user, err := database.GetUserWithEmail(database.DB, claims.Email)

	if err != nil {
		utils.RepondWithError(w, err, "Invalid refresh token", http.StatusBadRequest)
		return
	}

	if user.RefreshToken.String != refreshTokenRequest.RefreshToken {
		utils.RepondWithError(w, err, "Invalid refresh token. Please signin again", http.StatusBadRequest)
		return
	}

	// generating a new access token to be sent to the user
	accessToken := jwt.GenerateToken(user.Email, jwt.AccessTokenDuration)

	response := models.RefreshTokenResponse{
		AccessToken: accessToken,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
