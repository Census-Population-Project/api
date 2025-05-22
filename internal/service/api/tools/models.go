package tools

import "encoding/json"

type Optional[T any] struct {
	Value   *T
	Defined bool
}

func (t *Optional[T]) UnmarshalJSON(data []byte) error {
	t.Defined = true
	return json.Unmarshal(data, &t.Value)
}

func (t *Optional[T]) IsNullDefined() bool {
	return t.Defined && t.Value == nil
}

func (t *Optional[T]) HasValue() bool {
	return t.Defined && t.Value != nil
}

type TimeFormat string

const (
	TimeWithTimeZone    TimeFormat = "15:04:05-07:00"
	TimeWithoutTimeZone TimeFormat = "15:04:05"
)
