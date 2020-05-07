package main

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type infoActivated struct {
	ChannelID   string
	Activated   bool
	MessageID   string
	AutoRelease bool
}

//SelectVerifier will verify the level of the selected pokemon
func SelectVerifier(s *discordgo.Session, msg *discordgo.MessageCreate) {
	//Check if the person is allowed
	if !config.IsAllowedToUse {
		return
	}
	//Check if the author is pokecord
	if msg.Author.ID != "365975655608745985" {
		return
	}
	if !infoMenu.Activated {
		return
	}
	if infoMenu.AutoRelease {
		return
	}
	if msg.ChannelID != infoMenu.ChannelID {
		return
	}
	msgs, _ := s.ChannelMessages(msg.ChannelID, 5, msg.ID, "", "")
	ok := false
	for i := range msgs {
		if msgs[i].ID == infoMenu.MessageID {
			ok = true
			break
		}
	}
	if !ok {
		return
	}

	reg, err := regexp.Compile("[^0-9]")
	if err != nil {
		return
	}
	Level := strings.Split(msg.Content, "N")[0]
	Level = reg.ReplaceAllString(Level, "")

	if Level == "100" {
		logDebug("[DEBUG] AutoLeveler is searching for a new pokemon.")
		time.Sleep(3 * time.Second)
		//Will verify the next pokemon's level
		m, err := s.ChannelMessageSend(config.ChannelID, config.PrefixPokecord+"info")
		if err != nil {
			logDebug("[ERROR] ", err)
			return
		}
		infoMenu.ChannelID = m.ChannelID
		infoMenu.MessageID = m.ID
		infoMenu.Activated = true
	} else {
		logCyanLn(logs, "Autoleveler selected a new Pokemon.")
	}
}

//InfoVerifier will verify the infos of the current pokemon
func InfoVerifier(s *discordgo.Session, msg *discordgo.MessageCreate) {
	//Check if the person is allowed
	if !config.IsAllowedToUse {
		return
	}
	//Check if the author is pokecord
	if msg.Author.ID != "365975655608745985" {
		return
	}
	if !infoMenu.Activated {
		return
	}
	if infoMenu.AutoRelease {
		return
	}
	if msg.ChannelID != infoMenu.ChannelID {
		return
	}
	msgs, _ := s.ChannelMessages(msg.ChannelID, 5, msg.ID, "", "")
	ok := false
	for i := range msgs {
		if msgs[i].ID == infoMenu.MessageID {
			ok = true
			break
		}
	}
	if !ok {
		return
	}
	//Check if it's a p!info response
	Infos, err := parsePokemonInfo(msg)
	if err != nil {
		return
	}

	infoMenu.Activated = false

	if Infos.Level == "100" {
		logDebug("[DEBUG] AutoLeveler sending p!select")
		Number := 1
		if !Infos.Last {
			Number, _ = strconv.Atoi(Infos.Number)
			Number++
		}
		time.Sleep(2 * time.Second)
		//Select the next pokemon
		n := strconv.Itoa(Number)
		m, err := s.ChannelMessageSend(config.ChannelID, config.PrefixPokecord+"select "+n)
		if err != nil {
			logDebug("[ERROR]", err)
			return
		}
		infoMenu.ChannelID = m.ChannelID
		infoMenu.MessageID = m.ID
		infoMenu.Activated = true
	}
}

//AutoLeveler allows the bot to automatically level every pokemon possible.
func AutoLeveler(s *discordgo.Session, msg *discordgo.MessageCreate) {
	//Check if the person is allowed
	if !config.IsAllowedToUse {
		return
	}
	//Check if the author is pokecord
	if msg.Author.ID != "365975655608745985" {
		return
	}
	//Check if there is an embed
	if len(msg.Embeds) == 0 {
		return
	}
	//Check if there is a title containing "Congratulations"
	if !strings.Contains(msg.Embeds[0].Title, "Congratulations") {
		return
	}

	//Gets the level from the embed's content : "Your *Pokemon* is now level 40!"
	//Steps :
	//40
	reg, err := regexp.Compile("[^0-9]")
	if err != nil {
		return
	}
	NewLevel := reg.ReplaceAllString(msg.Embeds[0].Description, "")

	if NewLevel == "100" {
		time.Sleep(2 * time.Second)
		if len(priorityQueue) != 0 {
			logDebug("[DEBUG] AutoLeveler sending p!select")
			n := priorityQueue[0]
			m, err := s.ChannelMessageSend(config.ChannelID, config.PrefixPokecord+"select "+n)
			if err != nil {
				logDebug("[ERROR]", err)
				return
			}
			infoMenu.ChannelID = m.ChannelID
			infoMenu.MessageID = m.ID
			infoMenu.Activated = true
			priorityQueue = priorityQueue[1:]
			return
		}

		logDebug("[DEBUG] AutoLeveler sending p!info")
		m, err := s.ChannelMessageSend(config.ChannelID, config.PrefixPokecord+"info")
		if err != nil {
			logDebug("[ERROR]", err)
			return
		}
		infoMenu.ChannelID = m.ChannelID
		infoMenu.MessageID = m.ID
		infoMenu.Activated = true
	}
}
