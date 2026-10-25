package core_postgres

import "errors"

var (
	ErrNoRows             = errors.New("not rows")
	ErrViolatesForeignKey = errors.New("violates foreign key")
	ErrUnknown            = errors.New("unknown")
)
