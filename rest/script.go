package rest

import (
	"go_challenge/rest/handlers"
	"go_challenge/rest/models"
)

func prepareRestDB() {
	var db models.PostgresDBStruct
	postgresInstance := db.GetInstance()
	//defer db.CloseDB()

	postgresInstance.CreateTables()
	postgresInstance.InsertDBSeeds()
	//postgresInstance.InsertDummySeeds(5)
}

func RunRestApp() {
	prepareRestDB()
	handlers.StartServer()
}
