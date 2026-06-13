package server

import (
	"mini-SMF/pdu-session/internal/config"
	"mini-SMF/pdu-session/internal/server/handler"
	"net/http"
)

func addRoutes(
	config *config.Config,
	mux *http.ServeMux,
) {
	mux.Handle("/", handler.HandlerGetMessage())
	mux.Handle("/config", handler.HandlerGetServerConfigInfo(config))
}
