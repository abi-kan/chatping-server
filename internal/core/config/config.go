package core_config

import (
	core_pgx "chatping-server/internal/core/adapter/postgres/pgx"
	core_redis "chatping-server/internal/core/adapter/redis"
	core_logger "chatping-server/internal/core/logger"
	"fmt"
)

type Config struct {
	App      App
	Postgres core_pgx.Config
	Redis    core_redis.Config
	Logger   core_logger.Config
}

func NewConfig() Config {
	appConfig, err := newAppConfig()
	if err != nil {
		err = fmt.Errorf("get app config: %w", err)
		panic(err)
	}

	postgresConfig, err := core_pgx.NewConfig()
	if err != nil {
		err = fmt.Errorf("get postgress config: %w", err)
		panic(err)
	}

	redisConfig, err := core_redis.NewConfig()
	if err != nil {
		err = fmt.Errorf("get redis config: %w", err)
		panic(err)
	}

	loggerConfig, err := core_logger.NewConfig()
	if err != nil {
		err = fmt.Errorf("get logger config: %w", err)
		panic(err)
	}

	config := Config{
		App:      appConfig,
		Postgres: postgresConfig,
		Redis:    redisConfig,
		Logger:   loggerConfig,
	}

	return config
}
