package errors

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type publicError struct {
	err error
}

func (e publicError) Error() string {
	return e.err.Error()
}

func newPublicError(c codes.Code, format string, a ...any) error {
	return publicError{status.Errorf(c, format, a...)}
}

var (
	BadRequest                = newPublicError(codes.InvalidArgument, "bad request")
	MissingMetadata           = newPublicError(codes.InvalidArgument, "missing metadata")
	InvalidToken              = newPublicError(codes.Unauthenticated, "invalid token")
	SessionExpire             = newPublicError(codes.Unauthenticated, "session expired")
	RegistrationSessionExpire = newPublicError(codes.Unauthenticated, "registration session expire")
	Unauthorized              = newPublicError(codes.PermissionDenied, "permission denied")
	NotFound                  = newPublicError(codes.NotFound, "requested item was not found")
	DBError                   = newPublicError(codes.Internal, "database error")
	InternalError             = newPublicError(codes.Internal, "internal server error")
)

var publicErr = &publicError{}

func IsPublicError(err error) bool {
	return errors.As(err, publicErr)
}
