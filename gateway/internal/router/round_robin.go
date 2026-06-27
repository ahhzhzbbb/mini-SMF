package router

import (
	"errors"
	"mini-SMF/gateway/internal/registry"
	"sync"
	"sync/atomic"
)

type RoundRobin struct {
	current int64
	mu      sync.Mutex
}

func NewRoundRobin(current int64) *RoundRobin {
	return &RoundRobin{
		current: current,
	}
}

func (rr *RoundRobin) Next(reg *registry.Registry) (*registry.Instance, error) {
	n := len(reg.Instances)
	if n == 0 {
		return nil, errors.New("no active instances available")
	}

	idx := atomic.AddInt64(&rr.current, 1)

	return reg.Instances[idx%int64(n)], nil
}

// func RoundRobinLB(current *int, path string, registry *registry.Registry) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		if len(registry.Instances) == 0 {
// 			w.Write([]byte("There was not any active instance, try reload"))
// 		}

// 		*current = *current % (len(registry.Instances))
// 		instanceAddr := registry.Instances[*current].Address
// 		reqURL := "http://" + instanceAddr + path
// 		resp, err := http.Post(reqURL, "application/json", r.Body)
// 		if err != nil {
// 			http.Error(w, fmt.Sprintf("Proxy error: %v", err), http.StatusInternalServerError)
// 			return
// 		}
// 		defer resp.Body.Close()
// 		w.Write([]byte("RESPONSE: "))
// 		w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
// 		w.WriteHeader(resp.StatusCode)
// 		_, err = io.Copy(w, resp.Body)
// 		if err != nil {
// 			fmt.Printf("failed to get response: %v\n", err)
// 		}
// 		*current = (*current + 1) % (len(registry.Instances))
// 	})
// }
