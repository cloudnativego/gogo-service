package cfmgo

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//NewCollectionDialer -- dials a new mongo connection
func NewCollectionDialer(url string, dbname string, collectionname string) (collection Collection, err error) {
	var session *mgo.Session

	if session, err = mgo.Dial(url); err == nil {
		session.SetMode(mgo.Monotonic, true)
		db := session.DB(dbname)
		collection = &CollectionRepo{
			Col:     db.C(collectionname),
			session: session,
		}
	}
	return
}

//Find -- finds all records matching given selector
func (s *CollectionRepo) Find(params Params, result interface{}) (count int, err error) {
	count, err = s.Col.Find(params.Selector()).Count()
	if err != nil {
		return
	}
	err = s.Col.Find(params.Selector()).
		Select(params.Scope()).
		Limit(params.Limit()).
		Skip(params.Offset()).
		All(result)
	return
}

//FindOne -- finds record with given ID
func (s *CollectionRepo) FindOne(id string, result interface{}) (err error) {

	if bson.IsObjectIdHex(id) {
		hex := bson.ObjectIdHex(id)
		err = s.Col.FindId(hex).One(result)

	} else {
		err = ErrInvalidID
	}
	return
}

//FindAndModify -- execute a normal upsert
func (s *CollectionRepo) FindAndModify(selector interface{}, update interface{}, result interface{}) (info *mgo.ChangeInfo, err error) {
	change := mgo.Change{
		Update:    update.(bson.M),
		ReturnNew: false,
	}
	info, err = s.Col.Find(selector.(bson.M)).Apply(change, result)
	return
}

//UpsertID -- upserts the given object to the given id
func (s *CollectionRepo) UpsertID(id interface{}, update interface{}) (info *mgo.ChangeInfo, err error) {
	info, err = s.Col.UpsertId(id, update)
	return
}

//Close -- closes the connection
func (s *CollectionRepo) Close() {
	if s.session != nil {
		s.session.Close()
	}
}

//Count -- counts the collection records
func (s *CollectionRepo) Count() (int, error) {
	return s.Col.Count()
}

//Wake - will ping and reconnect if need be
func (s *CollectionRepo) Wake() {
	if s.session.Ping() != nil {
		s.session = s.session.Clone()
	}
}
