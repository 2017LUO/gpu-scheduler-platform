package collector

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
)

type PodGPUInfo struct {
	PodName   string   `json:"pod_name"`
	Namespace string   `json:"namespace"`
	GPUIDs    []string `json:"gpu_ids"`
}

type PodGPUUsageCollector struct {
	sourceFile string
}

func NewPodGPUUsageCollector() *PodGPUUsageCollector {
	return &PodGPUUsageCollector{
		sourceFile: "/var/lib/gpu-scheduler-agent/pod-gpu-usage.json",
	}
}

func (c *PodGPUUsageCollector) Collect(ctx context.Context) ([]PodGPUInfo, error) {
	_ = ctx

	data, err := os.ReadFile(c.sourceFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// 文件不存在时，认为当前没有可用映射信息
			return []PodGPUInfo{}, nil
		}
		return nil, fmt.Errorf("read pod gpu usage file %q: %w", c.sourceFile, err)
	}

	if strings.TrimSpace(string(data)) == "" {
		return []PodGPUInfo{}, nil
	}

	var items []PodGPUInfo
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, fmt.Errorf("unmarshal pod gpu usage file %q: %w", c.sourceFile, err)
	}

	out := make([]PodGPUInfo, 0, len(items))
	for _, item := range items {
		item.PodName = strings.TrimSpace(item.PodName)
		item.Namespace = strings.TrimSpace(item.Namespace)

		if item.PodName == "" || item.Namespace == "" {
			continue
		}

		item.GPUIDs = normalizeGPUIds(item.GPUIDs)
		out = append(out, item)
	}

	sort.Slice(out, func(i, j int) bool {
		if out[i].Namespace == out[j].Namespace {
			return out[i].PodName < out[j].PodName
		}
		return out[i].Namespace < out[j].Namespace
	})

	return out, nil
}

func normalizeGPUIds(ids []string) []string {
	if len(ids) == 0 {
		return []string{}
	}

	seen := make(map[string]struct{}, len(ids))
	out := make([]string, 0, len(ids))

	for _, id := range ids {
		id = strings.TrimSpace(id)
		if id == "" {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}

	sort.Strings(out)
	return out
}
