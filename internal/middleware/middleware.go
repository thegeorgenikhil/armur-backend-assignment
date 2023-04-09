package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/thegeorgenikhil/armur-backend-assignment/internal/database"
	"github.com/thegeorgenikhil/armur-backend-assignment/internal/models"
	"github.com/thegeorgenikhil/armur-backend-assignment/internal/utils"
	"github.com/thegeorgenikhil/armur-backend-assignment/pkg/jwt"
)

// AuthMiddleware is a middleware that checks if the user is authenticated by checking the Authorization header, it then passes the email to the request context which can be used by the other handlers
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")

		if authorizationHeader == "" {
			utils.RepondWithError(w, nil, "Authorization header not found", http.StatusUnauthorized)
			return
		}

		// The middleware expects the Authorization Header to be in the form: Bearer <token>
		accessToken := strings.TrimPrefix(authorizationHeader, "Bearer")
		accessToken = strings.TrimSpace(accessToken)

		if accessToken == "" {
			utils.RepondWithError(w, nil, "Access token not found", http.StatusUnauthorized)
			return
		}

		claims, err := jwt.ValidateToken(accessToken)

		if err != nil {
			utils.RepondWithError(w, err, err.Error(), http.StatusUnauthorized)
			return
		}

		userExists := database.CheckIfUserExists(database.DB, claims.Email)

		if !userExists {
			utils.RepondWithError(w, err, "User doesn't exist", http.StatusBadRequest)
			return
		}

		// Linter error when doing the below way : should not use built-in type string as key for value; define your own type to avoid collisions (SA1029)go-staticcheck
		// rWithEmail := r.WithContext(context.WithValue(r.Context(), "email", claims.Email))
		// Explanation: // https://stackoverflow.com/questions/40891345/fix-should-not-use-basic-type-string-as-key-in-context-withvalue-golint

		rWithEmail := r.WithContext(context.WithValue(r.Context(), models.CtxKeyForUserEmail{}, claims.Email))

		next.ServeHTTP(w, rWithEmail)
	})
}
