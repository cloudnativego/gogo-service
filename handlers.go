package main

import (
	"net/http"

	"code.google.com/p/go-uuid/uuid"

	"github.com/cloudnativego/gogo-engine"
	"github.com/unrolled/render"
)

func createMatchHandler(formatter *render.Render, repo matchRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		newMatch := gogo.NewMatch(5)
		repo.addMatch(newMatch)
		guid := uuid.New()
		w.Header().Add("Location", "/matches/"+guid)
		formatter.JSON(w, http.StatusCreated, &newMatchResponse{Id: guid})
	}
}
