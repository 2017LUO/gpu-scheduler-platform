package collector

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type DCGMCollector struct {
	fallback *NvidiaSMICollector
}

func NewDCGMCollector() *DCGMCollector {
	return &DCGMCollector{
		fallback: NewNvidiaSMICollector(),
	}
}

// Collect：优先使用 dcgmi 采集；若 dcgmi 不可用或解析失败，则回退到 nvidia-smi。
// 当前实现重点采集“基础 GPU 资源视图”，便于和现有 GPUInfo 对齐。
func (c *DCGMCollector) Collect(ctx context.Context) ([]GPUInfo, error) {
	gpus, err := c.collectWithDCGMI(ctx)
	if err == nil && len(gpus) > 0 {
		return gpus, nil
	}

	// 回退到 nvidia-smi，保证 agent 不中断
	fallbackGPUs, fallbackErr := c.fallback.Collect(ctx)
	if fallbackErr != nil {
		if err != nil {
			return nil, fmt.Errorf("dcgmi collect failed: %v; fallback nvidia-smi failed: %w", err, fallbackErr)
		}
		return nil, fallbackErr
	}
	return fallbackGPUs, nil
}

func (c *DCGMCollector) collectWithDCGMI(ctx context.Context) ([]GPUInfo, error) {
	// 说明：
	// dcgmi discovery -l 输出通常类似：
	// GPU 0: GPU-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	// GPU 1: GPU-yyyyyyyy-yyyy-yyyy-yyyy-yyyyyyyyyyyy
	//
	// 为了保持和现有 GPUInfo 对齐，这里只先拿 index + uuid，
	// 然后再结合 nvidia-smi 的静态信息补 model/memory。
	cmd := exec.CommandContext(ctx, "dcgmi", "discovery", "-l")
	out, err := cmd.Output()
	if err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return nil, fmt.Errorf("dcgmi not found in PATH")
		}
		return nil, fmt.Errorf("run dcgmi discovery -l: %w", err)
	}

	lines := splitNonEmptyLines(string(out))
	if len(lines) == 0 {
		return nil, fmt.Errorf("empty dcgmi discovery output")
	}

	// 先从 dcgmi 拿 index/uuid
	dcgmBasic := make(map[int]string)
	for _, line := range lines {
		index, uuid, ok := parseDCGMDiscoveryLine(line)
		if !ok {
			continue
		}
		dcgmBasic[index] = uuid
	}
	if len(dcgmBasic) == 0 {
		return nil, fmt.Errorf("no valid gpu entries parsed from dcgmi output")
	}

	// 再通过 nvidia-smi 补 model / memory.total / memory.free
	baseInfos, err := c.fallback.Collect(ctx)
	if err != nil {
		return nil, fmt.Errorf("dcgmi parsed basic info, but enrich with nvidia-smi failed: %w", err)
	}

	outGPUs := make([]GPUInfo, 0, len(baseInfos))
	for _, gpu := range baseInfos {
		uuid, ok := dcgmBasic[gpu.Index]
		if !ok {
			// dcgmi 没列出这个 index，就跳过
			continue
		}
		if uuid != "" {
			gpu.UUID = uuid
		}
		outGPUs = append(outGPUs, gpu)
	}

	if len(outGPUs) == 0 {
		return nil, fmt.Errorf("dcgmi discovery and nvidia-smi enrichment produced empty gpu list")
	}

	return outGPUs, nil
}

func parseDCGMDiscoveryLine(line string) (index int, uuid string, ok bool) {
	// 常见格式：
	// GPU 0: GPU-3d1f6d1d-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	line = strings.TrimSpace(line)
	if line == "" {
		return 0, "", false
	}
	if !strings.HasPrefix(line, "GPU ") {
		return 0, "", false
	}

	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return 0, "", false
	}

	left := strings.TrimSpace(parts[0])  // GPU 0
	right := strings.TrimSpace(parts[1]) // GPU-uuid

	leftFields := strings.Fields(left)
	if len(leftFields) != 2 {
		return 0, "", false
	}

	idx, err := strconv.Atoi(leftFields[1])
	if err != nil {
		return 0, "", false
	}

	return idx, right, true
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
