package errors

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	NotFound                  = errors.New("requested item was not found")
	SessionExpire             = errors.New("session expired")
	RegistrationSessionExpire = errors.New("registration session expire")
	DBError                   = errors.New("database error")
	BadRequest                = errors.New("bad request")
	MissingMetadata           = status.Errorf(codes.InvalidArgument, "missing metadata")
	Unauthorized              = status.Errorf(codes.PermissionDenied, "permission denied")
	InvalidToken              = status.Errorf(codes.Unauthenticated, "invalid token")
)
