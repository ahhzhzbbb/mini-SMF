package proxy

import (
	"mini-SMF/gateway/internal/config"
	"mini-SMF/gateway/internal/handler"
	"mini-SMF/gateway/internal/registry"
	"net/http"
)

func addRoute(
	config *config.Config,
	mux *http.ServeMux,
	registry *registry.Registry,
) {
	mux.Handle("/config", handler.HandlerGetProxyConfig(config))
	mux.Handle("/instances", handler.HandlerGetAllInstanceIp(registry))
}
