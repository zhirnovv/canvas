package toolkit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-playground/validator"
	"github.com/zhirnovv/canvas/api/jsonResponse"
)

type TestHeader struct {
	Key   string
	Value string
}

type HandlerTest struct {
	Name           string
	Method         string
	Body           string
	Headers        []TestHeader
	ExpectedError  *jsonResponse.ErrorResponse
	ExpectedResult *jsonResponse.SuccessResponse
}

var validate = validator.New()

// AssertResponse is a function which generates tests for a certain handler with a provided test table.
func TestHandler(t *testing.T, handler http.HandlerFunc, testTable []HandlerTest) {
	for _, testCase := range testTable {
		t.Run(testCase.Name, func(t *testing.T) {
			request := httptest.NewRequest(testCase.Method, "/auth", strings.NewReader(testCase.Body))
			for _, header := range testCase.Headers {
				request.Header.Add(header.Key, header.Value)
			}

			responseRecorder := httptest.NewRecorder()
			handler(responseRecorder, request)
			responseBody := responseRecorder.Body.String()

			if testCase.ExpectedError != nil {
				var response jsonResponse.ErrorResponse
				err := json.Unmarshal([]byte(responseBody), &response)

				if err != nil {
					t.Errorf(err.Error())
				}

				if testCase.ExpectedError.Meta.Type != response.Meta.Type {
					t.Errorf("Want type '%s', got type '%s'", testCase.ExpectedError.Meta.Type, response.Meta.Type)
				}

				if testCase.ExpectedError.Meta.Status != response.Meta.Status {
					t.Errorf("Want status '%d', got status '%d'", testCase.ExpectedError.Meta.Status, response.Meta.Status)
				}
			} else if testCase.ExpectedResult != nil {
				expectedJSON, _ := testCase.ExpectedResult.JSON()

				if string(expectedJSON) != responseBody {
					t.Errorf("Want body '%s', got body '%s'", expectedJSON, responseBody)
				}
			}

		})
	}
}
