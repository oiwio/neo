package handlers

import (
	"fmt"
	"neo/consumer"
	"zion/db"
	"zion/event"
)

func UserHandler(msg *consumer.Message) {
	var (
		userEvent *event.UserEvent
		err       error
	)

	err = msg.ReadJSON(&userEvent)
	if err == nil {
		switch userEvent.EventId {
		case event.EVENT_USER_UPDATE_PROFILE:
			go updateProfile(userEvent)
		}
		msg.Success()
	} else {
		log.Errorln(err)
	}
}

func updateProfile(userEvent *event.UserEvent) {
	var (
		err error
	)

	fmt.Println("ZHEER")
	//Create session for every request
	session := mgoSession.Copy()
	defer session.Close()

	err = db.UpdateUserProfile(session, userEvent.User)
	if err != nil {
		log.Errorln(err)
	}
	log.Infoln(userEvent.User.UserId, "update profile")
}
