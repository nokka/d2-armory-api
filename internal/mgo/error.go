package mgo

import (
	"fmt"
	"net"

	"github.com/nokka/d2-armory-api/internal/domain"
	"go.mongodb.org/mongo-driver/mongo"
)

// Error is a package specific error type that lets us define immutable errors.
type Error string

func (e Error) Error() string {
	return string(e)
}

func mongoErr(err error) error {
	// If the error is a timeout or temporary error, return temporary error.
	if err, ok := err.(net.Error); ok {
		if err.Temporary() || err.Timeout() {
			return fmt.Errorf("temporary error while performing query: %w", domain.ErrTemporary)
		}
	}

	// if the error is a dial error, return temporary.
	if err, ok := err.(*net.OpError); ok {
		if err.Op == "dial" {
			return fmt.Errorf("dial error: %w", domain.ErrTemporary)
		}
	}

	switch err {
	case mongo.ErrNoDocuments,
		mongo.ErrNilDocument:
		return fmt.Errorf("%w", domain.ErrNotFound)
	case mongo.ErrClientDisconnected,
		mongo.ErrUnacknowledgedWrite:
		return fmt.Errorf("%w", domain.ErrTemporary)
	}

	return fmt.Errorf("unspecified error: %s, %w", err, domain.ErrInternal)
}
