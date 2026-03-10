package reconciler

import "context"

type NodeController struct{}

func NewNodeController() *NodeController {
	return &NodeController{}
}

func (c *NodeController) Reconcile(ctx context.Context, key string) error {
	_ = ctx
	_ = key
	return nil
}
