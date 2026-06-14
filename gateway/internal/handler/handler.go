package handler

import (
	"fmt"
	"mini-SMF/gateway/internal/config"
	"mini-SMF/gateway/internal/registry"
	"net/http"
)

func HandlerGetProxyConfig(config *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Host: %s\nPort: %s\n", config.Host, config.Port)
	})
}

func HandlerGetAllInstanceIp(reg *registry.Registry) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		instanceAdressSlice := reg.GetAllInstances()
		for _, i := range instanceAdressSlice {
			i += "\n"
			w.Write([]byte(i))
		}
	})
}
