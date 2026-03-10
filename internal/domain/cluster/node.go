package cluster

type NodeState string

const (
	NodeStateUnknown  NodeState = "unknown"
	NodeStateReady    NodeState = "ready"
	NodeStateNotReady NodeState = "not_ready"
)

type Node struct {
	Name        string
	State       NodeState
	Schedulable bool

	Labels      map[string]string
	Annotations map[string]string

	GPUs     []GPU
	MIGs     []MIGDevice
	Topology Topology
}

func (n Node) HealthyGPUs() []GPU {
	out := make([]GPU, 0, len(n.GPUs))
	for _, g := range n.GPUs {
		if g.Healthy && !g.Allocated && !g.Reserved {
			out = append(out, g)
		}
	}
	return out
}
