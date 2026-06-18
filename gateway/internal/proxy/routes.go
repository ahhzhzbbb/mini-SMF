package proxy

import (
	"mini-SMF/gateway/internal/config"
	"mini-SMF/gateway/internal/handler"
	"mini-SMF/gateway/internal/registry"
	"mini-SMF/gateway/internal/router"
	"net/http"
)

func addRoute(
	config *config.Config,
	mux *http.ServeMux,
	registry *registry.Registry,
	loadBalancer router.LoadBalancer,
) {
	mux.Handle("GET /config", handler.HandlerGetProxyConfig(config))
	mux.Handle("GET /instances", handler.HandlerGetAllInstanceIp(registry))
	mux.Handle("POST /nsmf-pdusession/v1/sm-contexts", handler.HandlerPDUInstanceEstablishment(loadBalancer, "/nsmf-pdusession/v1/sm-contexts", registry))
	mux.Handle("GET /health", handler.HandlerHealthCheck("/heath", registry))
	mux.Handle("POST /register", handler.HandlerRegister(registry))
}
