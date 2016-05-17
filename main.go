package main

import (
	"os"

	"github.com/cloudfoundry-community/go-cfenv"
	service "github.com/cloudnativego/gogo-service/service"
)

func main() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "3000"
	}

	appEnv, _ := cfenv.Current()
	server := service.NewServer(appEnv)
	server.Run(":" + port)
}
