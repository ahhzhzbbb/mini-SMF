package router

import (
	"errors"
	"mini-SMF/gateway/internal/registry"
	"sync"
)

type LoadBased struct {
	mu sync.Mutex
}

func NewLoadBased() *LoadBased {
	return &LoadBased{}
}

func (l *LoadBased) Next(reg *registry.Registry) (*registry.Instance, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if len(reg.Instances) == 0 {
		return nil, errors.New("no active instances available")
	}
	result := reg.Instances[0]
	for _, instance := range reg.Instances {
		if instance.ActiveRequests < result.ActiveRequests {
			result = instance
		}
	}
	return result, nil
}
