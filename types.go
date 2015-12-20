package main

import "github.com/cloudnativego/gogo-engine"

type newMatchResponse struct {
	ID        string   `json:"id"`
	StartedAt int64    `json:"started_at"`
	GridSize  int      `json:"gridsize"`
	Players   []player `json:"players"`
}

type player struct {
	Color string `json:"color"`
	Name  string `json:"name"`
}

type newMatchRequest struct {
	GridSize int      `json:"gridsize"`
	Players  []player `json:"players"`
}

type matchRepository interface {
	addMatch(match gogo.Match) (err error)
	getMatches() []gogo.Match
}
