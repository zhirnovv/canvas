package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

var userAuthenticator = NewUserAuthenticator("secret")

func TestUserAuthenticatorTokenIssue(t *testing.T) {
	uuid := uuid.New()

	tokenString, err := userAuthenticator.IssueToken(uuid)

	if err != nil {
		t.Error(err)
	}

	hasSeparators := false

	for char := range tokenString {
		if char == 46 {
			hasSeparators = true
		}
	}

	if !hasSeparators {
		t.Error("JWT does not contain separators.")
	}
}

func TestUserAuthenticator(t *testing.T) {
	uuid := uuid.New()

	token, err := userAuthenticator.IssueToken(uuid)

	if err != nil {
		t.Error(err)
	}

	decodedUuid, err := userAuthenticator.DecodeToken(token)

	if err != nil {
		t.Error(err)
	}

	if uuid != decodedUuid {
		t.Errorf("UUID and decoded UUID don't match. UUID is %s, while decoded UUID is %s", uuid.String(), decodedUuid.String())
	}
}

func TestUserAuthenticatorMangledTokenDecode(t *testing.T) {
	uuid := uuid.New()

	token, err := userAuthenticator.IssueToken(uuid)

	if err != nil {
		t.Error(err)
	}

	token = token[:len(token)-1]

	_, err = userAuthenticator.DecodeToken(token)

	if err == nil {
		t.Error("The token was altered, but no error was thrown")
	}

}

func TestTokenExpiration(t *testing.T) {
	uuid := uuid.New()

	token, err := userAuthenticator.IssueToken(uuid)

	if err != nil {
		t.Error(err)
	}

	token = token[:len(token)-1]

	time.Sleep(time.Second * 6)

	_, err = userAuthenticator.DecodeToken(token)

	if err == nil {
		t.Error("The token is expired, but no error was thrown")
	}
}
