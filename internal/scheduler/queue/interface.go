package queue

import model "gpu-scheduler-platform/internal/repo/models"

type Interface interface {
	Push(*model.GPUJob) error
	Pop() (*model.GPUJob, bool)
	Len() int
	Peek() (*model.GPUJob, bool)
	Clear()
}
