package framework

import model "gpu-scheduler-platform/internal/repo/models"

type NodeGPUInventory map[string][]model.GPUDevice

const (
	StateKeyNodeGPUInventory StateKey = "node_gpu_inventory"
	StateKeySelectedNodeName StateKey = "selected_node_name"
	StateKeySelectedGPUUUIDs StateKey = "selected_gpu_uuids"
	StateKeyReservationID    StateKey = "reservation_id"
)

func ReadNodeGPUInventory(cs *CycleState) (NodeGPUInventory, bool) {
	if cs == nil {
		return nil, false
	}
	v, ok := cs.Read(StateKeyNodeGPUInventory)
	if !ok {
		return nil, false
	}
	out, ok := v.(NodeGPUInventory)
	return out, ok
}

func ReadSelectedGPUUUIDs(cs *CycleState) ([]string, bool) {
	if cs == nil {
		return nil, false
	}
	v, ok := cs.Read(StateKeySelectedGPUUUIDs)
	if !ok {
		return nil, false
	}
	out, ok := v.([]string)
	return out, ok
}

func ReadReservationID(cs *CycleState) (string, bool) {
	if cs == nil {
		return "", false
	}
	v, ok := cs.Read(StateKeyReservationID)
	if !ok {
		return "", false
	}
	out, ok := v.(string)
	return out, ok
}
