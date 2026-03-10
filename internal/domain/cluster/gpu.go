package cluster

type GPUHealth string

const (
	GPUHealthUnknown   GPUHealth = "unknown"
	GPUHealthHealthy   GPUHealth = "healthy"
	GPUHealthUnhealthy GPUHealth = "unhealthy"
)

type GPUType string

const (
	GPUTypeFull GPUType = "full"
	GPUTypeMIG  GPUType = "mig"
)

type GPU struct {
	ID            string
	UUID          string
	NodeName      string
	Index         int
	Model         string
	Vendor        string
	Type          GPUType
	MemoryMiB     int64
	FreeMemoryMiB int64
	Healthy       bool
	Health        GPUHealth

	MIGEnabled bool
	MIGProfile string

	Labels      map[string]string
	Annotations map[string]string

	Allocated bool
	Reserved  bool
}
