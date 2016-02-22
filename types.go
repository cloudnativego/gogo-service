package main

import "github.com/cloudnativego/gogo-engine"

type newMatchResponse struct {
	ID          string `json:"id"`
	StartedAt   int64  `json:"started_at"`
	GridSize    int    `json:"gridsize"`
	PlayerWhite string `json:"playerWhite"`
	PlayerBlack string `json:"playerBlack"`
	Turn        int    `json:"turn,omitempty"`
}

type matchDetailsResponse struct {
	ID          string   `json:"id"`
	StartedAt   int64    `json:"started_at"`
	GridSize    int      `json:"gridsize"`
	PlayerWhite string   `json:"playerWhite"`
	PlayerBlack string   `json:"playerBlack"`
	Turn        int      `json:"turn,omitempty"`
	GameBoard   [][]byte `json:"gameboard"`
}

type newMatchRequest struct {
	GridSize    int    `json:"gridsize"`
	PlayerWhite string `json:"playerWhite"`
	PlayerBlack string `json:"playerBlack"`
}

type boardPosition struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type newMoveRequest struct {
	Player   byte          `json:"player"`
	Position boardPosition `json:"position"`
}

type matchRepository interface {
	addMatch(match gogo.Match) (err error)
	getMatches() []gogo.Match
	getMatch(id string) (match gogo.Match, err error)
	updateMatch(id string, match gogo.Match) (err error)
}

func (request newMatchRequest) isValid() (valid bool) {
	valid = true
	if request.GridSize != 19 && request.GridSize != 13 && request.GridSize != 9 {
		valid = false
	}
	if request.PlayerWhite == "" {
		valid = false
	}
	if request.PlayerBlack == "" {
		valid = false
	}
	return valid
}
