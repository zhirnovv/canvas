package bodyParser

import (
	"encoding/json"
	"github.com/zhirnovv/canvas/api/jsonResponse"
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

// ParseJSON() unmarshals a JSON body of request r into a variable pointed at by body. Requires the request to possess a body and have Content-Type=application/json.
func ParseJSON(r *http.Request, schema interface{}) *jsonResponse.ErrorResponse {
	if !hasCorrectContentType(r) {
		badContentTypeError := jsonResponse.NewErrorResponse("/request/headers/contentType", http.StatusBadRequest, "Incorrect or missing content-type header", "The following request uses JSONParserMiddleware, which requires content-type to be application/json", r.RequestURI)

		return badContentTypeError
	}

	requestBody, readErr := ioutil.ReadAll(r.Body)

	if readErr != nil {
		bodyReadError := jsonResponse.NewErrorResponse("/request/body/couldNotRead", http.StatusBadRequest, "Could not read request body", "", r.RequestURI)

		return bodyReadError
	}

	var parsedBody map[string]interface{}

	parseErr := json.Unmarshal(requestBody, &parsedBody)

	if parseErr != nil {
		parseError := jsonResponse.NewErrorResponse("/request/body/parse/failed", http.StatusBadRequest, "Failed to parse request body.", parseErr.Error(), r.RequestURI)

		return parseError
	}

	validateErr := DecodeAndValidate(parsedBody, schema)

	if validateErr != nil { 
		validateError := jsonResponse.NewErrorResponse("/request/body/validate/failed", http.StatusBadRequest, "Failed to decode and validate request body.", validateErr.Error(), r.RequestURI)

		return validateError
	}

	return nil
}
