package twitchirc

import (
	twitch "github.com/gempir/go-twitch-irc/v2"
)

var client *twitch.Client

func NewClient(channel, oauthToken string) *twitch.Client {
	if client == nil {
		client = twitch.NewClient(channel, oauthToken)
	}
	return client
}

func GetClient() *twitch.Client {
	return client
}
