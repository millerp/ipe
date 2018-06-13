// Copyright 2016 Claudemiro Alves Feitosa Neto. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ipe

import (
	"sync"
	"github.com/globalsign/mgo/bson"
	"errors"
	"log"
)

// db represents a app database
// For now it there is only one memory database implementation
// but in the future I can write a sql implementation
type db interface {
	GetAppByAppID(appID string) (*app, error)
	GetAppByKey(key string) (*app, error)
	AddApp(*app) error
}

type memdb struct {
	sync.Mutex
	Apps []*app
}

func newMemdb() db {
	mongoDBAppCollection = mongoDBDefaultDatabase.C("apps")
	return &memdb{}
}

func (db *memdb) AddApp(a *app) error {
	db.Lock()
	count, err := mongoDBAppCollection.Find(bson.M{"appid": a.AppID}).Count()
	if err != nil {
		log.Println(err)
	}

	if count <= 0 {
		mongoDBAppCollection.Insert(a)
	}

	db.Apps = append(db.Apps, a)
	db.Unlock()
	return nil
}

// GetAppByAppID returns an App with by appID
func (db *memdb) GetAppByAppID(appID string) (*app, error) {
	for _, a := range db.Apps {
		if a.AppID == appID {
			return a, nil
		}
	}
	return nil, errors.New("App not found")

	result := new(app)
	err := mongoDBAppCollection.Find(bson.M{"appid": appID}).One(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetAppByKey returns an App with by key
func (db *memdb) GetAppByKey(key string) (*app, error) {
	for _, a := range db.Apps {
		if a.Key == key {
			return a, nil
		}
	}

	return nil, errors.New("App not found")

	result := new(app)
	err := mongoDBAppCollection.Find(bson.M{"key": key}).One(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
