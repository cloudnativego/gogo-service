package main

import (
	"testing"

	"github.com/cloudnativego/gogo-engine"
)

func TestAddMatchShowsUpInRepository(t *testing.T) {
	match := gogo.NewMatch(19)

	repo := newInMemoryRepository()
	err := repo.addMatch(match)
	if err != nil {
		t.Error("Got an error adding a match to repository, should not have.")
	}

	matches := repo.getMatches()
	if len(matches) != 1 {
		t.Errorf("Expected to have 1 match in the repository, got %d", len(matches))
	}
}

func TestNewRepositoryIsEmpty(t *testing.T) {
	repo := newInMemoryRepository()

	matches := repo.getMatches()
	if len(matches) != 0 {
		t.Errorf("Expected to have 0 matches in newly created repository, got %d", len(matches))
	}
}
