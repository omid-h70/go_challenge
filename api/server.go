package api

import (
	"bufio"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"go_challenge/_cmd"
	"go_challenge/_cmd/handlers"
	"go_challenge/_cmd/models"
	db "go_challenge/db/sqlc"
	"go_challenge/token"
	"go_challenge/util"
	"os"
	"strconv"
	"strings"
)

type Server struct {
	store      db.Store //### add it later
	router     *gin.Engine
	tokenMaker token.Maker
	config     *util.Config
}

func NewServer(config *util.Config, store db.Store) (*Server, error) {

	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannt create token %w", err)
	}

	server := &Server{
		store:      store,
		router:     gin.Default(),
		tokenMaker: tokenMaker,
		config:     config,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {

	//First routes Doesn't Need it
	server.router.POST("/users", server.createUser)
	server.router.POST("/users/login", server.loginUser)
	server.router.POST("/users/renew_access", server.renewAccessToken)

	//Below Handlers use MiddleWares
	authRouteGroup := server.router.Group("/").Use(authMiddleware(server.tokenMaker))
	authRouteGroup.POST("/accounts", server.createAccount)
	authRouteGroup.GET("/accounts/:id", server.getAccount)

	authRouteGroup.POST("/transfers", server.createTransfer)
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
	fmt.Println("First Please run docker-compose up in {$PWD}/_cmd/models path to have your postgres ready")
	fmt.Println("Press 1 or 2 or q to quit")
	fmt.Println("1. To Start localhost on port 8000")
	fmt.Println("2. To Cli Dummy Generator")

	reader := bufio.NewReader(os.Stdin)

	for {
		inStr, _ := reader.ReadString('\n')
		inStr = strings.TrimRight(inStr, "\r\n")
		if string(inStr) == "1" {
			_cmd.RunRestApp()
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
