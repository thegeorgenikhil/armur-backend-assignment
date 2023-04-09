package models

// ValidationError is the validation error model for the API, it is used to send validation errors to the client in case of any validation errors in the request body
type ValidationError struct {
	// Field is the field that is invalid
	Field string `json:"field"`
	// Msg is the error message
	Msg string `json:"msg"`
	// Param is the parameter that is invalid eg. len of a property
	Param string `json:"param,omitempty"`
}

// ValidationErrorResponse is the validation error response model for the API
type ValidationErrorResponse struct {
	Errors []ValidationError `json:"errors"`
}
