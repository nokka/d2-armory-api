package mongodb

import (
	"errors"
	"fmt"
	"net"

	"github.com/nokka/d2-armory-api/internal/domain"
	"gopkg.in/mgo.v2"
)

// Error is a package specific error type that lets us define immutable errors.
type Error string

func (e Error) Error() string {
	return string(e)
}

func mongoErr(err error) error {
	if err == mgo.ErrNotFound {
		return fmt.Errorf("no result returned: %w", domain.ErrNotFound)
	}

	if err, ok := err.(net.Error); ok {
		if err.Temporary() || err.Timeout() {
			return fmt.Errorf("temporary error while performing query: %w", domain.ErrTemporary)
		}
	}

	if err, ok := err.(*net.OpError); ok {
		if err.Op == "dial" {
			return fmt.Errorf("dial error: %w", domain.ErrTemporary)
		}
	}

	codeErr := func(code int) error {
		switch code {
		case 6, // HostUnreachable
			7,    // HostNotFound
			24,   // LockTimeout
			46,   // LockBusy
			50,   // ExceededTimeLimit
			89,   // NetworkTimeout
			91,   // ShutdownInProgress
			94,   // NotYetInitialized
			101,  // OutdatedClient
			107,  // LockFailed
			109,  // ConfigurationInProgress
			146,  // ExceededMemoryLimit
			164,  // InitialSyncActive
			202,  // NetworkInterfaceExceededTimeLimit
			208,  // TooManyLocks
			9001: // SocketException
			return domain.ErrTemporary

		case 84: // DuplicateKeyValue
			return domain.ErrConflict
		}
		return nil
	}

	switch err := err.(type) {
	case *mgo.BulkError:
		// Take first error
		for _, e := range err.Cases() {
			return mongoErr(e.Err)
		}
		return errors.New(err.Error())
	case *mgo.QueryError:
	case *mgo.LastError:
		return fmt.Errorf("unspecified error: %w", codeErr(err.Code))
	}

	return err
}
