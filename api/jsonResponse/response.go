package jsonResponse

import (
	"net/http"
)

// JSONResponse describes a JSON API response.
type JSONResponse interface {
	JSON() ([]byte, error)         // JSON() returns the response body in JSON format.
	WriteTo(w http.ResponseWriter) // WriteTo() allows to write the response to a ResponseWriter
}

// ResponseMeta represents unnecessary (but useful) information about the response.
type ResponseMeta struct {
	Instance string `json:"instance"`       // The origin of a response (useful for tracing errors)
	Status   int    `json:"status"`         // The status of the request (allows to bypass frontend handling via response.status)
	Type     string `json:"type,omitempty"` // Type of a response (used with errors to easily identify the exact reason of the problem, making handling easier)
}
