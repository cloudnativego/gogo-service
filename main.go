package main

import (
	"os"

	"github.com/cloudfoundry-community/go-cfenv"
)

func main() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "3000"
	}

	appEnv, _ := cfenv.Current()
	server := NewServer(appEnv)
	server.Run(":" + port)
}
