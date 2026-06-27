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
	ID              string
	ServiceName     string
	IpAddr          string
	Port            string
	Weight          float32
	ActiveRequests  int
	TimeoutRequests int64
	LastHeartbeat   time.Time
}

func NewInstance(id, serviceName, address, port string) *Instance {
	return &Instance{
		ID:              id,
		ServiceName:     serviceName,
		IpAddr:          address,
		Port:            port,
		TimeoutRequests: 3,
	}
}

type Registry struct {
	mu                  sync.RWMutex
	Instances           []*Instance
	CycleHeathCheckTime int          //second
	ClientHealthCheck   *http.Client //client representative to send health check request:V
	ActiveInstance      map[string]bool
}

func NewRegistry() *Registry {
	transport := &http.Transport{
		Protocols: new(http.Protocols),
	}

	transport.Protocols.SetUnencryptedHTTP2(true)

	return &Registry{
		Instances:           make([]*Instance, 0),
		CycleHeathCheckTime: 1,
		ClientHealthCheck: &http.Client{
			Timeout:   100 * time.Millisecond,
			Transport: transport,
		},
		ActiveInstance: make(map[string]bool),
	}
}

func (r *Registry) Load(serviceName string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	ips, err := net.LookupIP(serviceName)
	if err != nil {
		return err
	}
	port := os.Getenv("PDU_PORT")

	for _, ip := range ips {
		if r.ActiveInstance[ip.String()] == false {
			newIntance := NewInstance(fmt.Sprintf("%s-%s:%s", serviceName, ip, port), serviceName, ip.String(), port)
			r.Instances = append(r.Instances, newIntance)
			r.ActiveInstance[ip.String()] = true
		}
	}
	return nil
}

func (r *Registry) IsEmpty() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.Instances) == 0
}

func (r *Registry) GetAllInstances() []string {
	fmt.Println(len(r.Instances))
	var res []string
	for _, i := range r.Instances {
		res = append(res, i.IpAddr)
	}
	return res
}

func (r *Registry) GetIndexOfInstance(instance *Instance) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, ins := range r.Instances {
		if ins.IpAddr == instance.IpAddr {
			return i, nil
		}
	}
	return -1, errors.New("Your instance didnt exist in registry")
}

func (r *Registry) Remove(index int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if index < 0 || index >= len(r.Instances) {
		return errors.New("out of range")
	}

	ip := r.Instances[index].IpAddr
	r.ActiveInstance[ip] = false
	r.Instances = append(r.Instances[:index], r.Instances[index+1:]...)

	return nil
}

func (r *Registry) RemoveInstance(instance *Instance) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, ins := range r.Instances {
		if ins.IpAddr == instance.IpAddr {
			ip := r.Instances[i].IpAddr
			r.ActiveInstance[ip] = false
			r.Instances = append(r.Instances[:i], r.Instances[i+1:]...)
			return nil
		}
	}
	return errors.New("Your instance didnt exist in registry")
}

func (r *Registry) HealthCheck(path string) error {
	for _, instance := range r.Instances {
		address := net.JoinHostPort(instance.IpAddr, instance.Port)
		reqURL := "http://" + address + path //http://localhost:8080/heath

		resp, err := r.ClientHealthCheck.Get(reqURL)
		if err != nil {
			r.RemoveInstance(instance)
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			r.RemoveInstance(instance)
		}
	}
	return nil
}

func (r *Registry) Register(serviceName, ip, port string) error {
	// fmt.Println("registering...")
	// name, err := net.LookupAddr(net.JoinHostPort(ip, port))
	// if err != nil {
	// 	fmt.Printf("Cant find address with %s\n", fmt.Sprintf("%s:%s", ip, port))
	// 	return err
	// }
	// fmt.Println("found!!!")
	// if len(name) != 0 {
	// 	fmt.Println("Duplicate")
	// 	return errors.New("Duplicate address, this instance exited in registry")
	// }
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.ActiveInstance[ip] {
		return errors.New("this instance exited in registry")
	}

	newInstance := NewInstance(fmt.Sprintf("%s-%s:%s", serviceName, ip, port), serviceName, ip, port)
	r.Instances = append(r.Instances, newInstance)
	r.ActiveInstance[ip] = true
	fmt.Println("register successful")
	return nil
}
