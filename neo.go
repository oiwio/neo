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
	consumer.Connect(configuration.NSQ.Host)
	consumer.Start(true)
}
