package app

import (
	"context"
	name_service_grpc "template-grpc-go-service/internal/grpc"
	"template-grpc-go-service/internal/services/name_service"
	"template-grpc-go-service/internal/storage"
	"time"

	"github.com/sirupsen/logrus"
)

type App struct {
	GRPCServer *name_service_grpc.App
	Storage    storage.Repository
}

// Constructor APP creates gRPCServer, storage
func New(
	ctx context.Context,
	log *logrus.Logger,
	grpcPort int,
	dsn string,
	tokenTTL time.Duration,
) *App {
	storage, err := storage.NewStorage(ctx, dsn, log)
	if err != nil {
		panic(err)
	}

	//Создание сервисов и подключения
	name_service_Service := name_service.New(log, storage, tokenTTL)
	grpcApp := name_service_grpc.New(log, name_service_Service, grpcPort)

	return &App{
		GRPCServer: grpcApp,
		Storage:    storage,
	}
}
