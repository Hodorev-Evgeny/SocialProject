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
		return 0, fmt.Errorf(`path "%s" is required`, key)
	}

	valueInt, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid value for path %s: %e", key, core_errors.ErrorValidation)
	}

	return valueInt, nil
}
