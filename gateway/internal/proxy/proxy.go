package proxy

import (
	"mini-SMF/gateway/internal/config"
	"mini-SMF/gateway/internal/middleware"
	"mini-SMF/gateway/internal/registry"
	"mini-SMF/gateway/internal/router"
	"net/http"

	"github.com/rs/zerolog"
)

func NewProxy(
	config *config.Config,
	logger *zerolog.Logger,
	registry *registry.Registry,
	loadBalancer router.LoadBalancer,
) http.Handler {
	mux := http.NewServeMux()
	addRoute(
		config,
		mux,
		registry,
		loadBalancer,
	)
	var handler http.Handler
	handler = mux
	handler = middleware.LoggingMiddleware(logger, handler)
	handler = middleware.AuthMiddleware(handler)
	handler = middleware.CheckingTimeoutMiddleware(config.TimeoutDuration, handler)
	return handler
}
