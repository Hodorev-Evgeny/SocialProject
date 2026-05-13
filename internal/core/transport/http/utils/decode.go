package core_http_utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator"
)

var v = validator.New()

type customValidator interface {
	Validate() error
}

func DecodeJSON(data any, r *http.Request) error {
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		return fmt.Errorf("error decoding JSON: %w", err)
	}

	// проверка, есть ли кастомные валидации
	value, ok := data.(customValidator)
	if ok {
		return value.Validate()
	}

	return v.Struct(data)
}
