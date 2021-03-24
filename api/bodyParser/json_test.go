package bodyParser

import (
	"net/http"
	"testing"

	"github.com/zhirnovv/canvas/api/jsonResponse"
	"github.com/zhirnovv/canvas/api/toolkit"
)

func TestJsonParse(t *testing.T) {
	testTable := []toolkit.HandlerTest{
		{
			Name:           "Missing header",
			Method:         http.MethodPost,
			Body:           "",
			ExpectedError:  jsonResponse.NewErrorResponse("/request/headers/contentType", http.StatusBadRequest, "", "", "/auth"),
			ExpectedResult: nil,
		},
		{
			Name:   "Missing Body",
			Method: http.MethodPost,
			Body:   "",
			Headers: []toolkit.TestHeader{
				{
					Key:   "Content-Type",
					Value: "application/json",
				},
			},
			ExpectedError:  jsonResponse.NewErrorResponse("/request/body/parse/failed", http.StatusBadRequest, "", "", "/auth"),
			ExpectedResult: nil,
		},
		{
			Name:   "Incorrect JSON format",
			Method: http.MethodPost,
			Body:   "{ \"test\": \"2\", anotherTest: 2 }",
			Headers: []toolkit.TestHeader{
				{
					Key:   "Content-Type",
					Value: "application/json",
				},
			},
			ExpectedError:  jsonResponse.NewErrorResponse("/request/body/parse/failed", http.StatusBadRequest, "", "", "/auth"),
			ExpectedResult: nil,
		},
		{
			Name:   "Everything is correct",
			Method: http.MethodPost,
			Body:   "{ \"test\": ....2 }",
			Headers: []toolkit.TestHeader{
				{
					Key:   "Content-Type",
					Value: "application/json",
				},
			},
			ExpectedError:  nil,
			ExpectedResult: nil,
		},
	}

	toolkit.TestHandler(t, testHandler, testTable)
}
