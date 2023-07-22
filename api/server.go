package api

import (
	"bufio"
	"fmt"
	"github.com/gin-gonic/gin"
	"go_challenge/cmd"
	"go_challenge/cmd/handlers"
	"go_challenge/cmd/models"
	"os"
	"strconv"
	"strings"
)

type Server struct {
	//store *db.store ### add it later
	router *gin.Engine
}

func NewServer() *Server {
	server := &Server{
		router: gin.Default(),
	}

	server.router.POST("/accounts", server.createAccount)
	server.router.GET("/accounts/:id", server.getAccount)
	return server
}

func (server *Server) Start(addr string) error {
	return server.router.Run(addr)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func insertDummySeedsToDB(count int) {
	var db models.PostgresDBStruct
	postgresInstance := db.GetInstance()
	//defer db.CloseDB()

	postgresInstance.CreateTables()
	postgresInstance.InsertDummySeeds(count)
}

func RunCli() {
	var db models.PostgresDBStruct
	defer db.CloseDB()

	fmt.Println("\n\n\n\n\n")
	fmt.Println("Welcome to go_challenge mini app")
	fmt.Println("You Must Have Docker installed for running this app")
	fmt.Println("First Please run docker-compose up in {$PWD}/cmd/models path to have your postgres ready")
	fmt.Println("Press 1 or 2 or q to quit")
	fmt.Println("1. To Start localhost on port 8000")
	fmt.Println("2. To Cli Dummy Generator")

	reader := bufio.NewReader(os.Stdin)

	for {
		inStr, _ := reader.ReadString('\n')
		inStr = strings.TrimRight(inStr, "\r\n")
		if string(inStr) == "1" {
			cmd.RunRestApp()
			break
		} else if string(inStr) == "2" {

			fmt.Println("#########################################################")
			fmt.Println("##############  Cli Dummy Generator #####################")
			fmt.Println("#########################################################")
			fmt.Println("Please Enter Valid INT number")

			for {
				inStr, _ = reader.ReadString('\n')
				inStr = strings.TrimRight(inStr, "\r\n")

				if count, err := strconv.Atoi(inStr); err == nil {
					fmt.Println("You Can Check Data by REST verbs")
					insertDummySeedsToDB(count)
					handlers.StartServer()
					break
				} else if string(inStr) == "q" || string(inStr) == "Q" {
					break
				} else {
					fmt.Println("Error  => Please Enter Valid INT number")
				}
			}
			break
		} else if string(inStr) == "q" || string(inStr) == "Q" {
			fmt.Println("Done ...!")
			break
		}
	}
}
