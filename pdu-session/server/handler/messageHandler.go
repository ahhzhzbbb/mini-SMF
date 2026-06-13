package handler

import (
	"fmt"
	"mini-SMF/pdu-session/config"
	"net/http"
)

func HandlerGetMessage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!!!"))
	})
}

func HandlerGetServerConfigInfo(config *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Host: %s\nPort: %s\n", config.Host, config.Port)
	})
}
