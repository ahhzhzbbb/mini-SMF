package router

import "mini-SMF/gateway/internal/registry"

type LoadBalancer interface {
	Next(instances []*registry.Instance) (*registry.Instance, error)
}
