package handlers

import (
	"fmt"
	"neo/consumer"
	"zion/db"
	"zion/event"

	"github.com/ylywyn/jpush-api-go-client"
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
		err error
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

		var pf jpushclient.Platform
		pf.Add(jpushclient.ANDROID)
		pf.Add(jpushclient.IOS)

		var ad jpushclient.Audience
		var notice jpushclient.Notice
		var option jpushclient.Option
		option.ApnsProduction = true
		if feedEvent.Comment.Reference != nil {
			s := []string{feedEvent.Comment.Reference.UserId.Hex()}
			ad.SetAlias(s)
			notice.SetAlert(fmt.Sprintf("%s 回复你的评论：%s", feedEvent.Comment.Author.NickName, feedEvent.Comment.Content))
		} else {
			// feed, err := db.GetFeedById(session, feedEvent.Comment.FeedId)
			// if err != nil {
			// 	log.Errorln("评论失败")
			// }
			// s := []string{feed.UserId.Hex()}
			// ad.SetAlias(s)
			ad.All()
			notice.SetAlert(fmt.Sprintf("%s 回复了你的发布：%s", feedEvent.Comment.Author.NickName, feedEvent.Comment.Content))
		}
		payload := jpushclient.NewPushPayLoad()
		payload.SetPlatform(&pf)
		payload.SetAudience(&ad)
		payload.SetNotice(&notice)
		payload.SetOptions(&option)

		log.Infoln(payload)

		bytes, _ := payload.ToBytes()

		c := jpushclient.NewPushClient(configuration.JPush.Secret, configuration.JPush.AppKey)
		log.Infoln(c)
		str, err := c.Send(bytes)
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
