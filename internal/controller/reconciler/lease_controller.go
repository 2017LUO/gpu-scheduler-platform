package reconciler

import "context"

type LeaseController struct{}

func NewLeaseController() *LeaseController {
	return &LeaseController{}
}

func (c *LeaseController) Reconcile(ctx context.Context, key string) error {
	_ = ctx
	_ = key
	return nil
}
