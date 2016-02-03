package main

import (
	"testing"

	"github.com/cloudnativego/gogo-engine"
	"github.com/cloudnativego/gogo-service/fakes"
	"github.com/pivotal-pez/cfmgo"
)

var (
	fakeDBURI = "mongodb://fake.uri@addr:port/guid"
)

func TestAddMatchShowsUpInMongoRepository(t *testing.T) {
	var fakeMatches = []matchRecord{}
	var matchesCollection = cfmgo.Connect(
		fakes.FakeNewCollectionDialer(fakeMatches),
		fakeDBURI,
		MatchesCollectionName)

	match := gogo.NewMatch(19, "bob", "alfred")
	repo := newMongoMatchRepository(matchesCollection)
	err := repo.addMatch(match)
	if err != nil {
		t.Error("Got an error adding a match to mongo, should not have.")
	}

	matches, err := repo.getMatches()
	if err != nil {
		t.Errorf("Got an error retrieving matches: %v", err)
	}
	if len(matches) == 0 {
		t.Errorf("Expected matches length to be greater than 0")
	}
}

// func TestAddMatchToPopulateMongoRepository(t *testing.T) {
// 	// make sure no destroy of existing
// }

// func TestGetMatchRetrievesProperMatchFromMongo(t *testing.T) {

// }

// func TestMatchUpdateShowsInMatchListAndMatchDetails(t *testing.T) {
// 	// create match; puts in repo
// 	// modify turn count and game board; then call update
// 	// verify match details reflect new state
// }
