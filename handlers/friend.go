package handlers

import (
	"neo/consumer"
	"zion/db"
	"zion/event"
)

func FriendHandler(msg *consumer.Message) {
	var (
		friendEvent *event.FriendEvent
		err         error
	)

	err = msg.ReadJSON(&friendEvent)
	if err == nil {
		switch friendEvent.EventId {
		case event.EVENT_FRIEND_FOLLOW:
			go followUser(friendEvent)
		case event.EVENT_FRIEND_UNFOLLOW:
			go unfollowUser(friendEvent)
		}
		msg.Success()
	} else {
		log.Errorln(err)
	}
}

func followUser(friendEvent *event.FriendEvent) {
	var (
		err error
	)

	//Create session for every request
	session := mgoSession.Copy()
	defer session.Close()

	_, err = db.FollowUser(session, friendEvent.InitiatorId, friendEvent.ResponderId)
	if err != nil {
		log.Errorln(err)
	}
	log.Infoln(friendEvent.InitiatorId, "start to follow", friendEvent.ResponderId)
}

func unfollowUser(friendEvent *event.FriendEvent) {
	var (
		err error
	)

	//Create session for every request
	session := mgoSession.Copy()
	defer session.Close()

	err = db.UnfollowUser(session, friendEvent.InitiatorId, friendEvent.ResponderId)
	if err != nil {
		log.Errorln(err)
	}
	log.Infoln(friendEvent.InitiatorId, "unfollow", friendEvent.ResponderId)
}
