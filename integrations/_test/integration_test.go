package integrations_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	. "github.com/cloudnativego/gogo-service"
)

/*
	..GetMatches empty array
	..AddMatch
	..AddSecondMatch
	..GetMatches with data
	..GetMatchDetails for each
	AddMove to first
	GetMatchDetails for both
	AddMove to second
	GetMatchDetails for both
*/

func TestIntegration(t *testing.T) {
	server := NewServer()

	getMatchListRequest, _ := http.NewRequest("GET", "/matches", nil)
	//	addMoveRequest, _ := http.NewRequest("POST", "/matches/{id}/moves", nil)

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
