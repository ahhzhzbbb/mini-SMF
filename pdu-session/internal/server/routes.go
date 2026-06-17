package server

import (
	"mini-SMF/pdu-session/internal/config"
	"mini-SMF/pdu-session/internal/handler"
	"net/http"
)

func addRoutes(
	config *config.Config,
	mux *http.ServeMux,
) {
	mux.Handle("GET /", handler.HandlerGetMessage())
	mux.Handle("GET /config", handler.HandlerGetServerConfigInfo(config))
	mux.Handle("POST /nsmf-pdusession/v1/sm-contexts", handler.HandlerPDUSessionEstablishment(config))
	mux.Handle("GET /health", handler.HandlerGetHeath(config))
}
