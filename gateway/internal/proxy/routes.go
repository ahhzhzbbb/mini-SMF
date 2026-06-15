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
	current *int,
) {
	mux.Handle("/config", handler.HandlerGetProxyConfig(config))
	mux.Handle("/instances", handler.HandlerGetAllInstanceIp(registry))
	mux.Handle("POST /nsmf-pdusession/v1/sm-contexts", handler.HandlerPDUInstanceEstablishment(router.RoundRobinLB(current, registry)))
}
