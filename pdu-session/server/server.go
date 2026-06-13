package server

import (
	"mini-SMF/pdu-session/config"
	"mini-SMF/pdu-session/server/middleware"
	"net/http"

	"github.com/rs/zerolog"
)

func NewServer(
	config *config.Config,
	logger *zerolog.Logger,
) http.Handler {
	mux := http.NewServeMux()

	addRoutes(
		config,
		mux,
	)
	var handler http.Handler = mux
	handler = middleware.LoggingMiddleware(logger, handler)
	return handler
}
