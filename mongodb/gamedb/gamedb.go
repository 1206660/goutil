package gamedb

import (
	"github.com/name5566/leaf/db/mongodb"
	"github.com/name5566/leaf/log"
)

var mongoDB *mongodb.DialContext

func Connect(dbUrl string, dbMaxConnNum int) *mongodb.DialContext {
	db, err := mongodb.Dial(dbUrl, dbMaxConnNum)
	if err != nil {
		log.Fatal("dial %v mongodb error: %v", dbUrl, err)
	}
	mongoDB = db
	return mongoDB
}

func MongoDB() *mongodb.DialContext {
	return mongoDB
}
