package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"template-grpc-go-service/internal/app"
	"template-grpc-go-service/internal/config"
	"template-grpc-go-service/internal/lib/logger"
)

func main() {
	cfg := config.MustLoad()

	//if need more opt make create in lib/logger
	log := logger.LoggerSetup(true)

	//init service and start
	ctx := context.Background()
	app := app.New(ctx, log, cfg.GRPC.Port, cfg.Dsn, cfg.TokenTTL)
	log.Info("Start service")
	go func() {
		app.GRPCServer.MustRun()
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)

	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	// Waiting for SIGINT (pkill -2) or SIGTERM
	<-stop
	// initiate graceful shutdown
	app.GRPCServer.Stop()
	log.Info("GRPCserver stopped")
	app.Storage.Stop()
	log.Info("Postgres connection closed")
}
