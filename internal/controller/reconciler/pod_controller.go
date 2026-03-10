package reconciler

import "context"

type PodController struct{}

func NewPodController() *PodController {
	return &PodController{}
}

func (c *PodController) Reconcile(ctx context.Context, key string) error {
	_ = ctx
	_ = key
	return nil
}
