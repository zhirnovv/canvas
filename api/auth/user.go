package auth

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

// UserAuthenticator contains all necessary parameters for signing and parsing JSON Web Tokens.
type UserAuthenticator struct {
	secret        string // Signing key. Under no circumstances should this value be exportable.
	defaultClaims jwt.MapClaims
	signingMethod jwt.SigningMethod
}

type UserAuthenticatorClaims struct {
	jwt.StandardClaims
	UserID uuid.UUID `json:"userId"`
}

func NewUserAuthenticator(secret string) *UserAuthenticator {
	return &UserAuthenticator{
		secret: secret,
		defaultClaims: map[string]interface{}{
			"iss": "UserAuthenticator",
		},
		signingMethod: jwt.GetSigningMethod("HS256"),
	}
}

func (authenticator *UserAuthenticator) IssueToken(userId uuid.UUID) (string, error) {
	tokenClaims := make(jwt.MapClaims)
	for key, value := range authenticator.defaultClaims {
		tokenClaims[key] = value
	}
	tokenClaims["exp"] = time.Now().Add(time.Second * 5).Unix()
	tokenClaims["iat"] = time.Now().Unix()
	tokenClaims["userId"] = userId.String()

	token := jwt.NewWithClaims(authenticator.signingMethod, &tokenClaims)

	tokenString, err := token.SignedString([]byte(authenticator.secret))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (authenticator *UserAuthenticator) DecodeToken(tokenString string) (uuid.UUID, error) {
	tokenClaims := &UserAuthenticatorClaims{}
	
	token, err := jwt.ParseWithClaims(tokenString, tokenClaims, func(token *jwt.Token) (interface{}, error) {
		if token.Method == authenticator.signingMethod {
			return []byte(authenticator.secret), nil
		}

		return nil, errors.New("Incorrect signing method.")
	})

	if err == nil && token.Valid {
		return tokenClaims.UserID, nil
	}

	return uuid.UUID{}, err
}
