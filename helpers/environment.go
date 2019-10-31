package helpers

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"os"
)

var (
	Production  bool = false
	Development bool = false
	Testing     bool = false
)

func setEnvironment() {
	if TestingCmd() {
		Testing = true
		return
	}
	if devEnvar() {
		log.SetLevel(log.DebugLevel)
		Development = true
		return
	}
	Production = true
	log.SetFormatter(&log.JSONFormatter{})
	gin.SetMode(gin.ReleaseMode)
}

func env() string {
	val := os.Getenv("WAVE_ENV")
	if val == "" {
		return "production"
	}
	return val
}

func devEnvar() bool {
	return env() == "development"
}

// Determine if the command was generated by `go test`
func TestingCmd() bool {
	if len(os.Args) < 1 {
		return false
	}

	cmd := os.Args[0]

	if len(cmd) < 4 {
		return false
	}

	return cmd[len(cmd)-4:] == "test"
}
