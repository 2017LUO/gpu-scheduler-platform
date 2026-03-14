package collector

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"sort"
	"strings"
)

type GPULink struct {
	From string `json:"from"`
	To   string `json:"to"`
	Type string `json:"type"`
}

type TopologyCollector struct{}

func NewTopologyCollector() *TopologyCollector {
	return &TopologyCollector{}
}

func (c *TopologyCollector) Collect(ctx context.Context) ([]GPULink, error) {
	cmd := exec.CommandContext(ctx, "nvidia-smi", "topo", "-m")
	out, err := cmd.Output()
	if err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return nil, fmt.Errorf("nvidia-smi not found in PATH")
		}
		return nil, fmt.Errorf("run nvidia-smi topo -m: %w", err)
	}

	lines := splitNonEmptyLines(string(out))
	if len(lines) == 0 {
		return nil, fmt.Errorf("nvidia-smi topo -m returned empty output")
	}

	links, err := parseTopologyMatrix(lines)
	if err != nil {
		return nil, err
	}

	sort.Slice(links, func(i, j int) bool {
		if links[i].From == links[j].From {
			if links[i].To == links[j].To {
				return links[i].Type < links[j].Type
			}
			return links[i].To < links[j].To
		}
		return links[i].From < links[j].From
	})

	return links, nil
}

func parseTopologyMatrix(lines []string) ([]GPULink, error) {
	// 典型输出：
	//
	//         GPU0    GPU1    GPU2    GPU3    CPU Affinity    NUMA Affinity
	// GPU0     X      NV4     SYS     SYS     0-31            0
	// GPU1    NV4      X      SYS     SYS     0-31            0
	// GPU2    SYS     SYS      X      NV4     32-63           1
	// GPU3    SYS     SYS     NV4      X      32-63           1
	//
	// 我们只提取 GPU-GPU 关系，忽略 CPU/NUMA 等额外列。

	if len(lines) < 2 {
		return nil, fmt.Errorf("invalid topology matrix: not enough lines")
	}

	headerFields := strings.Fields(lines[0])
	if len(headerFields) == 0 {
		return nil, fmt.Errorf("invalid topology matrix header")
	}

	// 只保留表头中的 GPU 列
	gpuHeaders := make([]string, 0, len(headerFields))
	for _, f := range headerFields {
		if isGPUHeader(f) {
			gpuHeaders = append(gpuHeaders, f)
		} else {
			break
		}
	}
	if len(gpuHeaders) == 0 {
		return nil, fmt.Errorf("no gpu headers found in topology matrix")
	}

	seen := make(map[string]struct{}, len(gpuHeaders)*len(gpuHeaders))
	links := make([]GPULink, 0, len(gpuHeaders)*2)

	for _, line := range lines[1:] {
		fields := strings.Fields(line)
		if len(fields) < len(gpuHeaders)+1 {
			continue
		}

		rowGPU := fields[0]
		if !isGPUHeader(rowGPU) {
			continue
		}

		for colIdx, colGPU := range gpuHeaders {
			linkType := strings.TrimSpace(fields[colIdx+1])

			// 跳过自身和空值
			if rowGPU == colGPU || linkType == "" || strings.EqualFold(linkType, "X") {
				continue
			}

			fromID := normalizeGPUHeader(rowGPU)
			toID := normalizeGPUHeader(colGPU)

			// 无向图去重：gpu-0|gpu-1 这样的 key
			k1, k2 := fromID, toID
			if k1 > k2 {
				k1, k2 = k2, k1
			}
			key := k1 + "|" + k2

			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}

			links = append(links, GPULink{
				From: fromID,
				To:   toID,
				Type: normalizeLinkType(linkType),
			})
		}
	}

	// 没有拓扑关系时返回空列表是合理的，例如单卡节点
	return links, nil
}

func isGPUHeader(s string) bool {
	s = strings.TrimSpace(strings.ToUpper(s))
	return strings.HasPrefix(s, "GPU")
}

func normalizeGPUHeader(s string) string {
	// GPU0 -> gpu-0
	s = strings.TrimSpace(strings.ToUpper(s))
	s = strings.TrimPrefix(s, "GPU")
	return "gpu-" + s
}

func normalizeLinkType(s string) string {
	s = strings.TrimSpace(strings.ToUpper(s))
	switch {
	case strings.HasPrefix(s, "NV"):
		return s // NV1/NV2/NV4...
	case s == "PIX":
		return "PIX"
	case s == "PXB":
		return "PXB"
	case s == "PHB":
		return "PHB"
	case s == "NODE":
		return "NODE"
	case s == "SYS":
		return "SYS"
	default:
		return s
	}
}
