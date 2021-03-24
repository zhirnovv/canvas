package jsonResponse

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// SuccessResponse represents an API Response in accordance to JSON:APIv1 and RFC 7807 hybrid spec.
type SuccessResponse struct {
	Data map[string]interface{} `json:"data"`
	Meta ResponseMeta           `json:"meta"`
}

// NewSuccessResponse() returns a new SuccessResponse instance.
func NewSuccessResponse(data map[string]interface{}, status int, instance string) *SuccessResponse {
	return &SuccessResponse{
		Data: data,
		Meta: ResponseMeta{
			Type: "",
			Status: status,
			Instance: instance,
		},
	}
}

// JSON() implements the Response interface
func (successBody *SuccessResponse) JSON() ([]byte, error) {
	successBodyJSON, parseError := json.Marshal(successBody)

	if parseError != nil {
		return nil, fmt.Errorf("Error while parsing api success body: %v", parseError)
	}

	return successBodyJSON, nil
}

// WriteTo() implements the Response interface
func (successBody *SuccessResponse) WriteTo(w http.ResponseWriter) {
	successBodyJSON, parseErr := successBody.JSON()

	if parseErr != nil {
		var instance string

		if successBody.Meta.Instance != "" {
			instance = successBody.Meta.Instance
		} else {
			instance = "Unknown Instance"
		}

		apiParseError := NewErrorResponse("/internal/json/parse", http.StatusInternalServerError, "Error parsing JSON error body", parseErr.Error(), instance)
		apiParseError.WriteTo(w)

		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(successBody.Meta.Status)
	w.Write(successBodyJSON)
}
