package integrations_test

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/cloudnativego/cfmgo"

	. "github.com/cloudnativego/gogo-service"
)

func getRepo(collectionName string) (col cfmgo.Collection, err error) {
	host := os.Getenv("MONGO_PORT_27017_TCP_ADDR")
	port := os.Getenv("MONGO_PORT_27017_TCP_PORT")

	if len(host) == 0 {
		err = errors.New("Could not retrieve mongo host information.")
		return
	}
	if len(port) == 0 {
		port = "27017"
	}

	uri := fmt.Sprintf("mongodb://%s:%s/fake-guid", host, port)
	col = cfmgo.Connect(cfmgo.NewCollectionDialer, uri, collectionName)

	return
}

func TestAddMatchToRepo(t *testing.T) {
	matches, err := getRepo("matches")
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	repo := NewMongoMatchRepository(matches)
	if repo == nil {
		t.Error("Error creating Mongo repo")
	}
}
