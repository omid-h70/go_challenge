package gapi

import (
	"github.com/stretchr/testify/require"
	db "go_challenge/db/sqlc"
	"go_challenge/util"
	"go_challenge/worker"
	"testing"
	"time"
)

func newTestServer(t *testing.T, store db.Store, distributor worker.TaskDistributor) *Server {

	cfg := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(&cfg, store, distributor)
	require.NoError(t, err)
	require.NotNil(t, server)
	return server
}
