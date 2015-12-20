package main

import "github.com/cloudnativego/gogo-engine"

type inMemoryMatchRepository struct {
	matches []gogo.Match
}

// NewRepository creates a new in-memory match repository
func newInMemoryRepository() *inMemoryMatchRepository {
	repo := &inMemoryMatchRepository{}
	repo.matches = []gogo.Match{}
	return repo
}

func (repo *inMemoryMatchRepository) addMatch(match gogo.Match) (err error) {
	repo.matches = append(repo.matches, match)
	return err
}

func (repo *inMemoryMatchRepository) getMatches() []gogo.Match {
	return repo.matches
}
