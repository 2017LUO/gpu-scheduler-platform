package reconciler

import "context"

type GPUJobController struct{}

func NewGPUJobController() *GPUJobController {
	return &GPUJobController{}
}

func (c *GPUJobController) Reconcile(ctx context.Context, key string) error {
	_ = ctx
	_ = key
	return nil
}
