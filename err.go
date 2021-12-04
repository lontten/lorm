package lorm

import "github.com/pkg/errors"

var (
	ErrNil          = errors.New("nil")
	ErrContainEmpty = errors.New("slice empty")
	ErrNoPkOrUnique = errors.New(" ERROR: there is no unique or exclusion constraint matching the ON CONFLICT specification (SQLSTATE 42P10) ")
)
