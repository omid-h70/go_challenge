package main

import (
	"github.com/gin-gonic/gin"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	//To Disable Extra Gin Logs here
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
