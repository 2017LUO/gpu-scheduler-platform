package agent

import (
	"context"
	"fmt"
)

func (a *App) run(ctx context.Context) error {
	if a == nil {
		return fmt.Errorf("agent app is nil")
	}
	if a.Service == nil {
		return fmt.Errorf("agent service is nil")
	}

	runCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	errCh := make(chan error, 2)

	go func() {
		errCh <- a.Service.Run(runCtx)
	}()

	go func() {
		errCh <- bootstrapRunHTTPServer(a, runCtx)
	}()

	var firstErr error
	for i := 0; i < 2; i++ {
		select {
		case <-ctx.Done():
			return nil

		case err := <-errCh:
			// 有任意子任务退出后，先取消其余任务
			cancel()

			// nil 表示某个后台任务正常退出；继续等另一个任务退出
			if err == nil {
				continue
			}

			// 记录第一个非 nil 错误，但继续收另一个 goroutine 的退出
			if firstErr == nil {
				firstErr = err
			}
		}
	}

	return firstErr
}
