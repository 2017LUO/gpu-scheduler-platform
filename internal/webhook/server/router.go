package server

import (
	"net/http"

	"gpu-scheduler-platform/internal/webhook/mutating"
	"gpu-scheduler-platform/internal/webhook/validating"

	"go.uber.org/zap"
)

type MiddlewareWrapper func(http.Handler) http.Handler

type Router struct {
	logger *zap.Logger
}

func NewRouter(lg *zap.Logger) *Router {
	return &Router{logger: lg}
}

func (r *Router) Register(mux *http.ServeMux, wrap MiddlewareWrapper) {
	if mux == nil {
		return
	}

	gpuJobDefaults := http.HandlerFunc(mutating.NewGPUJobDefaultsHandler(r.logger).Handle)
	podDefaults := http.HandlerFunc(mutating.NewPodDefaultsHandler(r.logger).Handle)

	gpuJobValidate := http.HandlerFunc(validating.NewGPUJobValidateHandler(r.logger).Handle)
	gpuPolicyValidate := http.HandlerFunc(validating.NewGPUPolicyValidateHandler(r.logger).Handle)
	gpuQuotaValidate := http.HandlerFunc(validating.NewGPUQuotaValidateHandler(r.logger).Handle)
	podValidate := http.HandlerFunc(validating.NewPodValidateHandler(r.logger).Handle)

	if wrap != nil {
		gpuJobDefaults = wrap(gpuJobDefaults).ServeHTTP
		podDefaults = wrap(podDefaults).ServeHTTP
		gpuJobValidate = wrap(gpuJobValidate).ServeHTTP
		gpuPolicyValidate = wrap(gpuPolicyValidate).ServeHTTP
		gpuQuotaValidate = wrap(gpuQuotaValidate).ServeHTTP
		podValidate = wrap(podValidate).ServeHTTP
	}

	mux.Handle("/mutate/gpujobs", gpuJobDefaults)
	mux.Handle("/mutate/pods", podDefaults)

	mux.Handle("/validate/gpujobs", gpuJobValidate)
	mux.Handle("/validate/gpupolicies", gpuPolicyValidate)
	mux.Handle("/validate/gpuquotas", gpuQuotaValidate)
	mux.Handle("/validate/pods", podValidate)
}
