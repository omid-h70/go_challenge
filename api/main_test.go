package api

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	db "go_challenge/db/sqlc"
	"go_challenge/util"
	"os"
	"testing"
	"time"
)

func newTestServer(t *testing.T, store *db.Store) *Server {

	cfg := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(&cfg, store)
	require.NoError(t, err)
	require.NotNil(t, server)
	return server
}

func TestMain(m *testing.M) {
	//To Disable Extra Gin Logs here
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
