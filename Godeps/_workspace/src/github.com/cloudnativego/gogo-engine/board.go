package gogo

import "fmt"

// PerformMove accepts an intent to perform a move, validates it, and returns
// the corresponding change in match as a new match struct.
func (gameboard GameBoard) PerformMove(move Move) (outBoard GameBoard, err error) {
	//outBoard = NewBoard(cap(gameboard.Positions))
	outBoard = gameboard.copy()

	// TODO - this should eventually be part of another method that runs through
	// all the game rules in appropriate order.
	if gameboard.Positions[move.Position.X][move.Position.Y] != EmptyPosition {
		return outBoard, fmt.Errorf(RulesFailureSpaceOccupied, move)
	}

	outBoard.Positions[move.Position.X][move.Position.Y] = move.Player
	return outBoard, err
}
