package metrics

import (
	"bufio"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"net"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
)

// MetricsMiddleware middleware to collect metrics from http requests
func (o *ServiceMetrics) MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wi := &responseWriterInterceptor{
			code:           http.StatusOK,
			ResponseWriter: w,
		}

		startedAt := time.Now()
		next.ServeHTTP(wi, r)
		elapsed := time.Since(startedAt)

		o.httpRequestDurations.With(prometheus.Labels{
			"method": r.Method, "path": filepath.Dir(r.RequestURI), "code": strconv.Itoa(wi.code)}).
			Observe(elapsed.Seconds())

		o.httpRequestCounters.With(prometheus.Labels{
			"method": r.Method, "path": filepath.Dir(r.RequestURI), "code": strconv.Itoa(wi.code)}).
			Add(1)
	})
}

// responseWriterInterceptor is a simple wrapper to intercept set data
type responseWriterInterceptor struct {
	http.ResponseWriter
	code int
}

func (w *responseWriterInterceptor) WriteHeader(statusCode int) {
	w.code = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *responseWriterInterceptor) Write(p []byte) (int, error) {
	return w.ResponseWriter.Write(p)
}

func (w *responseWriterInterceptor) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := w.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("type assertion failed http.ResponseWriter not a http.Hijacker")
	}
	return h.Hijack()
}

func (w *responseWriterInterceptor) Flush() {
	f, ok := w.ResponseWriter.(http.Flusher)
	if !ok {
		return
	}

	f.Flush()
}
