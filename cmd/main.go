package main

import (
	"Workmate/internal/api"
	db "Workmate/internal/repositories/postgresql/sqlc"
	"Workmate/internal/repositories/redis"
	"Workmate/internal/service"
	"Workmate/util/config"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	Config, err := config.InitConfig(".")
	if err != nil {
		log.Fatalf("cannot load config %s", err)
	}
	connPool, err := pgxpool.New(context.Background(), Config.DBDSource)
	if err != nil {
		log.Fatalf("cannot connect to db %s", err)
	}
	queries := db.New(connPool)
	redisClient, err := redis.NewRedisClient(Config.RedisAddr)
	if err != nil {
		log.Fatalf("cannot connect to redis %v", err)
	}
	services := service.NewServices(queries, redisClient, Config)
	server := api.NewServer(router, services)
	runDBMigration(Config.MigrationUrl, Config.DBDSource)
	err = server.Run(Config.ServerAddress)
	if err != nil {
		log.Fatalf("cannot run server %v", err)
	}

}
func runDBMigration(migrationURL, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatalln("cannot find migration to up", err)
	}
	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalln("cannot start migration", err)
	}
}
