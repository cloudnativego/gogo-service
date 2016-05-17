package service

import (
	"testing"

	"github.com/cloudnativego/gogo-engine"
)

func TestAddMatchShowsUpInRepository(t *testing.T) {
	match := gogo.NewMatch(19, "bob", "alfred")

	repo := newInMemoryRepository()
	err := repo.addMatch(match)
	if err != nil {
		t.Error("Got an error adding a match to repository, should not have.")
	}

	matches, err := repo.getMatches()
	if err != nil {
		t.Errorf("Unexpected error in getMatches(): %s", err)
	}
	if len(matches) != 1 {
		t.Errorf("Expected to have 1 match in the repository, got %d", len(matches))
	}

	if matches[0].PlayerBlack != "bob" {
		t.Errorf("Player 1's name should have been bob, got %s", matches[0].PlayerBlack)
	}
	if matches[0].PlayerWhite != "alfred" {
		t.Errorf("Player 2's name should have been alfred, got %s", matches[0].PlayerWhite)
	}
}

func TestGetMatchRetrievesProperMatch(t *testing.T) {
	match := gogo.NewMatch(19, "bob", "alfred")

	repo := newInMemoryRepository()
	err := repo.addMatch(match)
	if err != nil {
		t.Error("Got an error adding a match to repository, should not have.")
	}

	target, err := repo.getMatch(match.ID)
	if err != nil {
		t.Errorf("Got an error when retrieving match from repo instead of success. Err: %s", err.Error())
	}

	if target.GridSize != 19 {
		t.Errorf("Got the wrong gridsize. Expected 19, got %d", target.GridSize)
	}
}

func TestNewRepositoryIsEmpty(t *testing.T) {
	repo := newInMemoryRepository()

	matches, err := repo.getMatches()
	if err != nil {
		t.Errorf("Unexpected error in getMatches(): %s", err)
	}
	if len(matches) != 0 {
		t.Errorf("Expected to have 0 matches in newly created repository, got %d", len(matches))
	}
}

func TestUpdateMatch(t *testing.T) {
	redHerring := gogo.NewMatch(13, "buckshank", "d'squarius")
	match := gogo.NewMatch(19, "bob", "alfred")

	repo := newInMemoryRepository()
	err := repo.addMatch(redHerring)
	if err != nil {
		t.Errorf("Error adding match: %s", err)
	}
	err = repo.addMatch(match)
	if err != nil {
		t.Errorf("Error adding match: %s", err)
	}

	match.TurnCount = 37
	err = repo.updateMatch(match.ID, match)
	if err != nil {
		t.Errorf("Error updating match: %s", err)
	}

	found, err := repo.getMatch(match.ID)
	if err != nil {
		t.Errorf("Error retrieving updated match: %s", err)
	}
	if found.GridSize != match.GridSize || found.PlayerWhite != match.PlayerWhite {
		t.Errorf("Retrieved incorrect match:\nexpected %+v\nreceived %+v", match, found)
	}
	if found.TurnCount != 37 {
		t.Errorf("Update failed: expected %d; received %d", 37, found.TurnCount)
	}
}
