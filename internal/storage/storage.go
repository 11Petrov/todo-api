package storage

import "errors"

var (
	ErrDBConnection = errors.New("failed to connect to database")
)
