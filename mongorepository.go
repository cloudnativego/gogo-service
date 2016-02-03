package main

import (
	"net/url"

	"github.com/cloudnativego/gogo-engine"
	"github.com/pivotal-pez/cfmgo"
	"github.com/pivotal-pez/cfmgo/params"
	"gopkg.in/mgo.v2/bson"
)

type mongoMatchRepository struct {
	Collection cfmgo.Collection
}

type matchRecord struct {
	ID          bson.ObjectId `bson:"_id,omitempty" json:"id"`
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
	_, err = r.Collection.UpsertID(mr.ID, mr)
	return
}

func (r *mongoMatchRepository) getMatches() (matches []matchRecord, err error) {
	r.Collection.Wake()

	params := params.Extract(url.Values{}) //this is clunky with repo approach.
	matches = make([]matchRecord, 0)

	_, err = r.Collection.Find(params, &matches)
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
		//uuid is not a valid bson.ObjectIdHex object
		//consider bson in favor of uuid?
		ID:          bson.NewObjectId(),
		TurnCount:   m.TurnCount,
		GridSize:    m.GridSize,
		StartTime:   m.StartTime.Format("2006-01-02 15:04:05"),
		PlayerBlack: m.PlayerBlack,
		PlayerWhite: m.PlayerWhite,
	}
	return
}
