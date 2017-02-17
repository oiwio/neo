package handlers

import (
	"neo/consumer"
	"zion/db"
	"zion/event"
)

func TagHandler(msg *consumer.Message) {
	var (
		tagEvent *event.TagEvent
		err      error
	)

	err = msg.ReadJSON(&tagEvent)
	if err == nil {
		switch tagEvent.EventId {
		case event.EVENT_TAG_CREATE:
			go createTag(tagEvent)
		case event.EVENT_TAG_REMOVE:
			go removeTag(tagEvent)
		}
		msg.Success()
	} else {
		log.Errorln(err)
	}
}

func createTag(tagEvent *event.TagEvent) {
	var (
		err error
	)

	//Create session for every request
	session := mgoSession.Copy()
	defer session.Close()

	_, err = db.NewTag(session, tagEvent.Tag)
	if err != nil {
		log.Errorln(err)
	}
	log.Infoln("create tag", tagEvent.Tag.TagId)
}

func removeTag(tagEvent *event.TagEvent) {

	log.Infoln("create tag", tagEvent.Tag.TagId)
}
