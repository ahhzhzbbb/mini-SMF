package proxy

import (
	"mini-SMF/gateway/internal/config"
	"mini-SMF/gateway/internal/middleware"
	"mini-SMF/gateway/internal/registry"
	"net/http"

	"github.com/rs/zerolog"
)

func NewProxy(
	config *config.Config,
	logger *zerolog.Logger,
	registry *registry.Registry,
) http.Handler {
	mux := http.NewServeMux()
	addRoute(
		config,
		mux,
		registry,
	)
	var handler http.Handler
	handler = mux
	handler = middleware.LoggingMiddleware(logger, handler)
	handler = middleware.AuthMiddleware(handler)
	return handler
}
