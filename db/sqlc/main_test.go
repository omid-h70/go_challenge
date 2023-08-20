package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"go_challenge/util"
	"log"
	"os"
	"testing"
)

// Entry Point Of All Uint Tests
// TestMain is Used For Testing a Package All together
var testQueries *Queries

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..") // Go For Current Path
	if err != nil {
		log.Fatal(err.Error())
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal(err.Error())
	}

	testQueries = New(conn)
	os.Exit(m.Run())
}
