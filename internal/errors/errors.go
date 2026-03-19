package errors

import "errors"

var (
	NotFound      = errors.New("requested item was not found")
	SessionExpire = errors.New("session expired")
	DBError       = errors.New("database error")
)
