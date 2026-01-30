package middlewares

import (
	"github.com/anzhy11/go-e-commerce/internal/config"
	"github.com/anzhy11/go-e-commerce/internal/server"
)

type Middlewares struct {
	cfg *config.Config
	srv *server.Server
}

func New(srv *server.Server, cfg *config.Config) *Middlewares {
	md := &Middlewares{
		cfg: cfg,
		srv: srv,
	}

	md.Authorization()

	return md
}
