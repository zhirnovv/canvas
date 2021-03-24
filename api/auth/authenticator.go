package auth

// Authenticator describes an instance capable of issuing and decoding JSON Web Tokens.
type Authenticator interface {
	Issue(payload interface{}) (string, error)      // Issue() should create an encoded token string.
	VerifyAndDecode(tokenString string) (interface{}, error) // VerifyAndDecode() should parse and validate the provided tokenString and return the encoded value.
}
