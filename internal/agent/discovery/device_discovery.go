package discovery

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type DiscoveredGPU struct {
	ID         string `json:"id"`
	UUID       string `json:"uuid"`
	Index      int    `json:"index"`
	Model      string `json:"model"`
	MemoryMiB  int64  `json:"memory_mib"`
	BusID      string `json:"bus_id,omitempty"`
	DevicePath string `json:"device_path,omitempty"`
	MIGCapable bool   `json:"mig_capable"`
}

type DeviceFileInfo struct {
	Path string `json:"path"`
	Kind string `json:"kind"`
}

type DeviceInventory struct {
	NodeName       string           `json:"node_name"`
	GPUCount       int              `json:"gpu_count"`
	TotalMemoryMiB int64            `json:"total_memory_mib"`
	GPUs           []DiscoveredGPU  `json:"gpus"`
	DeviceFiles    []DeviceFileInfo `json:"device_files,omitempty"`
}

type DeviceDiscovery struct{}

func NewDeviceDiscovery() *DeviceDiscovery {
	return &DeviceDiscovery{}
}

// Discover：发现节点的静态 GPU 资产信息。
// 适合在 agent 启动时调用一次，或低频刷新调用。
func (d *DeviceDiscovery) Discover(ctx context.Context, nodeName string) (*DeviceInventory, error) {
	gpus, err := d.discoverGPUs(ctx)
	if err != nil {
		return nil, err
	}

	deviceFiles := d.discoverDeviceFiles()

	inv := &DeviceInventory{
		NodeName:    nodeName,
		GPUs:        gpus,
		DeviceFiles: deviceFiles,
		GPUCount:    len(gpus),
	}

	var total int64
	for _, g := range gpus {
		total += g.MemoryMiB
	}
	inv.TotalMemoryMiB = total

	return inv, nil
}

func (d *DeviceDiscovery) discoverGPUs(ctx context.Context) ([]DiscoveredGPU, error) {
	// 先查基础静态信息
	cmd := exec.CommandContext(
		ctx,
		"nvidia-smi",
		"--query-gpu=index,uuid,name,memory.total,pci.bus_id",
		"--format=csv,noheader,nounits",
	)

	out, err := cmd.Output()
	if err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return nil, fmt.Errorf("nvidia-smi not found in PATH")
		}
		return nil, fmt.Errorf("run nvidia-smi discover: %w", err)
	}

	lines := splitNonEmptyLines(string(out))
	if len(lines) == 0 {
		return nil, fmt.Errorf("no gpu discovered from nvidia-smi")
	}

	// 再额外查 MIG mode；查不到不影响主流程
	migMap := d.discoverMIGCapability(ctx)

	gpus := make([]DiscoveredGPU, 0, len(lines))
	for _, line := range lines {
		parts := splitCSVLine(line)
		if len(parts) < 5 {
			continue
		}

		index, err := strconv.Atoi(parts[0])
		if err != nil {
			continue
		}

		memMiB, err := strconv.ParseInt(parts[3], 10, 64)
		if err != nil {
			memMiB = 0
		}

		gpu := DiscoveredGPU{
			ID:         fmt.Sprintf("gpu-%d", index),
			UUID:       parts[1],
			Index:      index,
			Model:      parts[2],
			MemoryMiB:  memMiB,
			BusID:      parts[4],
			DevicePath: fmt.Sprintf("/dev/nvidia%d", index),
			MIGCapable: migMap[index],
		}

		// 设备文件不存在时，不强行填
		if _, err := os.Stat(gpu.DevicePath); err != nil {
			gpu.DevicePath = ""
		}

		gpus = append(gpus, gpu)
	}

	sort.Slice(gpus, func(i, j int) bool {
		return gpus[i].Index < gpus[j].Index
	})

	if len(gpus) == 0 {
		return nil, fmt.Errorf("gpu discovery returned empty result")
	}

	return gpus, nil
}

// discoverMIGCapability：根据 mig.mode.current 判断 GPU 是否具备 MIG 能力。
// 规则：
// - 返回值为 "N/A"：通常表示不支持 MIG
// - 返回值为 "Enabled"/"Disabled"：通常表示支持 MIG
func (d *DeviceDiscovery) discoverMIGCapability(ctx context.Context) map[int]bool {
	cmd := exec.CommandContext(
		ctx,
		"nvidia-smi",
		"--query-gpu=index,mig.mode.current",
		"--format=csv,noheader,nounits",
	)

	out, err := cmd.Output()
	if err != nil {
		return map[int]bool{}
	}

	lines := splitNonEmptyLines(string(out))
	m := make(map[int]bool, len(lines))

	for _, line := range lines {
		parts := splitCSVLine(line)
		if len(parts) < 2 {
			continue
		}

		index, err := strconv.Atoi(parts[0])
		if err != nil {
			continue
		}

		mode := strings.TrimSpace(parts[1])
		if mode == "" || strings.EqualFold(mode, "N/A") {
			m[index] = false
			continue
		}

		m[index] = true
	}

	return m
}

func (d *DeviceDiscovery) discoverDeviceFiles() []DeviceFileInfo {
	paths := make([]string, 0, 16)

	// 常见 NVIDIA 设备节点
	patterns := []string{
		"/dev/nvidia*",
		"/dev/nvidia-caps/*",
		"/dev/nvidia-uvm*",
	}

	for _, p := range patterns {
		matches, err := filepath.Glob(p)
		if err != nil {
			continue
		}
		paths = append(paths, matches...)
	}

	seen := make(map[string]struct{}, len(paths))
	out := make([]DeviceFileInfo, 0, len(paths))

	for _, p := range paths {
		if _, ok := seen[p]; ok {
			continue
		}
		seen[p] = struct{}{}

		out = append(out, DeviceFileInfo{
			Path: p,
			Kind: classifyDeviceFile(p),
		})
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].Path < out[j].Path
	})

	return out
}

func classifyDeviceFile(path string) string {
	base := filepath.Base(path)

	switch {
	case strings.HasPrefix(base, "nvidia-caps"):
		return "mig-cap"
	case base == "nvidiactl":
		return "control"
	case strings.HasPrefix(base, "nvidia-uvm"):
		return "uvm"
	case strings.HasPrefix(base, "nvidia") && isAllDigits(strings.TrimPrefix(base, "nvidia")):
		return "gpu"
	default:
		return "other"
	}
}

func splitNonEmptyLines(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}

	raw := strings.Split(s, "\n")
	out := make([]string, 0, len(raw))
	for _, line := range raw {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		out = append(out, line)
	}
	return out
}

func splitCSVLine(line string) []string {
	parts := strings.Split(line, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		out = append(out, strings.TrimSpace(p))
	}
	return out
}

func isAllDigits(s string) bool {
	if s == "" {
		return false
	}
	for _, ch := range s {
		if ch < '0' || ch > '9' {
			return false
		}
	}
	return true
}
