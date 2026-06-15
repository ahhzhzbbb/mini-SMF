package registry

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

type Instance struct {
	ID             string
	ServiceName    string
	IpAddr         string
	Port           string
	Weight         float32
	ActiveRequests int
	LastHeartbeat  time.Time
}

func NewInstance(id, serviceName, address, port string) *Instance {
	return &Instance{
		ID:          id,
		ServiceName: serviceName,
		IpAddr:      address,
		Port:        port,
	}
}

type Registry struct {
	mu        sync.RWMutex
	Instances []*Instance
}

func NewRegistry() *Registry {
	return &Registry{
		Instances: make([]*Instance, 0),
	}
}

func (r *Registry) Load(serviceName string) error {
	ips, err := net.LookupIP(serviceName)
	if err != nil {
		return err
	}
	port := os.Getenv("PDU_PORT")

	count := 1
	for _, ip := range ips {
		newIntance := NewInstance(fmt.Sprintf("%s-%d", serviceName, count), serviceName, ip.String(), port)
		r.Instances = append(r.Instances, newIntance)
		count++
	}
	return nil
}

func (r *Registry) GetAllInstances() []string {
	fmt.Println(len(r.Instances))
	var res []string
	for _, i := range r.Instances {
		res = append(res, i.IpAddr)
	}
	return res
}

// func (r *Registry) Register(instance *Instance)

// func (r *Registry) Deregister(id string)

// func (r *Registry) Heartbeat(id string)

// func (r *Registry) Get(id string) (*Instance, bool)

// func (r *Registry) List() []*Instance
