package handlers

import (
	"fmt"
	"neo/consumer"
	"zion/db"
	"zion/event"
	"zion/push"

	"gopkg.in/mgo.v2/bson"
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
	log.Infoln("create feed", feedEvent.Feed.FeedId)
}

func pushFeedtoFollowedUsers(userId bson.ObjectId) error {
	var (
		err error
	)
	return err
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
	log.Infoln("remove feed", feedEvent.FeedId)
}

func createComment(feedEvent *event.FeedEvent) {
	var (
		err     error
		users   []string
		content string
		str     string
	)

	//Create session for every request
	session := mgoSession.Copy()
	defer session.Close()

	_, err = db.NewComment(session, feedEvent.Comment)
	if err != nil {
		log.Errorln(err)
	}
	log.Infoln("create comment", feedEvent.Comment.CommentId)

	go func() {
		//Create session for every request
		session := mgoSession.Copy()
		defer session.Close()

		if feedEvent.Comment.Reference != nil {
			users = []string{feedEvent.Comment.Reference.UserId.Hex()}
			content = fmt.Sprintf("%s 回复你的评论：%s", feedEvent.Comment.Author.NickName, feedEvent.Comment.Content)
		} else {
			feed, err := db.GetFeedById(session, feedEvent.Comment.FeedId)
			if err != nil {
				log.Errorln("评论失败")
			}
			users = []string{feed.UserId.Hex()}
			content = fmt.Sprintf("%s 回复了你的发布：%s", feedEvent.Comment.Author.NickName, feedEvent.Comment.Content)
		}
		str, err = push.JpushWithUserIds(users, content)
		if err != nil {
			log.Errorln(err)
		} else {
			log.Infoln("推送成功：", str)
		}
	}()
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
	log.Infoln("remove comment", feedEvent.CommentId)
}
