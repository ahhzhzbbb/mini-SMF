package registry

import (
	"errors"
	"fmt"
	"net"
	"net/http"
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
	mu                  sync.RWMutex
	Instances           []*Instance
	CycleHeathCheckTime int          //second
	ClientHealthCheck   *http.Client //client representative to send health check request:V
}

func NewRegistry() *Registry {
	return &Registry{
		Instances:           make([]*Instance, 0),
		CycleHeathCheckTime: 1,
		ClientHealthCheck: &http.Client{
			Timeout: 1 * time.Millisecond,
		},
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

func (r *Registry) Remove(index int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if index < 0 || index >= len(r.Instances) {
		return errors.New("out of range")
	}

	r.Instances = append(r.Instances[:index], r.Instances[index+1:]...)

	return nil
}

func (r *Registry) HealthCheck(path string) error {
	for idx, instance := range r.Instances {
		address := net.JoinHostPort(instance.IpAddr, instance.Port)
		reqURL := "http://" + address + path //http://localhost:8080/heath

		resp, err := r.ClientHealthCheck.Get(reqURL)
		if err != nil {
			r.Remove(idx)
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			r.Remove(idx)
		}
	}
	return nil
}

// func (r *Registry) Register(instance *Instance)

// func (r *Registry) Deregister(id string)

// func (r *Registry) Heartbeat(id string)

// func (r *Registry) Get(id string) (*Instance, bool)

// func (r *Registry) List() []*Instance
