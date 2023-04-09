package models

import "database/sql"

// https://stackoverflow.com/questions/40891345/fix-should-not-use-basic-type-string-as-key-in-context-withvalue-golint
type CtxKeyForUserEmail struct{}

// User is the user model
type User struct {
	UserId                string         `json:"_" db:"uid"`
	FirstName             string         `json:"first_name" db:"first_name" validate:"required,max=50"`
	LastName              string         `json:"last_name" db:"last_name" validate:"required,max=50"`
	Email                 string         `json:"email" db:"email" validate:"required,email"`
	Password              string         `json:"password" db:"password" validate:"required,min=3"`
	PhoneNumber           string         `json:"phone_number" db:"phone_number" validate:"required,len=10"`
	Country               string         `json:"country" db:"country" validate:"required"`
	JobType               string         `json:"job_type" db:"job_type" validate:"required"`
	IncomeRange           string         `json:"income_range" db:"income_range" validate:"required"`
	RefreshToken          sql.NullString `json:"refresh_token" db:"refresh_token"`
	UserVerificationToken sql.NullString `db:"user_verification_token"`
	IsVerified            int            `db:"is_verified"`
}

// Response is the response model for all the handlers sending a response
type Response struct {
	Msg string `json:"msg"`
}

// AuthenticatedUserResponse is the response model for the login handler
type AuthenticatedUserResponse struct {
	Msg            string `json:"msg"`
	AccessToken    string `json:"access_token"`
	RefreshToken   string `json:"refresh_token"`
	IsUserVerified bool   `json:"is_verified"`
}

// LoginUserRequest is the request model for the login handler
type LoginUserRequest struct {
	Email    string `json:"email" db:"email" validate:"required,email"`
	Password string `json:"password" db:"password" validate:"required,min=3"`
}

// LoginUserResponse is the response model for the logout handler
type LogoutUserResponse struct {
	Msg string `json:"msg"`
}

// RefreshTokenRequest is the request model for the refresh token handler
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// RefreshTokenResponse is the response model for the refresh token handler
type RefreshTokenResponse struct {
	AccessToken string `json:"access_token" vaiidate:"required"`
}

// UserVerificationStatusResponse is the request model for check user verification status handler
type UserVerificationStatusResponse struct {
	IsUserVerified bool `json:"is_user_verified"`
}

// ForgotPasswordRequest is the request model for the forgot password handler
type ForgotPasswordRequest struct {
	Email string `json:"email" db:"email" validate:"required,email"`
}

// ResetPasswordRequest is the request model for the reset password handler
type ResetPasswordRequest struct {
	Email      string `json:"email" db:"email" validate:"required,email"`
	ResetToken string `json:"token" db:"token" validate:"required"`
	Password   string `json:"password" db:"password" validate:"required,min=3"`
}

// GetUserResponse is the response model for the get user handler
type GetUserResponse struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Country     string `json:"country"`
	JobType     string `json:"job_type"`
	IncomeRange string `json:"income_range"`
	IsVerified  bool   `json:"is_verified"`
}
