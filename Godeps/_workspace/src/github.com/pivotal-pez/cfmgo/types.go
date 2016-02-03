package cfmgo

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type (
	//Collection - an interface representing a trimmed down collection object
	Collection interface {
		Wake()
		Close()
		Find(params Params, result interface{}) (count int, err error)
		FindOne(id string, result interface{}) (err error)
		UpsertID(selector interface{}, update interface{}) (info *mgo.ChangeInfo, err error)
		FindAndModify(selector interface{}, update interface{}, target interface{}) (info *mgo.ChangeInfo, err error)
		Count() (int, error)
	}

	//CollectionRepo - mgo collection adapter
	CollectionRepo struct {
		Col     *mgo.Collection
		session *mgo.Session
	}

	//CollectionDialer - a funciton type to dial for collections
	CollectionDialer func(url string, dbname string, collectionname string) (collection Collection, err error)

	//Params interface exposes mongodb-specific query parameters: Selector, Scope, Limit, and Offset
	Params interface {
		Selector() bson.M
		Scope() bson.M
		Limit() int
		Offset() int
	}
)
