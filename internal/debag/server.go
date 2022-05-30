package debag

import (
	"context"
	"net/http"

	"github.com/Dsmit05/metida/internal/debag/handlers"
	"github.com/Dsmit05/metida/internal/logger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const maxHeaderBytes = 1 << 20 // controls the maximum number of bytes in header

type DebagServer struct {
	s *http.Server
}

func NewDebagServer(cfg configDebagI) *DebagServer {
	r := http.NewServeMux()
	profiling := handlers.NewProfilingHandler()
	swagg := handlers.NewSwaggerHandler()
	r.Handle("/pprof/", profiling)
	r.Handle("/swagger/", swagg)
	r.Handle("/metrics/", promhttp.Handler())

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		cfg.GetConfigInfo(w, r)
	})

	s := &http.Server{
		Addr:           cfg.GetDebagAddr(),
		Handler:        r,
		ReadTimeout:    cfg.GetDebagReadTimeout(),
		WriteTimeout:   cfg.GetDebagWriteTimeout(),
		MaxHeaderBytes: maxHeaderBytes,
	}

	return &DebagServer{s}
}

func (o *DebagServer) Start() {
	if err := o.s.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			logger.Info("DebagServer", "Server closed under request")
		} else {
			logger.Error("DebagServer closed unexpect", err)
		}
	}
}

func (o *DebagServer) Stop(ctx context.Context) {
	stop := make(chan bool)

	go func() {
		o.s.SetKeepAlivesEnabled(false)
		if err := o.s.Shutdown(ctx); err != nil {
			logger.Error("DebagServer Shutdown failed", err)
		}
		stop <- true
	}()

	select {
	case <-ctx.Done():
		logger.Error("DebagServer context timeout", ctx.Err())
	case <-stop:
		logger.Info("DebagServer", "Stop")
	}

}
