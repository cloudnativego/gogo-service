package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/cloudnativego/cfmgo"
	"github.com/cloudnativego/cfmgo/params"
	"github.com/cloudnativego/gogo-engine"
	"gopkg.in/mgo.v2/bson"
)

type mongoMatchRepository struct {
	Collection cfmgo.Collection
}

type matchRecord struct {
	RecordID    bson.ObjectId `bson:"_id,omitempty" json:"id"`
	MatchID     string        `json:"match_id"`
	TurnCount   int           `json:"turn_count"`
	GridSize    int           `json:"grid_size"`
	StartTime   string        `json:"start_time"`
	GameBoard   [][]byte      `json:"game_board"`
	PlayerBlack string        `json:"player_black"`
	PlayerWhite string        `json:"player_white"`
}

//This is poorly named.  Name should communicate connecting to mongodb collection
//which may or may not be empty.
func NewMongoMatchRepository(col cfmgo.Collection) (repo *mongoMatchRepository) {
	repo = &mongoMatchRepository{
		Collection: col,
	}
	return
}

func (r *mongoMatchRepository) addMatch(match gogo.Match) (err error) {
	mr := convertMatchToMatchRecord(match)

	r.Collection.Wake()
	_, err = r.Collection.UpsertID(mr.RecordID, mr)
	return
}

func (r *mongoMatchRepository) getMatches() (matches []matchRecord, err error) {
	r.Collection.Wake()

	matches = make([]matchRecord, 0)

	_, err = r.Collection.Find(cfmgo.ParamsUnfiltered, &matches)
	return matches, err
}

func (r *mongoMatchRepository) getMatch(id string) (match gogo.Match, err error) {
	matchRecord, err := r.getMongoMatch(id)
	if err == nil {
		match = convertMatchRecordToMatch(matchRecord)
	}
	return
}

func (r *mongoMatchRepository) getMongoMatch(id string) (mongoMatch matchRecord, err error) {
	r.Collection.Wake()

	matches := make([]matchRecord, 0)
	query := bson.M{"match_id": id}
	params := &params.RequestParams{
		Q: query,
	}

	count, err := r.Collection.Find(params, &matches)
	if count == 0 {
		err = errors.New("Match not found")
	}
	if err == nil {
		mongoMatch = matches[0]
	}
	return
}

func (r *mongoMatchRepository) updateMatch(id string, match gogo.Match) (err error) {
	foundMatch, err := r.getMongoMatch(id)
	if err == nil {
		mr := convertMatchToMatchRecord(match)
		mr.RecordID = foundMatch.RecordID
		_, err = r.Collection.UpsertID(mr.RecordID, mr)
	}
	return
}

func convertMatchToMatchRecord(m gogo.Match) (mr *matchRecord) {
	mr = &matchRecord{
		RecordID:    bson.NewObjectId(),
		MatchID:     m.ID,
		TurnCount:   m.TurnCount,
		GridSize:    m.GridSize,
		StartTime:   m.StartTime.Format("2006-01-02 15:04:05"),
		PlayerBlack: m.PlayerBlack,
		PlayerWhite: m.PlayerWhite,
	}
	return
}

func convertMatchRecordToMatch(mr matchRecord) (m gogo.Match) {
	t, err := time.Parse("2006-01-02 15:04:05", mr.StartTime)
	if err != nil {
		fmt.Printf("Error parsing time value in Match Record: %v", err)
	} else {
		m = gogo.Match{
			ID:          mr.MatchID,
			TurnCount:   mr.TurnCount,
			GridSize:    mr.GridSize,
			StartTime:   t,
			PlayerBlack: mr.PlayerBlack,
			PlayerWhite: mr.PlayerWhite,
		}
	}
	return
}
