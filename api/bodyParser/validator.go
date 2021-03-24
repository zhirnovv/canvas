package bodyParser

import (
	"github.com/go-playground/validator"
	"github.com/mitchellh/mapstructure"
)

var validate = validator.New()

// DecodeAndValidate() decodes a parsed json body to a schema and validates it.
// schema has to be a pointer to a struct.
func DecodeAndValidate(body map[string]interface{}, schema interface{}) error {
	formDecodeErr := mapstructure.Decode(body, &schema)

	if formDecodeErr != nil {
		return formDecodeErr
	}

	validationErr := validate.Struct(schema)

	if validationErr != nil {
		return validationErr
	}

	return nil
}
