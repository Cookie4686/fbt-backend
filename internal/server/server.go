package server

import (
	"fbt/backend/internal/domain/auth"
	"fbt/backend/internal/domain/auth/service"
	"fbt/backend/internal/domain/bookkeeping"
	"fbt/backend/internal/interceptor"
	"fbt/backend/internal/util"
	"fmt"
	"net/http"
	"time"

	"connectrpc.com/connect"
)

const readHeaderTimeout = 2 * time.Second

func NewServer(d *util.Dependency) *http.Server {
	mux := http.NewServeMux()

	service := service.NewService(d)

	i := interceptor.NewInterceptorProvider(d, service)

	opts := connect.WithInterceptors(
		i.Logging(),
		i.Auth(),
		i.Validator(),
	)

	auth.RegisterService(mux, d, opts)
	bookkeeping.RegisterService(mux, d, opts)

	p := new(http.Protocols)
	p.SetHTTP1(true)

	// Use h2c so we can serve HTTP/2 without TLS.
	p.SetUnencryptedHTTP2(true)

	server := http.Server{
		Addr:              fmt.Sprintf(":%d", d.CFG.API.PORT),
		Handler:           mux,
		Protocols:         p,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	return &server
}
