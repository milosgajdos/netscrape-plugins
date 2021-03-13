package dgraph

import "errors"

var (
	// ErrUnknownOp is returned when attempting to use an unkown operation on store.
	ErrUnknownOp = errors.New("ErrUnknownOp")
)
