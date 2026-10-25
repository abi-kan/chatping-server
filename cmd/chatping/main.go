package main

import (
	core_postgres "chatping-server/internal/core/adapter/postgres"
	core_pgx "chatping-server/internal/core/adapter/postgres/pgx"
	core_redis "chatping-server/internal/core/adapter/redis"
	core_config "chatping-server/internal/core/config"
	core_logger "chatping-server/internal/core/logger"
	core_http_middleware "chatping-server/internal/core/transport/http/middleware"
	core_http_server "chatping-server/internal/core/transport/http/server"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

// @title Golang Chatping API
// @version 1.0
// @description Chatping Application REST-API scheme
// @host 127.0.0.1:5051
// @BasePath /api/v1
func main() {
	config := core_config.NewConfig()
	time.Local = config.App.TimeZone

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGALRM,
	)
	defer cancel()

	logger, err := core_logger.NewLogger(config.Logger)
	if err != nil {
		fmt.Println("failed to init application logger:", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("application time zone", zap.Any("zone", time.Local))
	logger.Debug("initializing postgres connection pool")
	pool, err := core_pgx.NewPool(
		ctx,
		config.Postgres,
	)
	if err != nil {
		logger.Fatal("failed to init postgres connection pool", zap.Error(err))
	}
	defer pool.Close()

	logger.Debug("initializing redis client")
	redis, err := core_redis.NewClient(ctx, config.Redis)
	if err != nil {
		logger.Fatal("failed to init redis client", zap.Error(err))
	}
	defer redis.Close()

	httpServer := createServer(logger, pool, redis, config)
	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run error", zap.Error(err))
	}
}

func createServer(
	logger *core_logger.Logger,
	pool core_postgres.Pool,
	redis *core_redis.Client,
	config core_config.Config,
) *core_http_server.Server {
	logger.Debug("initializing HTTP server")
	httpServer := core_http_server.NewServer(
		config.HTTP,
		logger,
		core_http_middleware.CORS(),
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)

	return httpServer
}
