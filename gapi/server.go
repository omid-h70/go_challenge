package gapi

import (
	"fmt"
	db "go_challenge/db/sqlc"
	"go_challenge/pb"
	"go_challenge/token"
	"go_challenge/util"
	"go_challenge/worker"
)

type Server struct {
	pb.UnimplementedSimpleBankServer
	store           db.Store //### add it later
	tokenMaker      token.Maker
	config          *util.Config
	taskDistributor worker.TaskDistributor
}

func NewServer(config *util.Config, store db.Store, distributor worker.TaskDistributor) (*Server, error) {

	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannt create token %w", err)
	}

	server := &Server{
		store:           store,
		tokenMaker:      tokenMaker,
		config:          config,
		taskDistributor: distributor,
	}

	return server, nil
}
