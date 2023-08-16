package helpers

import (
	"encoding/json"
	"io"
)

func BodyParser(body io.ReadCloser, v any) error {
	data, err := io.ReadAll(body)

	if err != nil {
		return err
	}

	err = json.Unmarshal(data, v)

	if err != nil {
		return err
	}

	return nil
}
