package lorm

import "github.com/pkg/errors"

var (
	ErrNil = errors.New("nil")
	ErrContainEmpty = errors.New("slice empty")
)

