package algorithm

import (
	"context"
	"sort"

	model "gpu-scheduler-platform/internal/repo/models"
	schedframework "gpu-scheduler-platform/internal/scheduler/framework"
)

func SelectNode(
	ctx context.Context,
	deps Dependencies,
	cs *schedframework.CycleState,
	job *model.GPUJob,
	nodes []*model.Node,
) (*model.Node, []string, map[string]int64, map[string][]string, *schedframework.Status) {
	if job == nil {
		return nil, nil, nil, nil, schedframework.NewStatus(schedframework.CodeError, "job is nil")
	}
	if len(nodes) == 0 {
		return nil, nil, nil, nil, schedframework.NewStatus(schedframework.CodeUnschedulable, "no schedulable nodes")
	}

	inventory := schedframework.NodeGPUInventory{}
	if deps.LoadNodeInventory != nil {
		var err error
		inventory, err = deps.LoadNodeInventory(ctx, nodes)
		if err != nil {
			return nil, nil, nil, nil, schedframework.AsError(err)
		}
	}
	cs.Write(schedframework.StateKeyNodeGPUInventory, inventory)

	candidates := make([]NodeScore, 0, len(nodes))
	filterReasons := make(map[string][]string, len(nodes))
	scoreMap := make(map[string]int64, len(nodes))

	for _, node := range nodes {
		if node == nil {
			continue
		}

		if st := deps.Framework.RunFilter(cs, job, node); !st.IsSuccess() {
			filterReasons[node.NodeName] = append([]string(nil), st.Reasons()...)
			continue
		}

		score, st := deps.Framework.RunScore(cs, job, node)
		if !st.IsSuccess() {
			filterReasons[node.NodeName] = append([]string(nil), st.Reasons()...)
			continue
		}

		candidates = append(candidates, NodeScore{
			Node:  node,
			Score: score,
		})
		scoreMap[node.NodeName] = score
	}

	if len(candidates) == 0 {
		return nil, nil, scoreMap, filterReasons, schedframework.NewStatus(schedframework.CodeUnschedulable, "all nodes filtered out")
	}

	sort.SliceStable(candidates, func(i, j int) bool {
		if candidates[i].Score != candidates[j].Score {
			return candidates[i].Score > candidates[j].Score
		}
		if candidates[i].Node.FreeMemoryMiB != candidates[j].Node.FreeMemoryMiB {
			return candidates[i].Node.FreeMemoryMiB > candidates[j].Node.FreeMemoryMiB
		}
		if candidates[i].Node.HealthyGPUCount != candidates[j].Node.HealthyGPUCount {
			return candidates[i].Node.HealthyGPUCount > candidates[j].Node.HealthyGPUCount
		}
		return candidates[i].Node.NodeName < candidates[j].Node.NodeName
	})

	topK := deps.TopK
	if topK <= 0 || topK > len(candidates) {
		topK = len(candidates)
	}

	candidateNames := make([]string, 0, topK)
	for i := 0; i < topK; i++ {
		candidateNames = append(candidateNames, candidates[i].Node.NodeName)
	}

	selected := candidates[0].Node
	cs.Write(schedframework.StateKeySelectedNodeName, selected.NodeName)

	return selected, candidateNames, scoreMap, filterReasons, nil
}
