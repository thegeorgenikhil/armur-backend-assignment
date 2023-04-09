package utils

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/thegeorgenikhil/armur-backend-assignment/internal/models"
)

// RespondWithError sets the writer with suitable error response
func RepondWithError(w http.ResponseWriter, err error, msg string, status int) {
	log.Println("Error: ", err)
	errorResponse := models.ErrorResponse{
		Message: msg,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(errorResponse)
}

// ResponsdWithValidationErrors sets the writer with suitable error response for all the valiation error
func RespondWithValidationErrors(w http.ResponseWriter, err error) {
	// creating a seperate error slice to return back all validation errors
	errors := []models.ValidationError{}
	for _, err := range err.(validator.ValidationErrors) {
		validationErr := models.ValidationError{
			Field: err.Field(),
			Msg:   err.Tag(),
			Param: err.Param(),
		}

		errors = append(errors, validationErr)
	}

	errorResponse := models.ValidationErrorResponse{
		Errors: errors,
	}
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(errorResponse)
}
