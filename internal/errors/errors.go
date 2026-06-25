package errors

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	BadRequest                = status.Errorf(codes.InvalidArgument, "bad request")
	MissingMetadata           = status.Errorf(codes.InvalidArgument, "missing metadata")
	InvalidToken              = status.Errorf(codes.Unauthenticated, "invalid token")
	SessionExpire             = status.Errorf(codes.Unauthenticated, "session expired")
	RegistrationSessionExpire = status.Errorf(codes.Unauthenticated, "registration session expire")
	Unauthorized              = status.Errorf(codes.PermissionDenied, "permission denied")
	NotFound                  = status.Errorf(codes.NotFound, "requested item was not found")
	DBError                   = status.Errorf(codes.Internal, "database error")
)
