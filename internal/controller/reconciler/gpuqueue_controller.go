package reconciler

import "context"

type GPUQueueController struct{}

func NewGPUQueueController() *GPUQueueController {
	return &GPUQueueController{}
}

func (c *GPUQueueController) Reconcile(ctx context.Context, key string) error {
	_ = ctx
	_ = key
	return nil
}
