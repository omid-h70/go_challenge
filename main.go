package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"go_challenge/api"
	"go_challenge/cmd"
	db "go_challenge/db/sqlc"
	"go_challenge/util"
	"log"
)

/*
	They will be loaded from Config File

const (

	dbDriver   = "postgres"
	dbSource   = "postgresql://root:secret@localhost:5432/test_db?sslmode=disable"
	serverAddr = "0.0.0.0:8080"

)
*/
func main() {

	config, err := util.LoadConfig(".") // Go For Current Path
	if err != nil {
		log.Fatal(err.Error())
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal(err.Error())
	}

	store := db.NewStore(conn)
	server, _ := api.NewServer(&config, store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("can't start server")
	}
}

func _main() {
	cmd.RunRestApp()
	return
}
