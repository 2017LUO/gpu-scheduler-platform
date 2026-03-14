package apiserver

import "time"

func (a *App) shutdownHTTP(timeout time.Duration) error {
	if timeout <= 0 {
		timeout = 20 * time.Second
	}
	return shutdownHTTPServer(a, timeout)
}
