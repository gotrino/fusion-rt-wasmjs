package binding

import (
	"net/http"
)

const (
	decoderError = -1
	encoderError = -2
)

// A Repository which represents CRUD (create read update delete) operations on an Entity based resource set.
// If any entity has a certain kind of ID, the repository implementation must unmarshal it from a string
// to support Load and Delete.
type Repository[T any] interface {
	List() ([]T, error)
	Load(id string) (T, error)
	// Delete removes the entity. It is no error, if an already deleted entry is removed again.
	Delete(id string) error
	// Save updates or creates the Entity.
	Save(t T) error
}

// ResourceRepository represents an aggregate which may or may not have an id.
type ResourceRepository[T any] interface {
	Load() (T, error)
	Save(t T) error
	Delete() error
}

type httpError struct {
	status int
	cause  error
}

func (e httpError) Error() string {
	return "http-error"
}

func (e httpError) Forbidden() bool {
	return e.status == http.StatusForbidden
}

func (e httpError) Unauthenticated() bool {
	return e.status == http.StatusUnauthorized
}

func (e httpError) InternalServerError() bool {
	return e.status >= http.StatusInternalServerError && e.status <= http.StatusVariantAlsoNegotiates
}

func (e httpError) ProtocolError() bool {
	return e.status == decoderError || e.status == encoderError
}

func (e httpError) Unwrap() error {
	return e.cause
}
