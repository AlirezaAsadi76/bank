package main

import (
	"database/sql"
	"firstprj/api"
	db "firstprj/db/sqlc"
	"firstprj/envi"
	_ "github.com/lib/pq"
	"log"
)

func main() {

	config, err := envi.LoadConfig(".")
	if err != nil {
		log.Fatal("envi can not load :", err)
	}
	conn, err := sql.Open(config.DbDriver, config.DbPath)

	if err != nil {
		log.Fatal("db can not connection")
	}
	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("server can not connection : ", err)
	}
}
