package main

import (
	"codeup.aliyun.com/6145b2b428003bdc3daa97c8/go-simba/go-simba-pkg.git/logger"
	"codeup.aliyun.com/6145b2b428003bdc3daa97c8/operation-platform/op-mno.git/pkg/interceptor"
	"codeup.aliyun.com/6145b2b428003bdc3daa97c8/operation-platform/op-mno.git/src/app/config"
	"codeup.aliyun.com/6145b2b428003bdc3daa97c8/operation-platform/op-mno.git/src/app/inject"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/automaxprocs/maxprocs"
)

func main() {

	// set GOMAXPROCS
	_, err := maxprocs.Set()
	if err != nil {
		slog.Error("failed set max procs", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	cfg, err := config.NewConfig()
	if err != nil {
		slog.Error("failed get config", err)
	}

	slog.Info("⚡ init app", "name", cfg.Name, "version", cfg.Version)

	// set up logrus
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logger.ConvertLogLevel(cfg.Log.Level))

	// integrate Logrus with the slog logger
	slog.New(logger.NewLogrusHandler(logrus.StandardLogger()))

	server := grpc.NewServer(grpc.UnaryInterceptor(interceptor.UnaryServerInterceptor()), grpc.StreamInterceptor(interceptor.StreamServerInterceptor()))

	go func() {
		defer server.GracefulStop()
		<-ctx.Done()
	}()
	cleanup := prepareApp(ctx, cancel, cfg, server)
	// gRPC Server.
	address := fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port)
	network := "tcp"

	l, err := net.Listen(network, address)
	if err != nil {
		slog.Error("failed to listen to address", err, "network", network, "address", address)
		cancel()
	}

	slog.Info("🌏 start server...", "address", address)

	defer func() {
		if err1 := l.Close(); err != nil {
			slog.Error("failed to close", err1, "network", network, "address", address)
		}
	}()

	err = server.Serve(l)
	if err != nil {
		slog.Error("failed start gRPC server", err, "network", network, "address", address)
		cancel()
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case v := <-quit:
		cleanup()
		slog.Info("signal.Notify", v)
	case done := <-ctx.Done():
		cleanup()
		slog.Info("ctx.Done", done)
	}
}

func prepareApp(ctx context.Context, cancel context.CancelFunc, cfg *config.Config, server *grpc.Server) func() {
	_, cleanup, err := inject.InitApp(cfg, server)
	reflection.Register(server)
	if err != nil {
		slog.Error("failed init app", err)
		cancel()
		<-ctx.Done()
	}
	return cleanup
}
