package router

import (
	"fmt"
	"io"
	"mini-SMF/gateway/internal/registry"
	"net/http"
)

func RoundRobinLB(current *int, registry *registry.Registry) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(registry.Instances) == 0 {
			w.Write([]byte("There was not any active instance, try reload"))
		}

		*current = *current % (len(registry.Instances))
		instanceAddr := registry.Instances[*current].Address
		reqURL := "http://" + instanceAddr + ":8081/nsmf-pdusession/v1/sm-contexts"
		resp, err := http.Post(reqURL, "application/json", r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Proxy error: %v", err), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		w.Write([]byte("RESPONSE: "))
		w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
		w.WriteHeader(resp.StatusCode)
		_, err = io.Copy(w, resp.Body)
		if err != nil {
			fmt.Printf("failed to get response: %v\n", err)
		}
		*current = (*current + 1) % (len(registry.Instances))
	})
}
