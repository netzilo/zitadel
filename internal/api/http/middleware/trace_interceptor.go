package middleware

import (
	"net/http"

	http_utils "github.com/caos/zitadel/internal/api/http"
	"github.com/caos/zitadel/internal/tracing"
)

func DefaultTraceHandler(handler http.Handler) http.Handler {
	return tracing.TraceHandler(handler, http_utils.Probes...)
}

func TraceHandler(ignoredMethods ...string) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return tracing.TraceHandler(handler, ignoredMethods...)
	}
}
