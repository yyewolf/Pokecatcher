package main

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/fatih/color"
)

func SpamFunc(Session *discordgo.Session, ChannelID string, interval int, message string) {
	SpamInterval = interval
	SpamMessage = message
	for {
		select {
		case <-SpamChannel:
			return
		default:
			break
		}
		_, err := Session.ChannelMessageSend(ChannelID, message)
		if err != nil {
			color.Red("Error while spamming. (Try to register a new channel ?)")
			SpamState = false
		}
		time.Sleep(time.Duration(interval) * time.Millisecond)
	}
}
