package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/zhirnovv/canvas/api/user"
)

// UserAuthenticator contains all necessary parameters for signing and parsing JSON Web Tokens.
type UserAuthenticator struct {
	Storage       *user.UserStorage // UserStorage attached to UserAuthenticator
	secret        string            // Signing key. Under no circumstances should this value be exportable.
	defaultClaims jwt.MapClaims
	signingMethod jwt.SigningMethod
}

type UserAuthenticatorClaims struct {
	jwt.StandardClaims
	UserID uuid.UUID `json:"userId"`
}

func NewUserAuthenticator(storage *user.UserStorage, secret string) *UserAuthenticator {
	return &UserAuthenticator{
		Storage: storage,
		secret:  secret,
		defaultClaims: map[string]interface{}{
			"iss": "UserAuthenticator",
		},
		signingMethod: jwt.GetSigningMethod("HS256"),
	}
}

func (authenticator *UserAuthenticator) Issue(payload interface{}) (string, error) {
	uuid, isValid := payload.(uuid.UUID)

	if !isValid {
		return "", fmt.Errorf("Payload %v is not a valid uuid", payload)
	}
		
	tokenClaims := make(jwt.MapClaims)
	for key, value := range authenticator.defaultClaims {
		tokenClaims[key] = value
	}
	tokenClaims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	tokenClaims["iat"] = time.Now().Unix()
	tokenClaims["userId"] = uuid.String()

	token := jwt.NewWithClaims(authenticator.signingMethod, &tokenClaims)

	tokenString, err := token.SignedString([]byte(authenticator.secret))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (authenticator *UserAuthenticator) VerifyAndDecode(tokenString string) (interface{}, error) {
	tokenClaims := &UserAuthenticatorClaims{}

	token, err := jwt.ParseWithClaims(tokenString, tokenClaims, func(token *jwt.Token) (interface{}, error) {
		if token.Method == authenticator.signingMethod {
			return []byte(authenticator.secret), nil
		}

		return nil, errors.New("Incorrect signing method.")
	})

	if err == nil && token.Valid {
		_, userExistsErr := authenticator.Storage.Read(tokenClaims.UserID)

		if userExistsErr != nil {
			return uuid.UUID{}, fmt.Errorf("User with uuid %s does not exist", tokenClaims.UserID)
		}

		return tokenClaims.UserID, nil
	}

	return uuid.UUID{}, err
}

