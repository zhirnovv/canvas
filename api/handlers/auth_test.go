package handlers

import (
	"net/http"
	"testing"

	"github.com/zhirnovv/canvas/api/auth"
	"github.com/zhirnovv/canvas/api/jsonResponse"
	"github.com/zhirnovv/canvas/api/toolkit"
	"github.com/zhirnovv/canvas/api/user"
)

func TestSignupHandler(t *testing.T) {
	testTable := []toolkit.HandlerTest{
		{
			Name:   "Missing header",
			Method: http.MethodPost,
			Body:   "",
			ExpectedError:  jsonResponse.NewErrorResponse("/request/headers/contentType", http.StatusBadRequest, "", "", "/auth"),
			ExpectedResult: nil,
		},
		{
			Name:   "Missing body",
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
			Name:   "Incorrect body",
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
			Name:   "Incorrect key",
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
			Name: "Incorrect value",
			Method: http.MethodPost,
			Headers: []toolkit.TestHeader{
				{
					Key:   "Content-Type",
					Value: "application/json",
				},
			},
			Body:   "{\"ntme\": \"/.,/,/.,/.\"}",
			ExpectedError:  jsonResponse.NewErrorResponse("/request/body/validate/failed", http.StatusBadRequest, "", "", "/auth"),
			ExpectedResult: nil,
		},
		{
			Name:   "Correct body",
			Method: http.MethodPost,
			Body:   "{\"name\": \"Johnny Test\"}",
			Headers: []toolkit.TestHeader{
				{
					Key:   "Content-Type",
					Value: "application/json",
				},
			},
			ExpectedError: nil,
			ExpectedResult: jsonResponse.NewSuccessResponse(
				map[string]interface{}{
					"data": map[string]string{
						"message": "Successfully created user.",
					},
				},
				http.StatusOK,
				"/auth",
			),
		},
	}

	userStorage := user.NewUserStorage()
	authenticator := auth.NewUserAuthenticator(userStorage, "test")

	handler := SignupHandler(userStorage, authenticator)

	toolkit.TestHandler(t, handler, testTable)
}
