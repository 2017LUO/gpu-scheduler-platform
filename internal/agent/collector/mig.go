package collector

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type MIGInfo struct {
	ID         string `json:"id"`
	ParentUUID string `json:"parent_uuid"`
	Profile    string `json:"profile"`
	MemoryMiB  int64  `json:"memory_mib"`
}

type MIGCollector struct{}

func NewMIGCollector() *MIGCollector {
	return &MIGCollector{}
}

func (c *MIGCollector) Collect(ctx context.Context) ([]MIGInfo, error) {
	cmd := exec.CommandContext(ctx, "nvidia-smi", "-L")
	out, err := cmd.Output()
	if err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return nil, fmt.Errorf("nvidia-smi not found in PATH")
		}
		return nil, fmt.Errorf("run nvidia-smi -L: %w", err)
	}

	lines := splitNonEmptyLines(string(out))
	if len(lines) == 0 {
		return nil, fmt.Errorf("nvidia-smi -L returned empty output")
	}

	migs, err := parseMIGLines(lines)
	if err != nil {
		return nil, err
	}

	sort.Slice(migs, func(i, j int) bool {
		if migs[i].ParentUUID == migs[j].ParentUUID {
			return migs[i].ID < migs[j].ID
		}
		return migs[i].ParentUUID < migs[j].ParentUUID
	})

	return migs, nil
}

func parseMIGLines(lines []string) ([]MIGInfo, error) {
	// 例子：
	// GPU 0: NVIDIA A100-SXM4-40GB (UUID: GPU-aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee)
	//   MIG 1g.5gb Device 0: (UUID: MIG-11111111-2222-3333-4444-555555555555)

	gpuUUIDRe := regexp.MustCompile(`UUID:\s*(GPU-[^)]+)\)`)
	migRe := regexp.MustCompile(`MIG\s+([0-9]+g\.[0-9]+gb)\s+Device\s+[0-9]+:\s+\(UUID:\s*(MIG-[^)]+)\)`)

	var currentParentUUID string
	out := make([]MIGInfo, 0)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// 先看是不是 GPU 行
		if strings.HasPrefix(line, "GPU ") {
			m := gpuUUIDRe.FindStringSubmatch(line)
			if len(m) >= 2 {
				currentParentUUID = strings.TrimSpace(m[1])
			} else {
				currentParentUUID = ""
			}
			continue
		}

		// 再看是不是 MIG 行
		if strings.HasPrefix(line, "MIG ") {
			m := migRe.FindStringSubmatch(line)
			if len(m) < 3 {
				continue
			}

			profile := strings.TrimSpace(m[1])
			migUUID := strings.TrimSpace(m[2])

			memMiB := parseMIGProfileMemoryMiB(profile)

			out = append(out, MIGInfo{
				ID:         migUUID,
				ParentUUID: currentParentUUID,
				Profile:    profile,
				MemoryMiB:  memMiB,
			})
		}
	}

	// MIG 没开时，返回空列表是合理的，不视为错误
	return out, nil
}

// parseMIGProfileMemoryMiB:
// 从 profile 如 "1g.5gb"、"2g.10gb"、"3g.20gb" 里提取显存容量并转成 MiB。
func parseMIGProfileMemoryMiB(profile string) int64 {
	// 匹配 ".5gb" / ".10gb" / ".20gb"
	re := regexp.MustCompile(`\.([0-9]+)gb$`)
	m := re.FindStringSubmatch(strings.ToLower(strings.TrimSpace(profile)))
	if len(m) < 2 {
		return 0
	}

	gb, err := strconv.ParseInt(m[1], 10, 64)
	if err != nil || gb < 0 {
		return 0
	}

	return gb * 1024
}
