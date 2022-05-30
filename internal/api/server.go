package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Dsmit05/metida/internal/logger"
	"github.com/Dsmit05/metida/internal/utils"
	"github.com/rs/cors"
)

const maxHeaderBytes = 1 << 20 // controls the maximum number of bytes in header

type ApiServer struct {
	s *http.Server
}

func NewApiServer(
	db repositoryI,
	managerToken cryptographyI,
	cfg configApiI,
	metric metricI) *ApiServer {

	ginBuilder := NewGinBuilder(db, managerToken, cfg).AddV1("/api/v1")

	serveMux := utils.RouterComposition(utils.Hanlde{
		Pattern: "/api/",
		Handler: ginBuilder,
	})

	serveMux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/" {
			http.NotFound(w, req)
			return
		}
		fmt.Fprintf(w, "If you have any questions, write to the mail example@google.com")
	})

	metricMiddl := metric.MetricsMiddleware(serveMux)

	corsProvided := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://" + cfg.GetDebagAddr()},
		AllowedHeaders:   []string{"Authorizations", "Content-Type"},
		AllowCredentials: true,
		Debug:            cfg.IfDebagOn(),
	})

	corsMiddl := corsProvided.Handler(metricMiddl)

	s := &http.Server{
		Addr:           cfg.GetApiAddr(),
		Handler:        corsMiddl,
		ReadTimeout:    cfg.GetApiReadTimeout(),
		WriteTimeout:   cfg.GetApiWriteTimeout(),
		MaxHeaderBytes: maxHeaderBytes,
	}

	return &ApiServer{s}
}

func (o *ApiServer) Start() {
	if err := o.s.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			logger.Info("ApiServer", "Server closed under request")
		} else {
			logger.Error("ApiServer closed unexpect", err)
		}
	}
}

func (o *ApiServer) Stop(ctx context.Context) {
	stop := make(chan bool)

	go func() {
		o.s.SetKeepAlivesEnabled(false)
		if err := o.s.Shutdown(ctx); err != nil {
			logger.Error("ApiServer Shutdown failed", err)
		}
		stop <- true
	}()

	select {
	case <-ctx.Done():
		logger.Error("ApiServer context timeout", ctx.Err())
	case <-stop:
		logger.Info("ApiServer", "Stop")
	}

}
