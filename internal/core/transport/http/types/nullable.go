package core_http_types

import (
	"encoding/json"
	"fmt"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
)

type Nullable[T any] struct {
	core_domain.Nullable[T]
}

func (n *Nullable[T]) UnmarshalJSON(data []byte) error {
	n.Set = true

	if string(data) == "null" {
		n.Value = nil

		return nil
	}

	var val T
	if err := json.Unmarshal(data, &val); err != nil {
		return fmt.Errorf("unable to unmarshal Nullable[%T]: %w", val, err)
	}

	n.Value = &val
	return nil
}

func (n *Nullable[T]) ToDomain() core_domain.Nullable[T] {
	return core_domain.Nullable[T]{
		Value: n.Value,
		Set:   n.Set,
	}
}
