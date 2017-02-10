package handlers

import (
	"neo/config"

	mgo "gopkg.in/mgo.v2"
)

var (
	mgoSession    *mgo.Session
	configuration config.Config
)

func init() {
	var (
		err error
	)

	configuration = config.New()

	// consumer.Connect(configuration.NSQ.Host)

	mgoSession, err = mgo.Dial(configuration.MongoDB.Host)
	if err != nil {
		panic(err)
	}
	mgoSession.SetMode(mgo.Monotonic, true)

}
