package gogo

func genMove(x int, y int, player byte) Move {
	return Move{Player: player, Position: Coordinate{X: x, Y: y}}
}

// Copy copies the state of an existing board into a new board.
func (gameboard GameBoard) copy() GameBoard {
	outBoard := newBoard(cap(gameboard.Positions))
	copy(outBoard.Positions, gameboard.Positions)
	return outBoard
}

// NewBoard creates a new gameboard of a given size. Gameboards must always be square.
func newBoard(size int) GameBoard {
	outBoard := GameBoard{}
	a := make([][]byte, size)
	for i := range a {
		a[i] = make([]byte, size)
	}
	outBoard.Positions = a
	return outBoard
}
