package option

import (
	"bytes"
	"encoding/json"
)

// Option wraps an optional value. A nil pointer represents None.
type Option[T any] struct {
	value *T
}

// Some constructs an Option containing the provided value.
func Some[T any](v T) Option[T] {
	return Option[T]{value: &v}
}

// None constructs an empty Option.
func None[T any]() Option[T] {
	return Option[T]{}
}

// FromPtr builds an Option from a pointer; nil becomes None.
func FromPtr[T any](ptr *T) Option[T] {
	return Option[T]{value: ptr}
}

// IsSome reports whether the Option stores a value.
func (o Option[T]) IsSome() bool {
	return o.value != nil
}

// IsNone reports whether the Option is empty.
func (o Option[T]) IsNone() bool {
	return o.value == nil
}

// IsZero allows json tags such as omitzero/omitempty to drop None.
func (o Option[T]) IsZero() bool {
	return o.value == nil
}

// Ptr returns the underlying pointer; callers must check IsSome first.
func (o Option[T]) Ptr() *T {
	return o.value
}

// Value returns the stored value and a flag indicating presence.
func (o Option[T]) Value() (T, bool) {
	if o.value == nil {
		var zero T
		return zero, false
	}
	return *o.value, true
}

// ValueOr returns the stored value when present, or fallback otherwise.
func (o Option[T]) ValueOr(fallback T) T {
	if o.value == nil {
		return fallback
	}
	return *o.value
}

// Map transforms the stored value when present.
func Map[T, U any](o Option[T], f func(T) U) Option[U] {
	if o.value == nil {
		return None[U]()
	}
	return Some(f(*o.value))
}

// FlatMap chains optional computations.
func FlatMap[T, U any](o Option[T], f func(T) Option[U]) Option[U] {
	if o.value == nil {
		return None[U]()
	}
	return f(*o.value)
}

// MarshalJSON serializes None as null and Some as the inner value.
func (o Option[T]) MarshalJSON() ([]byte, error) {
	if o.value == nil {
		return []byte("null"), nil
	}
	return json.Marshal(*o.value)
}

// UnmarshalJSON interprets null as None and otherwise stores the parsed value.
func (o *Option[T]) UnmarshalJSON(data []byte) error {
	if bytes.Equal(bytes.TrimSpace(data), []byte("null")) {
		o.value = nil
		return nil
	}
	var v T
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	o.value = &v
	return nil
}
