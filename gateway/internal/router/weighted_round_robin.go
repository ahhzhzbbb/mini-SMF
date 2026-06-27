package router

import (
	"fmt"
	"mini-SMF/gateway/internal/registry"
	"sync"
)

type WeightedRoundRobin struct {
	mu sync.Mutex
}

func NewWeightRoundRobin() *WeightedRoundRobin {
	return &WeightedRoundRobin{}
}

func (w *WeightedRoundRobin) Next(reg *registry.Registry) (*registry.Instance, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	winnerIndex := 0
	for i, instance := range reg.Instances {
		instance.CurrentWeight += instance.Weight
		fmt.Printf("%d => ", instance.CurrentWeight)
		if instance.CurrentWeight > reg.Instances[winnerIndex].CurrentWeight {
			winnerIndex = i
		}
	}
	fmt.Print("\n")
	reg.Instances[winnerIndex].CurrentWeight -= reg.SumOfWeight
	return reg.Instances[winnerIndex], nil
}
