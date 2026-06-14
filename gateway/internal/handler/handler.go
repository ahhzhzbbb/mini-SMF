package handler

import (
	"fmt"
	"mini-SMF/gateway/internal/config"
	"net/http"
)

func HandlerGetProxyConfig(config *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Host: %s\nPort: %s\n", config.Host, config.Port)
	})
}
