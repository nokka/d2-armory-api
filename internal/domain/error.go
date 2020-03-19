package domain

// Error is a package specific error type that lets us define immutable errors.
type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	// ErrRequest is returned when the request was invalid.
	ErrRequest = Error("invalid request data")

	// ErrNotFound is returned when a resource can't be find.
	ErrNotFound = Error("resource was not found")

	// ErrInvalidArgument is returned when one or more arguments are invalid.
	ErrInvalidArgument = Error("invalid argument")

	// ErrUnavailable is returned when the service is unavailable.
	ErrUnavailable = Error("service unavailable")

	// ErrTemporary is returned when the service is temporarily unavailable.
	ErrTemporary = Error("service is temporary unavailable")

	// ErrConflict is returned when there's a conflict with a resource.
	ErrConflict = Error("conflict error")
)
