package reconciler

import "context"

type GPUQuotaController struct{}

func NewGPUQuotaController() *GPUQuotaController {
	return &GPUQuotaController{}
}

func (c *GPUQuotaController) Reconcile(ctx context.Context, key string) error {
	_ = ctx
	_ = key
	return nil
}
