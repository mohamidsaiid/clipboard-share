package jsonParser 

import (
	"net/http"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

func ReadJSON(r *http.Response, dst any) error {
	// To limit the size of the recived response
	maxBytes := 1_048_576
	// Decode the req body into to targeted destaniation
	dec := json.NewDecoder(r.Body)
	// To disallow and unknown fileds to be in the response which cannot be mapped to the target destination
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)

	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		// Using errors.As() to specify the type of the error
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at charcter %d)", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at charcter %d)", unmarshalTypeError.Offset)

		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		case strings.HasPrefix(err.Error(), "json:unknown field"):
			fieldName := strings.TrimPrefix(err.Error(), "json:unknown field")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", maxBytes)

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}