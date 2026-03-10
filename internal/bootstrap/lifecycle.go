package bootstrap

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"go.uber.org/zap"
)

type Hook func(context.Context) error

type Lifecycle struct {
	mu      sync.Mutex
	onStart []Hook
	onStop  []Hook
	started bool
	stopped bool
	logger  *zap.Logger
}

func NewLifecycle(lg *zap.Logger) *Lifecycle {
	return &Lifecycle{
		logger: lg,
	}
}

func (l *Lifecycle) AppendOnStart(h Hook) {
	if h == nil {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.onStart = append(l.onStart, h)
}

func (l *Lifecycle) AppendOnStop(h Hook) {
	if h == nil {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.onStop = append(l.onStop, h)
}

func (l *Lifecycle) Start(ctx context.Context) error {
	l.mu.Lock()
	if l.started {
		l.mu.Unlock()
		return nil
	}
	hooks := append([]Hook(nil), l.onStart...)
	l.started = true
	l.mu.Unlock()

	for i, h := range hooks {
		if err := h(ctx); err != nil {
			return fmt.Errorf("run onStart hook #%d: %w", i, err)
		}
	}
	return nil
}

func (l *Lifecycle) Stop(ctx context.Context) error {
	l.mu.Lock()
	if l.stopped {
		l.mu.Unlock()
		return nil
	}
	hooks := reverseHooks(append([]Hook(nil), l.onStop...))
	l.stopped = true
	l.mu.Unlock()

	var firstErr error
	for i, h := range hooks {
		if err := h(ctx); err != nil && firstErr == nil {
			firstErr = fmt.Errorf("run onStop hook #%d: %w", i, err)
		}
	}
	return firstErr
}

func WithSignalCancel(parent context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(parent)

	ch := make(chan os.Signal, 2)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		defer signal.Stop(ch)
		select {
		case <-ctx.Done():
			return
		case <-ch:
			cancel()
		}
	}()

	return ctx, cancel
}

func reverseHooks(in []Hook) []Hook {
	for i, j := 0, len(in)-1; i < j; i, j = i+1, j-1 {
		in[i], in[j] = in[j], in[i]
	}
	return in
}
