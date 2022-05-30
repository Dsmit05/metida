package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Dsmit05/metida/internal/api"
	"github.com/Dsmit05/metida/internal/config"
	"github.com/Dsmit05/metida/internal/cryptography"
	"github.com/Dsmit05/metida/internal/debag"
	"github.com/Dsmit05/metida/internal/logger"
	"github.com/Dsmit05/metida/internal/metrics"
	"github.com/Dsmit05/metida/internal/repositories"
	"github.com/Dsmit05/metida/internal/utils"
)

// @title metida
// @version 1.0.0
// @description API template server

// @host localhost:8080
// @basePath /api/v1/

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorizations
func main() {
	ctx := context.Background()

	// Init settings from cmd flag
	flagCmd, err := config.NewCommandLine()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Init logger
	if err = logger.InitLogger(flagCmd.IfDebagOn(), flagCmd.GetLogPath()); err != nil {
		panic(err)
	}

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recover(): %v", r)
			logger.Error("main()", err)
		}
		logger.ZapLog.Sync()
	}()

	// Init config app
	cfg, err := config.NewConfig(flagCmd)
	if err != nil {
		logger.Error("config.NewConfig()", err)
		return
	}

	logger.Info("config.NewConfig", cfg.String())

	// Init metrics
	metric := metrics.NewServiceMetrics()

	// Init connect to db
	db, err := repositories.NewPostgresRepository(cfg, metric)
	if err != nil {
		logger.Error("repositories.NewPostgresRepository()", err)
		return
	}
	defer db.Close()

	managerToken := cryptography.NewManagerToken(cfg.Cryptography.Secret)

	// Start api server
	apiServer := api.NewApiServer(db, managerToken, cfg, metric)
	go apiServer.Start()

	// Start debag server
	debagServer := debag.NewDebagServer(cfg)
	go debagServer.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	utils.Shutdown(ctx, apiServer, debagServer)
}
