package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang/mock/mockgen/model"
	"github.com/hibiken/asynq"
	_ "github.com/lib/pq"
	logZe "github.com/rs/zerolog/log"
	"github.com/techschool/simplebank/api"
	db "github.com/techschool/simplebank/db/sqlc"
	"github.com/techschool/simplebank/util"
	"github.com/techschool/simplebank/worker"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	conn, err := sql.Open(config.DbDriver, config.DbSource)
	if err != nil {
		log.Fatal("Err wheifn open db: ", err)
	}

	//run db migration
	runDbMigratio(config.MigrationSource, config.DbSource)

	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}
	taskDistributor := worker.NewRedisTaskDistributor(&redisOpt)

	store := db.NewStore(conn)
	go runTaskProcessor(redisOpt, store)
	server, err := api.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal("Server cannot be started: ", err)
		return
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("Server cannot be started: ", err)
	}
}

func runTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store) {
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store)
	logZe.Info().Msg("starting task processor")
	err := taskProcessor.Start()
	if err != nil {
		log.Fatal("error when start task processor ", err)
	}
	logZe.Info().Msg("started task processor")

}

func runDbMigratio(migrationDir string, databaseSource string) {
	migration, err := migrate.New(migrationDir, databaseSource)
	if err != nil {
		log.Fatal("cannot create migration: ", err)
	}
	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("cannot migrate dtb: ", err)
	}
	fmt.Println("migrate successfully")
}
