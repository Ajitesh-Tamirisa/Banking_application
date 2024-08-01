package main

import (
	"banking_application/api"
	db "banking_application/db/sqlc"
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:password@localhost:5432/banking_app?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Cannot connect to Database:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("Cannot start server: ", err)
	}
}
