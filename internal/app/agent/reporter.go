package agent

import (
	"context"
	"fmt"
)

func (a *App) run(ctx context.Context) error {
	if a.Service == nil {
		return fmt.Errorf("agent service is nil")
	}

	errCh := make(chan error, 2)

	go func() {
		errCh <- a.Service.Run(ctx)
	}()

	go func() {
		errCh <- bootstrapRunHTTPServer(a, ctx)
	}()

	select {
	case <-ctx.Done():
		return nil
	case err := <-errCh:
		return err
	}
}
