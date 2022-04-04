package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq" //the underscore triggers the lexer to not remove the import if not used
	"github.com/rtpa25/ecomm-api-go/api"
	db "github.com/rtpa25/ecomm-api-go/db/sqlc"
	"github.com/rtpa25/ecomm-api-go/utils"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config vars", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal(err.Error(), "Cannot connect to database")
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("Could not loadup server", err)
	}

	err = server.Start(config.ServerAddress)

	if err != nil {
		log.Fatal("Cannot start server")
	}
}
