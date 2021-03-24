package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/zhirnovv/canvas/api/auth"
	"github.com/zhirnovv/canvas/api/jsonResponse"
)

type key string
const UUIDKey key = "userUUID"

// RequireUserAuthentication is a middleware that accepts an Authenticator.
// If the user authentication cookie does not exists of the Authenticator invalidates it, throw an API error.
func RequireUserAuthentication(authenticator auth.Authenticator, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var authResponse jsonResponse.JSONResponse

		userAuthCookie, err := r.Cookie("user_auth_token")

		if err != nil {
			authResponse = jsonResponse.NewErrorResponse("/auth/noToken", http.StatusUnauthorized, "Failed to read user authentication token", err.Error(), r.RequestURI)
			authResponse.WriteTo(w)
			return
		}

		decodedTokenValue, authErr := authenticator.VerifyAndDecode(userAuthCookie.Value)

		if authErr != nil {
			authResponse = jsonResponse.NewErrorResponse("/auth/invalidToken", http.StatusUnauthorized, "Invalid authentication token.", authErr.Error(), r.RequestURI)
			authResponse.WriteTo(w)
			return
		}

		uuid, isValidUUID := decodedTokenValue.(uuid.UUID)

		if !isValidUUID {
			authResponse = jsonResponse.NewErrorResponse("/auth/invalidTokenValue", http.StatusUnauthorized, fmt.Sprintf("Decoded token value %v is not a valid UUID", decodedTokenValue), "", r.RequestURI)
			authResponse.WriteTo(w)
			return
		}

		newContext := context.WithValue(r.Context(), UUIDKey, uuid)
		next(w, r.WithContext(newContext))
	}
}
