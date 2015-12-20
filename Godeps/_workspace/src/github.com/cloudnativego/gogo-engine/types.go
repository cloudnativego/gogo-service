package gogo

import "time"

const (
	// EmptyPosition is an empty spot on the board.
	EmptyPosition = 0

	// PlayerBlack indicates a position occupied by a black stone.
	PlayerBlack = 1

	// PlayerWhite indicates a position occupied by a white stone.
	PlayerWhite = 2

	// RulesFailureSpaceOccupied is the string displayed when a player attempts to play on an occupied spot.
	RulesFailureSpaceOccupied = "Cannot perform move (%s), the target position is already occupied."
)

// GameBoard is a two-dimensional array of stones.
type GameBoard struct {
	Positions [][]byte
}

// Match represents the state of an in-progress game of Go.
type Match struct {
	TurnCount   int
	GridSize    int
	ID          string
	StartTime   time.Time
	GameBoard   GameBoard
	PlayerBlack string
	PlayerWhite string
}

// Coordinate represents an x-y coordinate of an intersection on a Go board
type Coordinate struct {
	X int
	Y int
}

// Move represents an intent to perform a move by a player. Empty Position indicates a pass.
type Move struct {
	Player   byte
	Position Coordinate
}
