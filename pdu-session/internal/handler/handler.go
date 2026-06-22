package handler

import (
	"fmt"
	"mini-SMF/pdu-session/internal/config"
	"net/http"
	"os"
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

func HandlerPDUSessionEstablishment(config *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		instanceID, err := os.Hostname()
		if err != nil {
			instanceID = "unknown-node"
		}
		fmt.Fprintf(w, "Hello Client, I am instance: %s, (Listening on %s:%s)",
			instanceID, config.Host, config.Port)
	})
}

func HandlerGetHeath(config *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		instanceID, err := os.Hostname()
		if err != nil {
			instanceID = "unknown-node"
		}
		fmt.Fprintf(w, "Hello Client, I am instance: %s, (Listening on %s:%s)",
			instanceID, config.Host, config.Port)
	})
}
