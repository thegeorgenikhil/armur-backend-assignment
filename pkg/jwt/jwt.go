package jwt

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)

const (
	AccessTokenDuration         = time.Hour * 24
	RefreshTokenDuration        = time.Hour * 24 * 5
	UserActivationTokenDuration = time.Minute * 10
	PasswordResetTokenDuration  = time.Minute * 10
)

var (
	secretKey       string
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

// Claims contains the email and the standard claims for the jwt
type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

// GenerateToken creates a new token by taking in the user email and secret
func GenerateToken(email string, t time.Duration) string {
	claims := &Claims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Local().Add(t).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secretKey))

	if err != nil {
		log.Fatal(err)
	}

	return token
}

// ValidateToken returns the Claims struct if the given token is a valid token, error if not
func ValidateToken(signedToken string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		},
	)

	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payload, ok := token.Claims.(*Claims)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}

func init() {
	godotenv.Load()
	secretKey = os.Getenv("JWT_SECRET")
}
