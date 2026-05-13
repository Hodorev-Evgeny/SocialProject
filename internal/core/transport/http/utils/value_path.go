package core_http_utils

import (
	"fmt"
	"net/http"
	"strconv"

	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
)

func GetValuePathInt(r *http.Request, key string) (int, error) {
	value := r.PathValue(key)
	if value == "" {
		return 0, fmt.Errorf("path %q is required: %w", key, core_errors.ErrorValidation)
	}

	valueInt, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid value for path %s: %w", key, core_errors.ErrorValidation)
	}

	return valueInt, nil
}
