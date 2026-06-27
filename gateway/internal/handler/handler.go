package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mini-SMF/gateway/internal/config"
	"mini-SMF/gateway/internal/registry"
	"mini-SMF/gateway/internal/router"
	"net"
	"net/http"
	"sync/atomic"
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
		// Gateway trả 503: {"cause":
		// "NO_BACKEND_AVAILABLE"}.
		if reg.IsEmpty() {
			http.Error(w, "NO_BACKEND_AVAILABLE", http.StatusServiceUnavailable)
			return
		}

		instance, err := lb.Next(reg.Instances)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		instanceAddress := net.JoinHostPort(instance.IpAddr, instance.Port)
		reqURL := "http://" + instanceAddress + path

		req, err := http.NewRequestWithContext(r.Context(), "POST", reqURL, r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tr := &http.Transport{
			Protocols: new(http.Protocols),
		}

		tr.Protocols.SetUnencryptedHTTP2(true)

		client := &http.Client{
			Transport: tr,
		}

		resp, err := client.Do(req)
		if err != nil {
			if errors.Is(r.Context().Err(), context.DeadlineExceeded) {
				newTimeoutCount := atomic.AddInt64(&instance.TimeoutRequests, -1)
				fmt.Printf("timeout number: %d\n", newTimeoutCount)
				if newTimeoutCount <= 0 {
					if err := reg.RemoveInstance(instance); err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
				}
			}
			http.Error(w, err.Error(), http.StatusGatewayTimeout)
			return
		}
		defer resp.Body.Close()

		w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
		w.WriteHeader(resp.StatusCode)

		w.Write([]byte("RESPONSE: \n"))
		_, err = io.Copy(w, resp.Body)
		if err != nil {
			fmt.Printf("failed to get response: %v\n", err)
		}
		atomic.StoreInt64(&instance.TimeoutRequests, 3)
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

func HandlerRegister(reg *registry.Registry) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type registerRequest struct {
			ServiceName string `json:"service_name"`
			Ip          string `json:"ip"`
			Port        string `json:"port"`
		}

		var req registerRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := reg.Register(req.ServiceName, req.Ip, req.Port); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	})
}
