package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/cloudnativego/gogo-engine"
	"github.com/unrolled/render"
)

func createMatchHandler(formatter *render.Render, repo matchRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		payload, _ := ioutil.ReadAll(req.Body)
		var newMatchRequest newMatchRequest
		json.Unmarshal(payload, &newMatchRequest)

		newMatch := gogo.NewMatch(newMatchRequest.GridSize, newMatchRequest.Players[0].Name, newMatchRequest.Players[1].Name)
		repo.addMatch(newMatch)
		w.Header().Add("Location", "/matches/"+newMatch.ID)
		formatter.JSON(w, http.StatusCreated, &newMatchResponse{ID: newMatch.ID, GridSize: newMatch.GridSize, Players: newMatchRequest.Players})
	}
}
