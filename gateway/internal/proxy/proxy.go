package proxy

import (
	"mini-SMF/gateway/internal/config"
	"mini-SMF/gateway/internal/middleware"
	"net/http"

	"github.com/rs/zerolog"
)

func NewProxy(
	config *config.Config,
	logger *zerolog.Logger,
) http.Handler {
	mux := http.NewServeMux()
	addRoute(
		config,
		mux,
	)
	var handler http.Handler
	handler = mux
	handler = middleware.LoggingMiddleware(logger, handler)
	handler = middleware.AuthMiddleware(handler)
	return handler
}
