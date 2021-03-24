package jsonResponse

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ErrorResponseBody represents the error response body.
type ErrorResponseBody struct {
	Title  string `json:"title"`  // Short info on the error.
	Detail string `json:"detail"` // Further elaboration of the error.
}

// ErrorResponse represents an API error in accordance to JSON:APIv1 and RFC 7807 hybrid spec.
type ErrorResponse struct {
	Body ErrorResponseBody `json:"error"`
	Meta ResponseMeta      `json:"meta"`
}

// NewErrorResponse() returns a new ErrorResponse instance.
func NewErrorResponse(errorType string, status int, title string, detail string, instance string) *ErrorResponse {
	return &ErrorResponse{
		Body: ErrorResponseBody{
			Title:  title,
			Detail: detail,
		},
		Meta: ResponseMeta{
			Type:     errorType,
			Status:   status,
			Instance: instance,
		},
	}
}

// Error() implements the error interface.
func (err *ErrorResponse) Error() string {
	errorMessage := fmt.Sprintf("ERROR '%q' occured on '%q' with HTTP code %d: %s - %s", err.Meta.Type, err.Meta.Instance, err.Meta.Status, err.Body.Title, err.Body.Detail)

	return errorMessage
}

// JSON() implements the Response interface
func (errBody *ErrorResponse) JSON() ([]byte, error) {
	errorJSON, parseError := json.Marshal(errBody)

	if parseError != nil {
		return nil, fmt.Errorf("Error while parsing api error body: %v", parseError)
	}

	return errorJSON, nil
}

// WriteTo() implements the Response interface
func (errBody *ErrorResponse) WriteTo(w http.ResponseWriter) {
	errorJSON, parseError := errBody.JSON()

	if parseError != nil {
		var instance string

		if errBody.Meta.Instance != "" {
			instance = errBody.Meta.Instance
		} else {
			instance = "Unknown Instance"
		}

		apiParseError := NewErrorResponse("/internal/json/parse", http.StatusInternalServerError, "Error parsing JSON error body", parseError.Error(), instance)
		apiParseError.WriteTo(w)

		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(errBody.Meta.Status)
	w.Write(errorJSON)
}
