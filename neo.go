package main

import (
	"neo/config"
	"neo/consumer"
	"neo/handlers"
)

var (
	configuration config.Config
)

func init() {

	configuration = config.New()

}

func main() {
	consumer.Register("feed", "consume", 30, handlers.FeedHandler)
	consumer.Register("friend", "consume", 30, handlers.FriendHandler)
	consumer.Register("tag", "consume", 30, handlers.TagHandler)
	consumer.Register("user", "consume", 30, handlers.UserHandler)
	consumer.Connect(configuration.NSQ.Host)
	consumer.Start(true)
}
