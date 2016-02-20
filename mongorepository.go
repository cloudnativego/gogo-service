package main

import (
	"github.com/cloudnativego/cfmgo"
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
func newMongoMatchRepository(col cfmgo.Collection) (repo *mongoMatchRepository) {
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
	return
}

func (r *mongoMatchRepository) updateMatch(id string, match gogo.Match) error {
	return nil
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
