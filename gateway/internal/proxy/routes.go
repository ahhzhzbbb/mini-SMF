package proxy

import (
	"mini-SMF/gateway/internal/config"
	"mini-SMF/gateway/internal/handler"
	"net/http"
)

func addRoute(
	config *config.Config,
	mux *http.ServeMux,
) {
	mux.Handle("/config", handler.HandlerGetProxyConfig(config))
}
