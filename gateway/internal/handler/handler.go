package handler

import (
	"fmt"
	"io"
	"mini-SMF/gateway/internal/config"
	"mini-SMF/gateway/internal/registry"
	"mini-SMF/gateway/internal/router"
	"net"
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

func HandlerPDUInstanceEstablishment(lb router.LoadBalancer, path string, reg *registry.Registry) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(reg.Instances) == 0 {
			http.Error(w, "There was not any active instance, try reload", http.StatusServiceUnavailable)
			return
		}

		instance, err := lb.Next(reg.Instances)
		if err != nil {
			w.Write([]byte(err.Error()))
		}

		instanceAddress := net.JoinHostPort(instance.IpAddr, instance.Port)
		reqURL := "http://" + instanceAddress + path

		req, err := http.NewRequestWithContext(r.Context(), "POST", reqURL, r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
		w.WriteHeader(resp.StatusCode)

		w.Write([]byte("RESPONSE: "))
		_, err = io.Copy(w, resp.Body)
		if err != nil {
			fmt.Printf("failed to get response: %v\n", err)
		}
	})
}

func HandlerHealthCheck(path string, reg *registry.Registry) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(reg.Instances) == 0 {
			http.Error(w, "There was not any active instance, try reload", http.StatusServiceUnavailable)
			return
		}
		if err := reg.HealthCheck(path); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}
