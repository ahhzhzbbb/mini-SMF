package server

import (
	"errors"
	"strconv"
	"sync"
)

type IPAllocator struct {
	mu   sync.Mutex
	pool []string
	used map[string]bool
}

func NewIPAllocator() *IPAllocator {
	newPool := make([]string, 253)
	for i := 2; i <= 254; i++ {
		newPool[i-2] = "10.0.0." + strconv.Itoa(i)
	}
	return &IPAllocator{
		pool: newPool,
		used: make(map[string]bool),
	}
}

func (a *IPAllocator) Allocate() (string, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	for _, ip := range a.pool {
		if !a.used[ip] {
			a.used[ip] = true
			return ip, nil
		}
	}
	return "", errors.New("no IP available")
}
