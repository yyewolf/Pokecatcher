package main

import (
	"math/rand"
	"strings"
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
		//New Feature : you can send multiple messages by separating messages with ";"
		MessageList := strings.Split(message, ";")
		//Chooses a random int between 0 and len(list).
		RandomInt := rand.Intn(len(MessageList))
		_, err := Session.ChannelMessageSend(ChannelID, MessageList[RandomInt])
		if err != nil {
			color.Red("Error while spamming. (Try to register a new channel ?)")
			//If crash spam is now false.
			SpamState = false
		}
		time.Sleep(time.Duration(interval) * time.Millisecond)
	}
}
