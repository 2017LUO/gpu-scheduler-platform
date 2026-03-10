package router

import "net/http"

func ApplyMiddlewares(h http.Handler) http.Handler {
	if h == nil {
		return http.NotFoundHandler()
	}
	return h
}
