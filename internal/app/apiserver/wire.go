package apiserver

import "time"

func (a *App) shutdownHTTP(timeoutSeconds any) error {
	timeout, ok := timeoutSeconds.(time.Duration)
	if !ok {
		timeout = 20 * time.Second
	}
	return shutdownHTTPServer(a, timeout)
}
