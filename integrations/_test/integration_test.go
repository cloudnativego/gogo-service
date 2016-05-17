package integrations_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloudfoundry-community/go-cfenv"
	. "github.com/cloudnativego/gogo-service/service"
)

var (
	appEnv, _       = cfenv.Current()
	server          = NewServer(appEnv)
	firstMatchBody  = []byte("{\n  \"gridsize\": 19,\n  \"playerWhite\": \"L'Carpetron Dookmarriott\",\n  \"playerBlack\": \"Hingle McCringleberry\"\n}")
	secondMatchBody = []byte("{\n  \"gridsize\": 19,\n  \"playerWhite\": \"Devoin Shower-Handel\",\n  \"playerBlack\": \"J'Dinkalage Morgoone\"\n}")
)

func TestIntegration(t *testing.T) {
	// Get empty match list
	emptyMatches, err := getMatchList(t)
	if len(emptyMatches) > 0 {
		t.Errorf("Expected get match list to return an empty array; received %d", len(emptyMatches))
	}

	// Add first match
	matchResponse, err := addMatch(t, firstMatchBody)
	if matchResponse.PlayerBlack != "Hingle McCringleberry" {
		t.Errorf("Didn't get expected black stone player name from creation, got '%s'", matchResponse.PlayerBlack)
	}

	matches, err := getMatchList(t)
	if err != nil {
		t.Errorf("Error getting match list, %v", err)
	}
	if len(matches) != 1 {
		t.Errorf("Expected 1 active match, got %d", len(matches))
	}
	if matches[0].PlayerWhite != "L'Carpetron Dookmarriott" {
		t.Errorf("Player white name was wrong, got %s", matches[0].PlayerWhite)
	}

	// Add second match
	matchResponse, err = addMatch(t, secondMatchBody)
	if matchResponse.PlayerBlack != "J'Dinkalage Morgoone" {
		t.Errorf("Didn't get expected black stone player name from creation, got '%s'", matchResponse.PlayerBlack)
	}

	matches, err = getMatchList(t)
	if err != nil {
		t.Errorf("Error getting match list, %v", err)
	}
	if len(matches) != 2 {
		t.Errorf("Expected 2 active match, got %d", len(matches))
	}
	if matches[1].PlayerWhite != "Devoin Shower-Handel" {
		t.Errorf("Player white name was wrong, got %s", matches[1].PlayerWhite)
	}

	// Get match details (first match)
	firstMatch, err := getMatchDetails(t, matches[0].ID)
	if firstMatch.GridSize != 19 {
		t.Errorf("Expected match gridsize to be 19; received %d", firstMatch.GridSize)
	}

	secondMatch := matches[1]

	// Add Move
	addMoveToMatch(t, firstMatch.ID, []byte("{\n  \"player\": 2,\n  \"position\": {\n    \"x\": 3,\n    \"y\": 10\n  }\n}"))

	updatedFirstMatch, err := getMatchDetails(t, firstMatch.ID)
	if err != nil {
		t.Errorf("Error getting match details, %v", err)
	}
	if updatedFirstMatch.GameBoard[3][10] != 2 {
		t.Errorf("Expected gameboard position 3,10 to be 2, received: %d", updatedFirstMatch.GameBoard[3][10])
	}

	originalSecondMatch, _ := getMatchDetails(t, secondMatch.ID)
	if originalSecondMatch.GameBoard[3][10] != 0 {
		t.Errorf("Expected gameboard position 3,10 to be 0, received: %d", originalSecondMatch.GameBoard[3][10])
	}

	addMoveToMatch(t, secondMatch.ID, []byte("{\n  \"player\": 1,\n  \"position\": {\n    \"x\": 3,\n    \"y\": 10\n  }\n}"))

	updatedFirstMatch, _ = getMatchDetails(t, firstMatch.ID)
	if updatedFirstMatch.GameBoard[3][10] != 2 {
		t.Errorf("Expected gameboard position 3,10 to be 2, received: %d", updatedFirstMatch.GameBoard[3][10])
	}

	updatedSecondMatch, _ := getMatchDetails(t, secondMatch.ID)
	if updatedSecondMatch.GameBoard[3][10] != 1 {
		t.Errorf("Expected gameboard position 3,10 to be 1, received: %d", updatedSecondMatch.GameBoard[3][10])
	}
}

// ----------------- Utility Functions ------------

func getMatchList(t *testing.T) (matches []newMatchResponse, err error) {
	getMatchListRequest, _ := http.NewRequest("GET", "/matches", nil)
	recorder := httptest.NewRecorder()
	server.ServeHTTP(recorder, getMatchListRequest)
	matches = make([]newMatchResponse, 0)
	err = json.Unmarshal(recorder.Body.Bytes(), &matches)
	if err != nil {
		t.Errorf("Error unmarshaling match list, %v", err)
	} else {
		if recorder.Code != 200 {
			t.Errorf("Expected match list code to be 200, got %d", recorder.Code)
		} else {
			fmt.Println("\tQueried Match List OK")
		}
	}
	return
}

func addMatch(t *testing.T, body []byte) (reply newMatchResponse, err error) {
	recorder := httptest.NewRecorder()
	createMatchRequest, _ := http.NewRequest("POST", "/matches", bytes.NewBuffer(body))
	server.ServeHTTP(recorder, createMatchRequest)
	if recorder.Code != 201 {
		t.Errorf("Error creating new match, expected 201 code, got %d", recorder.Code)
	}
	var matchResponse newMatchResponse
	err = json.Unmarshal(recorder.Body.Bytes(), &matchResponse)
	if err != nil {
		t.Errorf("Error unmarshaling new match response: %v", err)
	} else {
		fmt.Println("\tAdded Match OK")
	}
	reply = matchResponse
	return
}

func getMatchDetails(t *testing.T, ID string) (match matchDetailsResponse, err error) {
	recorder := httptest.NewRecorder()
	matchURL := fmt.Sprintf("/matches/%s", ID)
	getMatchDetailsRequest, _ := http.NewRequest("GET", matchURL, nil)
	server.ServeHTTP(recorder, getMatchDetailsRequest)
	if recorder.Code != 200 {
		t.Errorf("Error getting match details: %d", recorder.Code)
	}
	err = json.Unmarshal(recorder.Body.Bytes(), &match)
	if err != nil {
		t.Errorf("Error unmarshaling match details: %v", err)
	} else {
		fmt.Println("\tQueried Match Details OK")
	}
	return
}

func addMoveToMatch(t *testing.T, ID string, body []byte) {
	recorder := httptest.NewRecorder()
	requestString := fmt.Sprintf("/matches/%s/moves", ID)
	matchMove := bytes.NewBuffer(body)
	addMoveRequest, _ := http.NewRequest("POST", requestString, matchMove)
	server.ServeHTTP(recorder, addMoveRequest)
	if recorder.Code != 201 {
		t.Errorf("Error adding move to match: %d", recorder.Code)
	} else {
		fmt.Println("\tAdded Move to Match OK")
	}
	return
}

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
