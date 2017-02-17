package handlers

import (
	"neo/config"
	"os"

	"github.com/Sirupsen/logrus"
	mgo "gopkg.in/mgo.v2"
)

var (
	log           *logrus.Logger
	mgoSession    *mgo.Session
	configuration config.Config
)

func init() {
	var (
		err error
	)
	log = logrus.New()

	configuration = config.New()

	log = logrus.New()
	log.Formatter = new(logrus.JSONFormatter)
	log.Formatter = new(logrus.TextFormatter) // default

	file, err := os.OpenFile(configuration.Log.LogPath, os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		log.Out = file
	} else {
		log.Info("Failed to log to file, using default stderr")
	}
	log.Level = logrus.DebugLevel

	mgoSession, err = mgo.Dial(configuration.MongoDB.Host)
	if err != nil {
		panic(err)
	}
	mgoSession.SetMode(mgo.Monotonic, true)

}
