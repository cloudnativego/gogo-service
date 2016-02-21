package integrations_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/cloudnativego/cfmgo"

	. "github.com/cloudnativego/gogo-service"
)

func TestAddMatchToRepo(t *testing.T) {
	mongoHost := os.Getenv("MONGO_PORT_27017_TCP_ADDR")
	mongoPort := os.Getenv("MONGO_PORT_27017_TCP_PORT")

	if len(mongoHost) == 0 {
		t.Fatal("Could not retrieve mongo host information")
	}
	if len(mongoPort) == 0 {
		mongoPort = "27017"
	}

	mongoURI := fmt.Sprintf("mongodb://%s:%s/fake-guid", mongoHost, mongoPort)
	fmt.Println("MONGOURI:\t" + mongoURI)
	matchesCollection := cfmgo.Connect(cfmgo.NewCollectionDialer, mongoURI, "matches")

	repo := NewMongoMatchRepository(matchesCollection)
	if repo == nil {
		t.Error("Error creating Mongo repo")
	}
}
