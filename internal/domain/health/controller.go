// Package health for checking server health
package health

import (
	"context"
	healthv1 "fbt/backend/gen/proto/go/health/v1"
	"fbt/backend/gen/proto/go/health/v1/healthv1connect"
	"fbt/backend/internal/util"
	"net/http"

	"connectrpc.com/connect"
)

func RegisterService(mux *http.ServeMux, d *util.Dependency, opts ...connect.HandlerOption) *http.ServeMux {
	mux.Handle(NewServiceHandler(d, opts...))

	return mux
}

type con struct {
	d *util.Dependency
}

func NewServiceHandler(d *util.Dependency, opts ...connect.HandlerOption) (string, http.Handler) {
	return healthv1connect.NewHealthServiceHandler(&con{d: d}, opts...)
}

func (s *con) Status(ctx context.Context, req *healthv1.HealthServiceStatusRequest) (*healthv1.HealthServiceStatusResponse, error) {
	dbErr := s.d.DB.Ping(ctx)

	return &healthv1.HealthServiceStatusResponse{
		IsServerUp:   true,
		IsDatabaseUp: dbErr == nil,
	}, nil
}
