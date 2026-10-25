package core_pgx

import (
	core_postgres "chatping-server/internal/core/adapter/postgres"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type pgxRows struct {
	pgx.Rows
}

type pgxRow struct {
	pgx.Row
}

func (r pgxRow) Scan(dest ...any) error {
	err := r.Row.Scan(dest...)
	if err != nil {
		return mapError(err)
	}

	return err
}

type pgxCommandTag struct {
	pgconn.CommandTag
}

func mapError(sourceErr error) error {
	const pgxViolatesForeignKeyCode = "23503"

	if errors.Is(sourceErr, pgx.ErrNoRows) {
		return core_postgres.ErrNoRows
	}

	if pgErr, ok := errors.AsType[*pgconn.PgError](sourceErr); ok {
		if pgErr.Code == pgxViolatesForeignKeyCode {
			return fmt.Errorf("%v: %w", sourceErr, core_postgres.ErrViolatesForeignKey)
		}
	}

	return fmt.Errorf("%v: %w", sourceErr, core_postgres.ErrUnknown)
}
