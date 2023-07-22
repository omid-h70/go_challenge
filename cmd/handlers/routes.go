package handlers

import (
	"encoding/json"
	"fmt"
	"go_challenge/cmd/models"
	"log"
	"net/http"
	"strings"
)

const serverAddr string = "localhost:8000"

func StartServer() {
	mux := mux.NewRouter()

	//wiring Application Parts
	fmt.Println("Server is Listening On " + serverAddr)
	mux.HandleFunc("/request", getRecords)
	mux.HandleFunc("/report", getTradeReport)
	log.Fatal(http.ListenAndServe("localhost:8000", mux))

	fmt.Println("Urls You can reach By GET And json Content-Type")
	fmt.Println("http://localhost:8000/request?table=instrument | trade")
	fmt.Println("http://localhost:8000/report?name=GOOGL | AAPL")
}

func getAllTradeRecords(w http.ResponseWriter, r *http.Request) {
	var db models.PostgresDBStruct
	postgresInstance := db.GetInstance()
	//it sends to Response Page To Write
	if r.Header.Get("Content-Type") == "application/json" {
		w.Header().Add("Content-Type", "application/json")
		if r.URL.Query().Get("table") == "trade" {
			jsonLog := postgresInstance.GetDBLog(models.E_TARDE_TABLE)
			w.Write([]byte(jsonLog))
		} else if r.URL.Query().Get("table") == "instrument" {
			jsonLog := postgresInstance.GetDBLog(models.E_INSTROMENT_TABLE)
			w.Write([]byte(jsonLog))
		}
	}
}

func getRecords(w http.ResponseWriter, r *http.Request) {
	//it sends to Response Page To Write
	getAllTradeRecords(w, r)
}

func getTradeReport(w http.ResponseWriter, r *http.Request) {
	var db models.PostgresDBStruct
	postgresInstance := db.GetInstance()
	//it sends to Response Page To Write
	if r.Header.Get("Content-Type") == "application/json" {
		w.Header().Add("Content-Type", "application/json")

		stReport := postgresInstance.GetTradeReport(strings.ToUpper(r.URL.Query().Get("name")))
		jsonLog, _ := json.Marshal(stReport)
		w.Write(jsonLog)
	}
}
