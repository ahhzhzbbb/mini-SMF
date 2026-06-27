package router

import "mini-SMF/gateway/internal/registry"

type LoadBalancer interface {
	Next(reg *registry.Registry) (*registry.Instance, error)
}
