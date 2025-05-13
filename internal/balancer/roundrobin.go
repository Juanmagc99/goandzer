package balancer

import "sync"

type RoundRobin struct {
	targets []string
	mu      sync.Mutex
	idx     int
}

func NewRoundRobin(targets []string) *RoundRobin {
	return &RoundRobin{targets: targets}
}

func (r *RoundRobin) Next() string {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(r.targets) == 0 {
		return ""
	}

	target := r.targets[r.idx%len(r.targets)]
	r.idx++

	return target
}

func (r *RoundRobin) SetTargets(ts []string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.targets = ts
}
