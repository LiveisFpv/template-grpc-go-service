package name_service_grpc

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type App struct {
	log        *logrus.Logger
	gRPCServer *grpc.Server
	port       int
}

// App constructor with logger and Service
func New(log *logrus.Logger, name_service_Service Name_service, port int) *App {

	recoverOpts := []recovery.Option{
		recovery.WithRecoveryHandler(
			func(p interface{}) (err error) {
				//Logging panic with leve error
				log.Errorf("Recovered from panic %s", p)
				//Return to client internal error
				return status.Errorf(codes.Internal, "internal error")
			}),
	}

	//logging all, logging data, loggind payload, user didnt know that we know
	loggingOpts := []grpclog.Option{
		grpclog.WithLogOnEvents(
			grpclog.PayloadReceived, grpclog.PayloadSent,
		),
	}
	//Create grpcServer with interseptors(logger, recover)
	gRPCServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		recovery.UnaryServerInterceptor(recoverOpts...),
		grpclog.UnaryServerInterceptor(InterceptorLogger(log), loggingOpts...),
	))

	//Тута и мне осознать надо
	Register(gRPCServer, name_service_Service)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func InterceptorLogger(l *logrus.Logger) grpclog.Logger {
	return grpclog.LoggerFunc(func(ctx context.Context, lvl grpclog.Level, msg string, fields ...any) {
		var logrusLevel logrus.Level
		switch lvl {
		case grpclog.LevelDebug:
			logrusLevel = logrus.DebugLevel
		case grpclog.LevelInfo:
			logrusLevel = logrus.InfoLevel
		case grpclog.LevelWarn:
			logrusLevel = logrus.WarnLevel
		case grpclog.LevelError:
			logrusLevel = logrus.ErrorLevel
		default:
			logrusLevel = logrus.InfoLevel
		}

		// Преобразование поля `fields` в структуру с ключами
		logFields := make(map[string]any)
		for i := 0; i < len(fields); i += 2 {
			if i+1 < len(fields) {
				key, ok := fields[i].(string)
				if !ok {
					key = "unknown"
				}
				logFields[key] = fields[i+1]
			}
		}

		// Логирование с полями
		l.WithFields(logrus.Fields{
			"details": logFields,
		}).Log(logrusLevel, msg)
	})
}

// Start gRPC server
func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

// Run runs gRPC server.
func (a *App) Run() error {
	const op = "grpcapp.Run"

	// Создаём listener, который будет слушить TCP-сообщения, адресованные
	// Нашему gRPC-серверу
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	a.log.Info("grpc server started", slog.String("addr", l.Addr().String()))

	// Запускаем обработчик gRPC-сообщений
	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// Stop gRPC server
func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.log.Infof("op %s", op)
	a.log.Info("stopping gRPC server", slog.Int("port", a.port))

	a.gRPCServer.GracefulStop()
}
