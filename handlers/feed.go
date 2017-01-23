package handlers

import (
	"fmt"
	"neo/consumer"
	"zion/db"
	"zion/event"

	log "github.com/Sirupsen/logrus"
)

func FeedHandler(msg *consumer.Message) {
	var (
		feedEvent *event.FeedEvent
		err       error
	)

	err = msg.ReadJSON(&feedEvent)
	if err == nil {
		switch feedEvent.EventId {
		case event.EVENT_FEED_CREATE:
			go createFeed(feedEvent)
		case event.EVENT_FEED_UPDATE:
			go fmt.Println("feed update")
		case event.EVENT_FEED_REMOVE:
			go removeFeed(feedEvent)
		case event.EVENT_FEED_COMMENT_POST:
			go createComment(feedEvent)
		case event.EVENT_FEED_COMMENT_REMOVE:
			go removeComment(feedEvent)
		}
		msg.Success()
	} else {
		log.Errorln(err)
	}
}

func createFeed(feedEvent *event.FeedEvent) {
	var (
		err error
	)

	//Create session for every request
	session := mgoSession.Copy()
	defer session.Close()

	_, err = db.NewFeed(session, feedEvent.Feed)
	if err != nil {
		log.Errorln(err)
	}
	log.Infoln("create feed ", feedEvent.Feed.FeedId)
}

func removeFeed(feedEvent *event.FeedEvent) {
	var (
		err error
	)

	//Create session for every request
	session := mgoSession.Copy()
	defer session.Close()

	err = db.DeleteFeed(session, feedEvent.FeedId)
	err = db.DeleteCommentByFeedId(session, feedEvent.FeedId)
	if err != nil {
		log.Errorln(err)
	}
	log.Infoln("remove feed ", feedEvent.FeedId)
}

func createComment(feedEvent *event.FeedEvent) {
	var (
		err error
	)

	//Create session for every request
	session := mgoSession.Copy()
	defer session.Close()

	_, err = db.NewComment(session, feedEvent.Comment)
	if err != nil {
		log.Errorln(err)
	}
	log.Infoln("create comment ", feedEvent.Comment.CommentId)
}

func removeComment(feedEvent *event.FeedEvent) {
	var (
		err error
	)

	//Create session for every request
	session := mgoSession.Copy()
	defer session.Close()

	err = db.DeleteComment(session, feedEvent.CommentId)
	if err != nil {
		log.Errorln(err)
	}
	log.Infoln("remove comment ", feedEvent.CommentId)
}
