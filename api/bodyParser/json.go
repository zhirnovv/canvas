package bodyParser

import (
	"encoding/json"
	"github.com/zhirnovv/gochat/api/error"
	"io/ioutil"
	"net/http"
)

func hasCorrectContentType(r *http.Request) bool {
	contentType, ok := r.Header["Content-Type"]

	if !ok {
		return false
	}

	for _, value := range contentType {
		if value == "application/json" {
			return true
		}
	}

	return false
}

// ParseJSON() unmarshals a request body in JSON format. Requires the request to possess a body and have Content-Type=application/json
func ParseJSON(r *http.Request) (map[string]interface{}, *error.APIError) {
	if !hasCorrectContentType(r) {
		badContentTypeError := error.NewAPIError("/request/headers/contentType", http.StatusBadRequest, "Incorrect or missing content-type header", "The following request uses JSONParserMiddleware, which requires content-type to be application/json", r.RequestURI)

		return nil, badContentTypeError
	}

	requestBody, readErr := ioutil.ReadAll(r.Body)

	if readErr != nil {
		bodyReadError := error.NewAPIError("/request/body/couldNotRead", http.StatusBadRequest, "Could not read request body", "", r.RequestURI)

		return nil, bodyReadError
	}

	var parsedBody map[string]interface{}

	parseErr := json.Unmarshal(requestBody, &parsedBody)

	if parseErr != nil {
		parseError := error.NewAPIError("/request/body/parse/failed", http.StatusBadRequest, "Failed to parse request body.", parseErr.Error(), r.RequestURI)

		return nil, parseError
	}

	return parsedBody, nil
}
