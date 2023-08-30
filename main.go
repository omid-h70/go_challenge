package main

import (
	"context"
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hibiken/asynq"
	_ "github.com/lib/pq"
	"github.com/rakyll/statik/fs"
	"github.com/redis/go-redis/v9"
	"go_challenge/api"
	db "go_challenge/db/sqlc"
	_ "go_challenge/doc/statik" //you can replace it by go1.16 embed feature
	"go_challenge/gapi"
	"go_challenge/pb"
	"go_challenge/util"
	"go_challenge/worker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	zlog "github.com/rs/zerolog"
	"log"
	"net"
	"net/http"
	"os"
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

	if config.Env == "development" {
		//to have more human friendly logs
		log.Logger = log.Output(zerolog.consoleWriter({Out: os.Stderr})
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal.Msg(err.Error())
	}

	/* missing dependency */
	m, err := migrate.New(
		"./db/migration",
		config.DBSource)
	if err != nil {
		log.FatalMsg("cant create new migration")
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange{
		if err != nil {
		log.Fatal.Msg("cant run migration up command", err)
	}
	log.Println("db migrated successfully")

	store := db.NewStore(conn)
	server, _ := api.NewServer(&config, store)

	err = server.Start(config.HttpServerAddress)
	if err != nil {
		log.Fatal("can't start server", err)
	}

	redisOpt := async.RedisClientOpt{
		Addr: config.RedisAddr,
	}
	tskDistributr := worker.NewRedisTaskDistributor(redisOpt)

	runGrpcServer(config, store, tskDistributr)
	go runGatewayServer(config, store)
	go runTaskProcessor(redisOpt, store)
}

func runTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store){
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store)
	log.Info().Msg("start task processor")
	if err := taskProcessor.Start(); err != nil {
		log.Fatal().Err(err).Msg("failed to start processing redis tasks")
	}
}

func runGrpcServer(config util.Config, store db.Store, dstb worker.TaskDistributor) {
	server, err := gapi.NewServer(&config, store, dstb)
	if err != nil {
		log.Fatal("can't start grpc server")
	}

	grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger)

	pb.RegisterSimpleBankServer(grpcServer, server)
	//it allows user to explore self document
	reflection.Register(grpcServer)

	grpcListener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("can't start grpc listener")
	}
	log.Printf("start GRPC server on %s", config.GRPCServerAddress)

	err = grpcServer.Serve(grpcListener)
	if err != nil {
		log.Fatal("can't start server ", err)
	}
}

func runGatewayServer(config util.Config, store db.Store, dstb worker.TaskDistributor) {
	server, err := gapi.NewServer(&config, store, dstb)
	if err != nil {
		log.Fatal("can't start grpc server")
	}

	grpcMux := runtime.NewServeMux()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal("can't register grpc gateway")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	/* to Serve Manually hosted free swagger ?  v1
	fs := http.FileServer(http.Dir("./doc/swagger"))
	mux.Handle("/swagger/", http.StripPrefix("/swagger/", fs))
	*/

	statikFS, err := fs.New()
	if err != nil {
		log.Fatal("can't create statik FS")
	}

	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))
	mux.Handle("/swagger/", swaggerHandler)

	grpcListener, err := net.Listen("tcp", config.HttpServerAddress)
	if err != nil {
		log.Fatal("can't start grpc listener", err)
	}
	log.Printf("start Http server on %s", config.GRPCServerAddress)

	err = http.Serve(grpcListener, mux)
	if err != nil {
		log.Fatal("can't start server ", err)
	}
}

func runGinServer(config util.Config, store db.Store) {
	server, _ := api.NewServer(&config, store)

	err := server.Start(config.HttpServerAddress)
	if err != nil {
		log.Fatal("can't start server", err)
	}
}
