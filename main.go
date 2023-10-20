package main

import (
	"context"
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hibiken/asynq"
	"github.com/rakyll/statik/fs"
	zlog "github.com/rs/zerolog/log"
	"go_challenge/api"
	db "go_challenge/db/sqlc"
	"go_challenge/gapi"
	"go_challenge/mail"
	"go_challenge/pb"
	"go_challenge/util"
	"go_challenge/worker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
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
		//zlog.Logger = zlog.Output(log.consoleWriter({Out: os.Stderr})
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		zlog.Fatal().Msg(err.Error())
	}

	/* missing dependency */
	m, err := migrate.New(
		"./db/migration",
		config.DBSource)
	if err != nil {
		zlog.Fatal().Msg("cant create new migration")
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		if err != nil {
			zlog.Fatal().Err(err).Msg("cant run migration up command")
		}
		log.Println("db migrated successfully")

		store := db.NewStore(conn)
		server, _ := api.NewServer(&config, store)

		err = server.Start(config.HttpServerAddress)
		if err != nil {
			log.Fatal("can't start server", err)
		}

		redisOpt := asynq.RedisClientOpt{
			Addr: config.RedisAddr,
		}
		tskDistributor := worker.NewRedisTaskDistributor(redisOpt)
		gmail := mail.NewGmailSender("", "", "")

		runGrpcServer(config, store, tskDistributor)
		go runGatewayServer(config, store, tskDistributor)
		go runTaskProcessor(redisOpt, store, gmail)
	}
}

func runTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store, mailer mail.EmailSender) {
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store, mailer)
	zlog.Info().Msg("start task processor")
	if err := taskProcessor.Start(); err != nil {
		zlog.Fatal().Err(err).Msg("failed to start processing redis tasks")
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
