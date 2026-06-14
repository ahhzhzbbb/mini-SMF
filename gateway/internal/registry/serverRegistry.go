package registry

import (
	"fmt"
	"net"
	"sync"
	"time"
)

type Instance struct {
	ID            string
	ServiceName   string
	Address       string
	LastHeartbeat time.Time
}

func NewInstance(id, serviceName, address string) *Instance {
	return &Instance{
		ID:          id,
		ServiceName: serviceName,
		Address:     address,
	}
}

type Registry struct {
	mu       sync.RWMutex
	Services map[string]([]*Instance)
}

func NewRegistry() *Registry {
	return &Registry{
		Services: make(map[string]([]*Instance)),
	}
}

func (r *Registry) Load(serviceName string) error {
	ips, err := net.LookupIP(serviceName)
	if err != nil {
		return err
	}

	count := 1
	for _, ip := range ips {
		newIntance := NewInstance(fmt.Sprintf("%s-%d", serviceName, count), serviceName, ip.String())
		r.Services[serviceName] = append(r.Services[serviceName], newIntance)
		count++
	}
	return nil
}

func (r *Registry) Register(instance *Instance)

func (r *Registry) Deregister(id string)

func (r *Registry) Heartbeat(id string)

func (r *Registry) Get(id string) (*Instance, bool)

func (r *Registry) List() []*Instance
