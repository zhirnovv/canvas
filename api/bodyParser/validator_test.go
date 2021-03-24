package bodyParser

import (
	"net/http"
	"testing"

	"github.com/zhirnovv/canvas/api/jsonResponse"
	"github.com/zhirnovv/canvas/api/toolkit"
)

type testForm struct {
	Name   string `mapstruct:"name" validate:"eq=4,alphanum"`
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	var form testForm

	parseErr := ParseJSON(r, &form)

	if parseErr != nil {
		parseErr.WriteTo(w)
	}

	return
}

func TestDecodeAndValidate(t *testing.T) {
	testTable := []toolkit.HandlerTest{
		{
			Name:   "Incorrect Body",
			Method: http.MethodPost,
			Body:   "{}",
			Headers: []toolkit.TestHeader{
				{
					Key:   "Content-Type",
					Value: "application/json",
				},
			},
			ExpectedError:  jsonResponse.NewErrorResponse("/request/body/validate/failed", http.StatusBadRequest, "", "", "/auth"),
			ExpectedResult: nil,
		},
		{
			Name:   "Incorrect Key",
			Method: http.MethodPost,
			Body:   "{\"ntme\": \"test\"}",
			Headers: []toolkit.TestHeader{
				{
					Key:   "Content-Type",
					Value: "application/json",
				},
			},
			ExpectedError:  jsonResponse.NewErrorResponse("/request/body/validate/failed", http.StatusBadRequest, "", "", "/auth"),
			ExpectedResult: nil,
		},
		{
			Name:   "Incorrect Value",
			Method: http.MethodPost,
			Headers: []toolkit.TestHeader{
				{
					Key:   "Content-Type",
					Value: "application/json",
				},
			},
			Body:           "{\"ntme\": \"/.,/,/.,/.\"}",
			ExpectedError:  jsonResponse.NewErrorResponse("/request/body/validate/failed", http.StatusBadRequest, "", "", "/auth"),
			ExpectedResult: nil,
		},
	}

	toolkit.TestHandler(t, testHandler, testTable)
}
