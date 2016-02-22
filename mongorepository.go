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

//MongoMatchRepository object contains the matches collection and methods to
//act upon that data.
type MongoMatchRepository struct {
	Collection cfmgo.Collection
}

//MatchRecord is a struct representing a match in mongo.
type MatchRecord struct {
	RecordID    bson.ObjectId `bson:"_id,omitempty" json:"id"`
	MatchID     string        `bson:"match_id",json:"match_id"`
	TurnCount   int           `bson:"turn_count",json:"turn_count"`
	GridSize    int           `bson:"grid_size",json:"grid_size"`
	StartTime   string        `bson:"start_time",json:"start_time"`
	GameBoard   [][]byte      `bson:"game_board",json:"game_board"`
	PlayerBlack string        `bson:"player_black",json:"player_black"`
	PlayerWhite string        `bson:"player_white",json:"player_white"`
}

//NewMongoMatchRepository instantiates a new match repository object.
func NewMongoMatchRepository(col cfmgo.Collection) (repo *MongoMatchRepository) {
	repo = &MongoMatchRepository{
		Collection: col,
	}
	return
}

//AddMatch inserts a new record into the repo.
func (r *MongoMatchRepository) AddMatch(match gogo.Match) (err error) {
	r.Collection.Wake()
	mr := convertMatchToMatchRecord(match)
	_, err = r.Collection.UpsertID(mr.RecordID, mr)
	return
}

//GetMatch returns a match record from the repo based on the ID.
func (r *MongoMatchRepository) GetMatch(id string) (match gogo.Match, err error) {
	r.Collection.Wake()
	MatchRecord, err := r.getMongoMatch(id)
	if err == nil {
		match = convertMatchRecordToMatch(MatchRecord)
	}
	return
}

//GetMatches returns all matches in the repo.
//FIXME: Return value of []MatchRecord seems incorrect; should be []gogo.Match
func (r *MongoMatchRepository) GetMatches() (matches []MatchRecord, err error) {
	r.Collection.Wake()
	matches = make([]MatchRecord, 0)
	_, err = r.Collection.Find(cfmgo.ParamsUnfiltered, &matches)
	return
}

//UpdateMatch replaces the match state for the given ID with current match state.
func (r *MongoMatchRepository) UpdateMatch(id string, match gogo.Match) (err error) {
	r.Collection.Wake()
	foundMatch, err := r.getMongoMatch(id)
	if err == nil {
		mr := convertMatchToMatchRecord(match)
		mr.RecordID = foundMatch.RecordID
		_, err = r.Collection.UpsertID(mr.RecordID, mr)
	}
	return
}

func (r *MongoMatchRepository) getMongoMatch(id string) (mongoMatch MatchRecord, err error) {
	var matches []MatchRecord
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

func convertMatchToMatchRecord(m gogo.Match) (mr *MatchRecord) {
	mr = &MatchRecord{
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

func convertMatchRecordToMatch(mr MatchRecord) (m gogo.Match) {
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
