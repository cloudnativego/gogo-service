package integrations_test

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/cloudnativego/cfmgo"
	"github.com/cloudnativego/gogo-engine"

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
	matchesCollection, err := getRepo("matches")
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	repo := NewMongoMatchRepository(matchesCollection)
	match := gogo.NewMatch(19, "buckshank", "d'squarius")
	err = repo.AddMatch(match)
	if err != nil {
		t.Errorf("Error adding match to mongo: %v", err)
	}

	matches, err := repo.GetMatches()
	if err != nil {
		t.Errorf("Error retrieving matches: %v", err)
	}
	if len(matches) == 0 {
		t.Errorf("Expected matches length to be greater than 0; received %d", len(matches))
	}

	foundMatch, err := repo.GetMatch(match.ID)
	if err != nil {
		t.Errorf("Unable to find match with ID: %v... %s", match.ID, err)
	}
	if foundMatch.GridSize != match.GridSize || foundMatch.PlayerBlack != match.PlayerBlack {
		t.Errorf("Unexpected match results: %v", foundMatch)
	}
}
