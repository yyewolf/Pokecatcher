package main

import (
	"math/rand"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func spamFunc(Session *discordgo.Session, ChannelID string, interval int, message string) {
	defer func() {
		if r := recover(); r != nil {
			logRedLn(logs, "Error while spamming. Wait a bit and refresh the web page !")
		}
	}()
	spamInterval = interval
	spamMessage = message
	for {
		select {
		case <-spamChannel:
			return
		default:
			break
		}
		//You can send multiple messages by separating messages with ";"
		MessageList := strings.Split(message, ";")
		//Chooses a random int between 0 and len(list).
		RandomInt := rand.Intn(len(MessageList))
		_, err := Session.ChannelMessageSend(ChannelID, MessageList[RandomInt])
		if err != nil {
			logDebug("[ERROR] ", err)
			logRedLn(logs, "Error while spamming. (Try to register a new channel ?)")
			//If crash spam is now false.
			spamState = false
		}
		time.Sleep(time.Duration(interval) * time.Millisecond)
	}
}
