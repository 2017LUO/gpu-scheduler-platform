package reconciler

import "context"

type GPUPolicyController struct{}

func NewGPUPolicyController() *GPUPolicyController {
	return &GPUPolicyController{}
}

func (c *GPUPolicyController) Reconcile(ctx context.Context, key string) error {
	_ = ctx
	_ = key
	return nil
}
