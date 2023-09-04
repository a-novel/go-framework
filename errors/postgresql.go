package errors

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/uptrace/bun/driver/pgdriver"
	"strings"
)

var (
	// ErrConstraintViolation is thrown when a psql query violates any constraint.
	ErrConstraintViolation = fmt.Errorf("%w: record does not satisfy some of the column constraints", ErrInvalidEntity)
	// ErrUniqConstraintViolation is thrown when a psql query violates a unique constraint.
	ErrUniqConstraintViolation = fmt.Errorf("%w: some unique columns have duplicates", ErrConstraintViolation)
	// ErrMissingRelation is thrown when passing a bad reference in a column.
	ErrMissingRelation = fmt.Errorf("%w: a required relation is missing", ErrInvalidEntity)
	// ErrNotFound is thrown when a psql query returns no result.
	ErrNotFound = fmt.Errorf("could not find any record matching the request")
)

// HandlePGError extends pg library typed errors. Only a few errors are typed to be targeted with errors.Is, and some
// pretty common errors aren't. This handler parses postgres errors in a more test-friendly way.
func HandlePGError(err error) error {
	if err == nil {
		return nil
	}

	if err == sql.ErrNoRows {
		return ErrNotFound
	}

	pgErr, ok := err.(pgdriver.Error)
	if ok {
		if pgErr.IntegrityViolation() {
			// This error has a special treatment because, in most case, it is not checked upfront by the service.
			// Other constraint violation should be prevented by appropriate type checking in the service layer.
			// https://www.postgresql.org/docs/current/errcodes-appendix.html
			if strings.Contains(err.Error(), "SQLSTATE=23505") {
				return ErrUniqConstraintViolation
			}

			return ErrConstraintViolation
		} else if pgErr.StatementTimeout() {
			return errors.Join(ErrTimeout, err)
		}
	}

	// Catch custom agora errors.
	if strings.Contains(err.Error(), "AGORA=MISSINGRELATION") {
		return ErrMissingRelation
	}

	return err
}

func ForceRowsUpdate(res sql.Result) error {
	rows, err := res.RowsAffected()

	if err != nil {
		return fmt.Errorf("failed to check rows affected by the operation: %w", err)
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}
