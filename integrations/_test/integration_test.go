package integrations_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	. "github.com/cloudnativego/gogo-service"
)

func TestIntegration(t *testing.T) {
	server := NewServer()

	getMatchListRequest, _ := http.NewRequest("GET", "/matches", nil)

	firstMatchBody := []byte("{\n  \"gridsize\": 19,\n  \"playerWhite\": \"L'Carpetron Dookmarriott\",\n  \"playerBlack\": \"Hingle McCringleberry\"\n}")
	secondMatchBody := []byte("{\n  \"gridsize\": 19,\n  \"playerWhite\": \"Devoin Shower-Handel\",\n  \"playerBlack\": \"J'Dinkalage Morgoone\"\n}")

	//GetMatches with empty repo
	recorder := httptest.NewRecorder()
	server.ServeHTTP(recorder, getMatchListRequest)
	if recorder.Code != 200 {
		t.Errorf("Error getting match list: %d", recorder.Code)
	}
	if strings.TrimSpace(recorder.Body.String()) != "[]" {
		t.Errorf("Expected get match list to return an empty array; received %s", recorder.Body.String())
	}

	recorder = httptest.NewRecorder()
	createMatchRequest, _ := http.NewRequest("POST", "/matches", bytes.NewBuffer(firstMatchBody))
	server.ServeHTTP(recorder, createMatchRequest)
	if recorder.Code != 201 {
		t.Errorf("Error creating new match, expected 201 code, got %d", recorder.Code)
	}
	var matchResponse matchDetailsResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &matchResponse)
	if matchResponse.PlayerBlack != "Hingle McCringleberry" {
		t.Errorf("Didn't get expected black stone player name from creation, got '%s'", matchResponse.PlayerBlack)
	}

	recorder = httptest.NewRecorder()
	server.ServeHTTP(recorder, getMatchListRequest)
	matches := make([]newMatchResponse, 0)
	err = json.Unmarshal(recorder.Body.Bytes(), &matches)
	if err != nil {
		t.Errorf("Error unmarshaling match list, %v", err)
	}
	if len(matches) != 1 {
		t.Errorf("Expected 1 active match, got %d", len(matches))
	}
	if matches[0].PlayerWhite != "L'Carpetron Dookmarriott" {
		t.Errorf("Player white name was wrong, got %s", matches[0].PlayerWhite)
	}

	recorder = httptest.NewRecorder()
	createMatchRequest, _ = http.NewRequest("POST", "/matches", bytes.NewBuffer(secondMatchBody))
	server.ServeHTTP(recorder, createMatchRequest)
	if recorder.Code != 201 {
		t.Errorf("Error creating new match, expected 201 code, got %d", recorder.Code)
	}

	err = json.Unmarshal(recorder.Body.Bytes(), &matchResponse)
	if matchResponse.PlayerBlack != "J'Dinkalage Morgoone" {
		t.Errorf("Didn't get expected black stone player name from creation, got '%s'", matchResponse.PlayerBlack)
	}

	recorder = httptest.NewRecorder()
	server.ServeHTTP(recorder, getMatchListRequest)
	matches = make([]newMatchResponse, 0)
	err = json.Unmarshal(recorder.Body.Bytes(), &matches)
	if err != nil {
		t.Errorf("Error unmarshaling match list, %v", err)
	}
	if len(matches) != 2 {
		t.Errorf("Expected 2 active match, got %d", len(matches))
	}
	if matches[1].PlayerWhite != "Devoin Shower-Handel" {
		t.Errorf("Player white name was wrong, got %s", matches[1].PlayerWhite)
	}

	recorder = httptest.NewRecorder()
	getMatchDetailsRequest, _ := http.NewRequest("GET", "/matches/"+matches[0].ID, nil)
	server.ServeHTTP(recorder, getMatchDetailsRequest)
	if recorder.Code != 200 {
		t.Errorf("Error getting match details: %d", recorder.Code)
	}
	var firstMatch matchDetailsResponse
	err = json.Unmarshal(recorder.Body.Bytes(), &firstMatch)
	if err != nil {
		t.Errorf("Error unmarshaling match details: %v", err)
	}
	if firstMatch.GridSize != 19 {
		t.Errorf("Expected match gridsize to be 19; received %d", firstMatch.GridSize)
	}

	secondMatch := matches[1]

	recorder = httptest.NewRecorder()
	requestString := fmt.Sprintf("/matches/%s/moves", firstMatch.ID)
	matchMove := bytes.NewBuffer([]byte("{\n  \"player\": 2,\n  \"position\": {\n    \"x\": 3,\n    \"y\": 10\n  }\n}"))
	addMoveRequest, _ := http.NewRequest("POST", requestString, matchMove)
	server.ServeHTTP(recorder, addMoveRequest)
	if recorder.Code != 201 {
		t.Errorf("Error adding move to match: %d", recorder.Code)
	}

	recorder = httptest.NewRecorder()
	requestString = fmt.Sprintf("/matches/%s", firstMatch.ID)
	getMatchDetailsRequest, _ = http.NewRequest("GET", requestString, nil)
	server.ServeHTTP(recorder, getMatchDetailsRequest)
	if recorder.Code != 200 {
		t.Errorf("Error getting match details: %d", recorder.Code)
	}

	var updatedFirstMatch matchDetailsResponse
	err = json.Unmarshal(recorder.Body.Bytes(), &updatedFirstMatch)
	if err != nil {
		t.Errorf("Error unmarshaling match details: %v", err)
	}
	if updatedFirstMatch.GameBoard[3][10] != 2 {
		t.Errorf("Expected gameboard position 3,10 to be 2, received: %d", updatedFirstMatch.GameBoard[3][10])
	}

	recorder = httptest.NewRecorder()
	requestString = fmt.Sprintf("/matches/%s", secondMatch.ID)
	getMatchDetailsRequest, _ = http.NewRequest("GET", requestString, nil)
	server.ServeHTTP(recorder, getMatchDetailsRequest)
	if recorder.Code != 200 {
		t.Errorf("Error getting match details: %d", recorder.Code)
	}

	var originalSecondMatch matchDetailsResponse
	err = json.Unmarshal(recorder.Body.Bytes(), &originalSecondMatch)
	if err != nil {
		t.Errorf("Error unmarshaling match details: %v", err)
	}
	if originalSecondMatch.GameBoard[3][10] != 0 {
		t.Errorf("Expected gameboard position 3,10 to be 0, received: %d", originalSecondMatch.GameBoard[3][10])
	}

	recorder = httptest.NewRecorder()
	requestString = fmt.Sprintf("/matches/%s/moves", secondMatch.ID)
	matchMove = bytes.NewBuffer([]byte("{\n  \"player\": 1,\n  \"position\": {\n    \"x\": 3,\n    \"y\": 10\n  }\n}"))
	addMoveRequest, _ = http.NewRequest("POST", requestString, matchMove)
	server.ServeHTTP(recorder, addMoveRequest)
	if recorder.Code != 201 {
		t.Errorf("Error adding move to match: %d", recorder.Code)
	}

	recorder = httptest.NewRecorder()
	requestString = fmt.Sprintf("/matches/%s", firstMatch.ID)
	getMatchDetailsRequest, _ = http.NewRequest("GET", requestString, nil)
	server.ServeHTTP(recorder, getMatchDetailsRequest)
	if recorder.Code != 200 {
		t.Errorf("Error getting match details: %d", recorder.Code)
	}

	err = json.Unmarshal(recorder.Body.Bytes(), &updatedFirstMatch)
	if err != nil {
		t.Errorf("Error unmarshaling match details: %v", err)
	}
	if updatedFirstMatch.GameBoard[3][10] != 2 {
		t.Errorf("Expected gameboard position 3,10 to be 2, received: %d", updatedFirstMatch.GameBoard[3][10])
	}

	recorder = httptest.NewRecorder()
	requestString = fmt.Sprintf("/matches/%s", secondMatch.ID)
	getMatchDetailsRequest, _ = http.NewRequest("GET", requestString, nil)
	server.ServeHTTP(recorder, getMatchDetailsRequest)
	if recorder.Code != 200 {
		t.Errorf("Error getting match details: %d", recorder.Code)
	}

	var updatedSecondMatch matchDetailsResponse
	err = json.Unmarshal(recorder.Body.Bytes(), &updatedSecondMatch)
	if err != nil {
		t.Errorf("Error unmarshaling match details: %v", err)
	}
	if updatedSecondMatch.GameBoard[3][10] != 1 {
		t.Errorf("Expected gameboard position 3,10 to be 1, received: %d", updatedSecondMatch.GameBoard[3][10])
	}
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
