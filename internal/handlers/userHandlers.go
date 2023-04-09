package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/thegeorgenikhil/armur-backend-assignment/internal/database"
	"github.com/thegeorgenikhil/armur-backend-assignment/internal/models"
	"github.com/thegeorgenikhil/armur-backend-assignment/internal/utils"
)

// GetUser returns the user details, if the user is verified
func GetUser(w http.ResponseWriter, r *http.Request) {
	email := r.Context().Value(models.CtxKeyForUserEmail{}).(string)

	user, err := database.GetUserWithEmail(database.DB, email)

	if err != nil {
		utils.RepondWithError(w, err, "Some issues with the sever. Please try again", http.StatusInternalServerError)
		return
	}

	if user.IsVerified == 0 {
		utils.RepondWithError(w, err, "User is not verified", http.StatusUnauthorized)
		return
	}

	response := models.GetUserResponse{
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		PhoneNumber: user.PhoneNumber,
		Country:     user.Country,
		JobType:     user.JobType,
		IncomeRange: user.IncomeRange,
		IsVerified:  user.IsVerified != 0,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
