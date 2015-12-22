package gogo

import (
	"fmt"
	"testing"
)

func TestEmptyBoardAcceptsMoves(t *testing.T) {
	myMove := genMove(1, 1, PlayerWhite)

	emptyBoard := newBoard(19)

	newBoard, err := emptyBoard.performMove(myMove)

	if err != nil {
		t.Errorf("Should not have received an error when placing simple positions.")
	}

	if newBoard.Positions[1][1] != PlayerWhite {
		t.Errorf("Putting a white stone at 1,1 should result in a board with that stone in place.")
	}

	for r := range newBoard.Positions {
		for c := range newBoard.Positions[r] {
			if r != 1 && c != 1 {
				if newBoard.Positions[r][c] != EmptyPosition {
					t.Errorf("Encountered non-empty board position when it should have been empty.")
				}
			}
		}
	}
}

// Ensure board retains state between moves.
func TestBoardRemembersMultipleMoves(t *testing.T) {
	moveA := genMove(1, 1, PlayerBlack)
	moveB := genMove(2, 2, PlayerWhite)
	moveC := genMove(3, 3, PlayerBlack)

	emptyBoard := newBoard(19)
	newBoard, _ := emptyBoard.performMove(moveA)
	newBoard, _ = newBoard.performMove(moveB)
	newBoard, _ = newBoard.performMove(moveC)

	if newBoard.Positions[1][1] != PlayerBlack {
		t.Errorf("Expected black stone at 1,1 and it wasn't there.")
	}
	if newBoard.Positions[2][2] != PlayerWhite {
		t.Errorf("Expected white stone at 1,1 and it wasn't there.")
	}
	if newBoard.Positions[3][3] != PlayerBlack {
		t.Errorf("Expected black stone at 3,3 and it wasn't there.")
	}
}

func TestCannotPlayOnOccupiedPosition(t *testing.T) {
	myMove := genMove(5, 5, PlayerWhite)
	emptyBoard := newBoard(19)

	emptyBoard.Positions[5][5] = PlayerBlack

	_, err := emptyBoard.performMove(myMove)

	expectedError := fmt.Sprintf(RulesFailureSpaceOccupied, myMove)

	if err == nil {
		t.Errorf("GoGo should have prevented a move onto an occupied position, but didn't return an error.")
	} else {
		if err.Error() != expectedError {
			t.Errorf("Expected error '%s' but instead got '%s' when playing on an occupied position.", expectedError, err.Error())
		}
	}
}
