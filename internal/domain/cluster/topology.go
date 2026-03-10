package cluster

type LinkType string

const (
	LinkUnknown LinkType = "unknown"
	LinkPCIe    LinkType = "pcie"
	LinkNVLink  LinkType = "nvlink"
)

type TopologyLink struct {
	FromGPU string
	ToGPU   string
	Type    LinkType
	Weight  int
}

type Topology struct {
	NodeName string
	Links    []TopologyLink
}

func (t Topology) LinkBetween(a, b string) (TopologyLink, bool) {
	for _, l := range t.Links {
		if (l.FromGPU == a && l.ToGPU == b) || (l.FromGPU == b && l.ToGPU == a) {
			return l, true
		}
	}
	return TopologyLink{}, false
}
