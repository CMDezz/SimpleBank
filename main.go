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
	_ "github.com/lib/pq"
	"github.com/techschool/simplebank/api"
	db "github.com/techschool/simplebank/db/sqlc"
	"github.com/techschool/simplebank/util"
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

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("Server cannot be started: ", err)
		return
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("Server cannot be started: ", err)
	}
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
