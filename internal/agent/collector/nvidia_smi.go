package collector

import (
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type GPUInfo struct {
	ID            string `json:"id"`
	UUID          string `json:"uuid"`
	Index         int    `json:"index"`
	Model         string `json:"model"`
	MemoryMiB     int64  `json:"memory_mib"`
	FreeMemoryMiB int64  `json:"free_memory_mib"`
	Healthy       bool   `json:"healthy"`
	Health        string `json:"health"`
}

type NvidiaSMICollector struct{}

func NewNvidiaSMICollector() *NvidiaSMICollector {
	return &NvidiaSMICollector{}
}

func (c *NvidiaSMICollector) Collect(ctx context.Context) ([]GPUInfo, error) {
	cmd := exec.CommandContext(
		ctx,
		"nvidia-smi",
		"--query-gpu=index,uuid,name,memory.total,memory.free",
		"--format=csv,noheader,nounits",
	)

	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("run nvidia-smi: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	gpus := make([]GPUInfo, 0, len(lines))

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.Split(line, ",")
		if len(parts) < 5 {
			continue
		}

		index, _ := strconv.Atoi(strings.TrimSpace(parts[0]))
		uuid := strings.TrimSpace(parts[1])
		name := strings.TrimSpace(parts[2])
		totalMem, _ := strconv.ParseInt(strings.TrimSpace(parts[3]), 10, 64)
		freeMem, _ := strconv.ParseInt(strings.TrimSpace(parts[4]), 10, 64)

		gpus = append(gpus, GPUInfo{
			ID:            fmt.Sprintf("gpu-%d", index),
			UUID:          uuid,
			Index:         index,
			Model:         name,
			MemoryMiB:     totalMem,
			FreeMemoryMiB: freeMem,
			Healthy:       true,
			Health:        "healthy",
		})
	}

	return gpus, nil
}
