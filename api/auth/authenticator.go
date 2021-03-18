package auth

// Authenticator describes a structure capable of issuing and decoding JSON Web Tokens.
type Authenticator interface {
	IssueToken(payload interface{}) (string, error)      // IssueToken() should create an encoded token string.
	DecodeToken(tokenString string) (interface{}, error) // DecodeToken() should parse and validate the provided tokenString and return the encoded value.
}
