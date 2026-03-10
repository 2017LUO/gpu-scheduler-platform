package algorithm

import (
	"sort"

	"gpu-scheduler-platform/internal/domain/cluster"
	"gpu-scheduler-platform/internal/domain/job"
)

func SelectNode(snapshot *cluster.Snapshot, j job.Job) (*PlacementDecision, bool) {
	if snapshot == nil {
		return nil, false
	}

	nodes := snapshot.ReadyNodes()
	scores := make([]NodeScore, 0, len(nodes))

	for _, n := range nodes {
		gpus, ok := SelectGPUs(n, j.Requirement)
		if !ok {
			continue
		}
		score := scoreNode(n, gpus, j)
		scores = append(scores, NodeScore{
			Node:  n,
			Score: score,
		})
	}

	if len(scores) == 0 {
		return nil, false
	}

	sort.Slice(scores, func(i, j int) bool {
		if scores[i].Score == scores[j].Score {
			return scores[i].Node.Name < scores[j].Node.Name
		}
		return scores[i].Score > scores[j].Score
	})

	best := scores[0]
	gpus, ok := SelectGPUs(best.Node, j.Requirement)
	if !ok {
		return nil, false
	}

	return &PlacementDecision{
		Job:    j,
		Node:   best.Node,
		GPUs:   gpus,
		Reason: "best-fit",
	}, true
}

func scoreNode(n cluster.Node, gpus []cluster.GPU, j job.Job) int {
	score := 0

	if j.Requirement.RequireSameNode {
		score += 100
	}

	// 倾向刚好能放下的小节点，保留大块资源
	freeCount := 0
	for _, g := range n.GPUs {
		if !g.Allocated && !g.Reserved && (!j.Requirement.RequireHealthy || g.Healthy) {
			freeCount++
		}
	}
	score += 100 - freeCount

	// 倾向更小显存满足
	var totalFreeMem int64
	for _, g := range gpus {
		totalFreeMem += g.FreeMemoryMiB
	}
	score += int(1_000_000 / max64(totalFreeMem, 1))

	if j.Requirement.RequireNVLink && hasNVLinkBetweenAll(n, gpus) {
		score += 200
	}

	return score
}

func hasNVLinkBetweenAll(n cluster.Node, gpus []cluster.GPU) bool {
	if len(gpus) <= 1 {
		return true
	}
	for i := 0; i < len(gpus); i++ {
		for j := i + 1; j < len(gpus); j++ {
			link, ok := n.Topology.LinkBetween(gpus[i].ID, gpus[j].ID)
			if !ok || link.Type != cluster.LinkNVLink {
				return false
			}
		}
	}
	return true
}

func max64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
