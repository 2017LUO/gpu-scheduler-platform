package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	appcfg "gpu-scheduler-platform/internal/config"

	"go.uber.org/zap"
)

func NewHTTPServer(cfg appcfg.HTTPServerConfig, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:              cfg.Addr,
		Handler:           handler,
		ReadTimeout:       cfg.ReadTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
		ReadHeaderTimeout: 5 * time.Second,
	}
}

func NewHTTPSServer(cfg appcfg.HTTPSServerConfig, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:              cfg.Addr,
		Handler:           handler,
		ReadTimeout:       cfg.ReadTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		ReadHeaderTimeout: 5 * time.Second,
	}
}

func RunHTTPServer(ctx context.Context, lg *zap.Logger, srv *http.Server) error {
	errCh := make(chan error, 1)

	go func() {
		if lg != nil {
			lg.Info("http server starting", zap.String("addr", srv.Addr))
		}
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- fmt.Errorf("listen and serve: %w", err)
			return
		}
		errCh <- nil
	}()

	select {
	case <-ctx.Done():
		return nil
	case err := <-errCh:
		return err
	}
}

func RunHTTPSServer(ctx context.Context, lg *zap.Logger, srv *http.Server, certFile, keyFile string) error {
	errCh := make(chan error, 1)

	go func() {
		if lg != nil {
			lg.Info("https server starting", zap.String("addr", srv.Addr))
		}
		if err := srv.ListenAndServeTLS(certFile, keyFile); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- fmt.Errorf("listen and serve tls: %w", err)
			return
		}
		errCh <- nil
	}()

	select {
	case <-ctx.Done():
		return nil
	case err := <-errCh:
		return err
	}
}

func ShutdownHTTPServer(lg *zap.Logger, srv *http.Server, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if lg != nil {
		lg.Info("http server shutting down", zap.String("addr", srv.Addr), zap.Duration("timeout", timeout))
	}
	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutdown http server: %w", err)
	}
	return nil
}
