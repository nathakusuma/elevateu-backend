package main

import (
	"github.com/nathakusuma/elevateu-backend/internal/infra/database"
	"github.com/nathakusuma/elevateu-backend/internal/infra/env"
	"github.com/nathakusuma/elevateu-backend/internal/infra/redis"
	"github.com/nathakusuma/elevateu-backend/internal/infra/server"
	"github.com/nathakusuma/elevateu-backend/pkg/log"
)

func main() {
	env.NewEnv()
	log.NewLogger()

	srv := server.NewHTTPServer()
	postgresDB := database.NewPostgresPool(
		env.GetEnv().DBHost,
		env.GetEnv().DBPort,
		env.GetEnv().DBUser,
		env.GetEnv().DBPass,
		env.GetEnv().DBName,
		env.GetEnv().DBSSLMode,
	)
	redisClient := redis.NewRedisPool(
		env.GetEnv().RedisHost,
		env.GetEnv().RedisPort,
		env.GetEnv().RedisPass,
		env.GetEnv().RedisDB,
	)
	defer postgresDB.Close()
	defer redisClient.Close()

	srv.MountMiddlewares()
	srv.MountRoutes(postgresDB, redisClient)
	srv.Start(env.GetEnv().AppPort)
}
