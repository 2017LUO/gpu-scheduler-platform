package cluster

import "time"

type Snapshot struct {
	Version   string
	Nodes     []Node
	CreatedAt time.Time
}

func (s Snapshot) ReadyNodes() []Node {
	out := make([]Node, 0, len(s.Nodes))
	for _, n := range s.Nodes {
		if n.State == NodeStateReady && n.Schedulable {
			out = append(out, n)
		}
	}
	return out
}

func (s Snapshot) FindNode(name string) (Node, bool) {
	for _, n := range s.Nodes {
		if n.Name == name {
			return n, true
		}
	}
	return Node{}, false
}
