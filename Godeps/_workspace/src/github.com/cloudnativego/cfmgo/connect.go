// A MongoDB integration package for Cloud Foundry.
package cfmgo

import (
	"fmt"

	"gopkg.in/mgo.v2"
)

//Connect to the specified database and return a Collection object for the specified collection.
func Connect(dialer CollectionDialer, URI string, collectionName string) (collection Collection) {
	var (
		err      error
		dialInfo *mgo.DialInfo
	)

	if dialInfo, err = mgo.ParseURL(URI); err != nil || dialInfo.Database == "" {
		panic(fmt.Sprintf("cannot parse given URI %s due to error: %s", URI, err.Error()))
	}

	if collection, err = dialer(URI, dialInfo.Database, collectionName); err != nil {
		panic(fmt.Sprintf("cannot dial connection due to error: %s URI:%s col:%s db:%s", err.Error(), URI, collectionName, dialInfo.Database))
	}
	return
}
