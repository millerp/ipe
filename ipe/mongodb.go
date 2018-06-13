package ipe

import (
	"github.com/globalsign/mgo"
)

var mongoDBSession *mgo.Session
var mongoDBDefaultDatabase *mgo.Database
var mongoDBAppCollection *mgo.Collection

func StartMongoDBInstance(host string, dbname string) {
	var err error
	mongoDBSession, err = mgo.Dial(host)
	if err != nil {
		panic(err)
	}
	//defer mongoDBSession.Close()
	setDefaultDb(dbname)
}

func setDefaultDb(dbName string) (*mgo.Database) {
	mongoDBDefaultDatabase = mongoDBSession.DB(dbName)
	return mongoDBDefaultDatabase
}
