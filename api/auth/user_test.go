package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/zhirnovv/canvas/api/user"
)

var userStorage = user.NewUserStorage()
var userAuthenticator = NewUserAuthenticator(userStorage, "secret")

func TestUserAuthenticatorTokenIssue(t *testing.T) {
	newUser, err := userAuthenticator.Storage.Create("test user")

	if err != nil {
		t.Error(err)
	}

	tokenString, err := userAuthenticator.Issue(newUser.ID)

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
	newUser, err := userAuthenticator.Storage.Create("test user")

	if err != nil {
		t.Error(err)
	}

	token, err := userAuthenticator.Issue(newUser.ID)

	if err != nil {
		t.Error(err)
	}

	decodedUuid, err := userAuthenticator.VerifyAndDecode(token)

	if err != nil {
		t.Error(err)
	}

	if newUser.ID != decodedUuid {
		t.Errorf("UUID and decoded UUID don't match. UUID is %s, while decoded UUID is %s", newUser.ID.String(), decodedUuid.(uuid.UUID).String())
	}
}

func TestUserAuthenticatorMangledTokenDecode(t *testing.T) {
	newUser, err := userAuthenticator.Storage.Create("test user")

	if err != nil {
		t.Error(err)
	}

	token, err := userAuthenticator.Issue(newUser.ID)

	if err != nil {
		t.Error(err)
	}

	token = token[:len(token)-1]

	_, err = userAuthenticator.VerifyAndDecode(token)

	if err == nil {
		t.Error("The token was altered, but no error was thrown")
	}

}

func TestTokenExpiration(t *testing.T) {
	newUser, err := userAuthenticator.Storage.Create("test user")

	if err != nil {
		t.Error(err)
	}

	token, err := userAuthenticator.Issue(newUser.ID)

	if err != nil {
		t.Error(err)
	}

	token = token[:len(token)-1]

	time.Sleep(time.Second * 6)

	_, err = userAuthenticator.VerifyAndDecode(token)

	if err == nil {
		t.Error("The token is expired, but no error was thrown")
	}
}
