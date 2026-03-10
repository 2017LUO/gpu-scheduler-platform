package cluster

type MIGDevice struct {
	ID            string
	ParentGPUUUID string
	NodeName      string
	Profile       string
	GI            int
	CI            int
	MemoryMiB     int64
	Healthy       bool
	Allocated     bool
	Reserved      bool
}
