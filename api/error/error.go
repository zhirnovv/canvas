package error

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// APIError represents an API error in accordance to RFC 7807.
type APIError struct {
	Type     string `json:"type"`     // Unique error type.
	Status   int    `json:"status"`   // HTTP status code to return.
	Title    string `json:"title"`    // Error title.
	Detail   string `json:"detail"`   // More information on the error.
	Instance string `json:"instance"` // Response .
}

// Error() implements the error interface for APIError
func (apiErr *APIError) Error() string {
	errorMessage := fmt.Sprintf("ERROR '%q' occured on '%q' with HTTP code %d: %s - %s", apiErr.Type, apiErr.Instance, apiErr.Status, apiErr.Title, apiErr.Detail)

	return errorMessage
}

// JSONError() returns an APIError instance in JSON format.
func (apiErr *APIError) JSONError() ([]byte, error) {
	errorJSON, parseError := json.Marshal(apiErr)

	if parseError != nil {
		return nil, fmt.Errorf("Error while parsing api error body: %v", parseError)
	}

	return errorJSON, nil
}

// WriteTo() sends the error to the request origin.
func (apiErr *APIError) WriteTo(w http.ResponseWriter) {
	errorJSON, parseError := apiErr.JSONError()

	if parseError != nil {
		var instance string

		if apiErr.Instance != "" {
			instance = apiErr.Instance
		} else {
			instance = "Unknown Instance"
		}

		apiParseError := NewAPIError("/internal/json/parse", http.StatusInternalServerError, "Error parsing JSON error body", parseError.Error(), instance)
		apiParseError.WriteTo(w)

		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(errorJSON)
}

func NewAPIError(errorType string, status int, title string, detail string, instance string) *APIError {
	return &APIError{
		Type:     errorType,
		Status:   status,
		Title:    title,
		Detail:   detail,
		Instance: instance,
	}
}
