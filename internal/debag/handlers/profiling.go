package handlers

import (
	"net/http"
	"net/http/pprof"
)

// ProfilingHandler .
type ProfilingHandler struct {
	r *http.ServeMux
}

func NewProfilingHandler() *ProfilingHandler {
	r := http.NewServeMux()
	// Регистрация pprof-обработчиков
	r.HandleFunc("/pprof/", pprof.Index)
	r.HandleFunc("/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/pprof/profile", pprof.Profile)
	r.HandleFunc("/pprof/symbol", pprof.Symbol)
	r.HandleFunc("/pprof/trace", pprof.Trace)

	return &ProfilingHandler{r}
}

func (o *ProfilingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	o.r.ServeHTTP(w, r)
}
