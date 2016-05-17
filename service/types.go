package service

import "github.com/cloudnativego/gogo-engine"

type newMatchResponse struct {
	ID          string `json:"id"`
	StartedAt   int64  `json:"started_at"`
	GridSize    int    `json:"gridsize"`
	PlayerWhite string `json:"playerWhite"`
	PlayerBlack string `json:"playerBlack"`
	Turn        int    `json:"turn,omitempty"`
}

func (m *newMatchResponse) copyMatch(match gogo.Match) {
	m.ID = match.ID
	m.StartedAt = match.StartTime.Unix()
	m.GridSize = match.GridSize
	m.PlayerWhite = match.PlayerWhite
	m.PlayerBlack = match.PlayerBlack
	m.Turn = match.TurnCount
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

func (m *matchDetailsResponse) copyMatch(match gogo.Match) {
	m.ID = match.ID
	m.StartedAt = match.StartTime.Unix()
	m.GridSize = match.GridSize
	m.PlayerWhite = match.PlayerWhite
	m.PlayerBlack = match.PlayerBlack
	m.Turn = match.TurnCount
	m.GameBoard = match.GameBoard.Positions
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
	getMatches() (matches []gogo.Match, err error)
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
