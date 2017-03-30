package handlers

import (
	"neo/consumer"
	"zion/db"
	"zion/event"

	"fmt"

	"gopkg.in/mgo.v2/bson"
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
		case event.EVENT_ADD_TAGS:
			go addTags(tagEvent)
		case event.EVENT_TAG_FOLLOW:
			go followTag(tagEvent)
		case event.EVENT_TAG_UNFOLLOW:
			go unfollowTag(tagEvent)
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
	log.Infoln("remove tag", tagEvent.Tag.TagId)
}

func addTags(tagEvent *event.TagEvent) {
	var (
		err     error
		tag     *db.Tag
		feedTag *db.FeedTag
	)

	//Create session for every request
	session := mgoSession.Copy()
	defer session.Close()

	feedTag = new(db.FeedTag)
	tags := make([]db.FeedTag, 0, 1)
	for _, tagId := range tagEvent.TagIds {
		tag, err = db.GetTagById(session, bson.ObjectIdHex(tagId))
		if err != nil {
			log.Errorln(err)
			return
		}
		feedTag.TagId = tag.TagId
		feedTag.Name = tag.Name
		feedTag.AddUser = tagEvent.AddUser
		tags = append(tags, *feedTag)
	}
	fmt.Println(tags)
	err = db.AddTags(session, tagEvent.FeedId, tags)
	if err != nil {
		log.Errorln(err)
	}
	log.Infoln(tagEvent.AddUser, "add tags", tagEvent.TagIds)
}

func followTag(tagEvent *event.TagEvent) {
	var (
		err error
	)

	//Create session for every request
	session := mgoSession.Copy()
	defer session.Close()

	_, err = db.FollowTag(session, tagEvent.InitiatorId, tagEvent.TagId)
	if err != nil {
		log.Errorln(err)
	}
	log.Infoln(tagEvent.InitiatorId, "start to follow tag", tagEvent.TagId)
}

func unfollowTag(tagEvent *event.TagEvent) {
	var (
		err error
	)

	//Create session for every request
	session := mgoSession.Copy()
	defer session.Close()

	err = db.UnFollowTag(session, tagEvent.InitiatorId, tagEvent.TagId)
	if err != nil {
		log.Errorln(err)
	}
	log.Infoln(tagEvent.InitiatorId, "unfollow tag", tagEvent.TagId)
}
