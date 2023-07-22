package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"go_challenge/api"
	"go_challenge/cmd"
	"go_challenge/util"
	"log"
)

const (
	dbDriver   = "postgres"
	dbSource   = "postgresql://root:secret@localhost:5432/test_db?sslmode=disabled"
	serverAddr = "0.0.0.0:8080"
)

func main() {

	config, err := util.LoadConfig(".") // Go For Current Path
	if err != nil {
		log.Fatal(err.Error())
	}

	_, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal(err.Error())
	}

	//store := db.NewStore(conn)
	server := api.NewServer( /*store*/)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("can't start server")
	}
}

func _main() {
	cmd.RunRestApp()
	return
}
