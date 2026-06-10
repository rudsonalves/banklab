package sharedhttp

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

// DecodeJSON decodes exactly one JSON value and rejects unknown fields.
func DecodeJSON(body io.Reader, target any) error {
	decoder := json.NewDecoder(body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(target); err != nil {
		return err
	}

	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		if err != nil {
			return err
		}

		return fmt.Errorf("request body must contain a single JSON value")
	}

	return nil
}
