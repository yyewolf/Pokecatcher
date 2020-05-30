package main

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type infoActivated struct {
	ChannelID        string
	Activated        bool
	MessageID        string
	AutoRelease      bool
	SelectedFromList bool
	SelectedIndex    string
}

//SelectVerifier will verify the level of the selected pokemon
func SelectVerifier(s *discordgo.Session, msg *discordgo.MessageCreate) {
	//Check if the person is allowed
	if !config.IsAllowedToUse {
		return
	}
	//Check if the author is pokecord
	if msg.Author.ID != botID {
		return
	}
	//Check if the person wants to autolevel
	if !config.AutoLeveling {
		return
	}
	//Check that the leveler is activated
	if !infoMenu.Activated {
		return
	}
	//Verify that it's the right channel
	if msg.ChannelID != infoMenu.ChannelID {
		return
	}
	//Check that it is not in release mode
	if infoMenu.AutoRelease {
		return
	}
	//Verifies that it's the answer to the user's message
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

	//Do auto leveler stuff
	reg, err := regexp.Compile("[^0-9]")
	if err != nil {
		return
	}
	Level := strings.Split(msg.Content, "N")[0]
	Level = reg.ReplaceAllString(Level, "")
	Name := strings.Split(msg.Content, ".")[0]
	max := len(strings.Split(Name, Level+" ")) - 1
	Name = strings.Split(Name, Level+" ")[max]

	if Level == "" {
		return
	}
	//Sets the current box
	currentPokemonImg.Image = decodedImages[Name]
	currentPokemonLevel.SetText("Level : " + Level)
	currentPokemonImg.Refresh()

	lvl, err := strconv.Atoi(Level)
	if err != nil {
		return
	}
	maxlvl, err := strconv.Atoi(config.AutoLevelMax)
	if err != nil {
		return
	}

	if lvl >= maxlvl {
		//Updates the level in the List
		if infoMenu.SelectedFromList {
			c := pokemonList[infoMenu.SelectedIndex]
			c.Level = config.AutoLevelMax
			pokemonList[infoMenu.SelectedIndex] = c
		}
		logDebug("[DEBUG] AutoLeveler is searching for a new pokemon.")
		time.Sleep(3 * time.Second)
		//Will verify the next pokemon's level
		m, err := s.ChannelMessageSend(config.ChannelID, config.PrefixPokecord+"info")
		if err != nil {
			logDebug("[ERROR] ", err)
			return
		}
		//Activate
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
	if msg.Author.ID != botID {
		return
	}
	//Check if the person wants to autolevel
	if !config.AutoLeveling {
		return
	}
	//Check that the leveler is activated
	if !infoMenu.Activated {
		return
	}
	//Check that it is not in release mode
	if infoMenu.AutoRelease {
		return
	}
	//Verify that it's the right channel
	if msg.ChannelID != infoMenu.ChannelID {
		return
	}
	//Verifies that it's the answer to the user's message
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
	logDebug("1")
	//Check if it's a p!info response
	Infos, err := parsePokemonInfo(msg)
	if err != nil {
		logDebug("[ERROR]", err)
		return
	}
	logDebug("2")
	//Sets the current box
	currentPokemonImg.Image = Infos.Image
	currentPokemonLevel.SetText("Level : " + Infos.Level)
	currentPokemonImg.Refresh()

	infoMenu.Activated = false
	infoMenu.MessageID = ""

	lvl, err := strconv.Atoi(Infos.Level)
	if err != nil {
		return
	}
	maxlvl, err := strconv.Atoi(config.AutoLevelMax)
	if err != nil {
		return
	}

	if lvl >= maxlvl {
		logDebug("[DEBUG] AutoLeveler sending p!select")
		Number := 1
		if !Infos.Last {
			Number, _ = strconv.Atoi(Infos.Number)
			Number++
		}
		time.Sleep(2 * time.Second)
		//Select the next pokemon
		n := strconv.Itoa(Number)
		//Tries to select a legendary if possible
		for i := range pokemonList {
			if pokemonList[i].Level == config.AutoLevelMax {
				continue
			}
			for j := range legendaries {
				if pokemonList[i].Name == legendaries[j] {
					n = pokemonList[i].NewNumber
					infoMenu.SelectedFromList = true
					infoMenu.SelectedIndex = i
					break
				}
			}
			break
		}
		//Sends the message
		m, err := s.ChannelMessageSend(config.ChannelID, config.PrefixPokecord+"select "+n)
		if err != nil {
			logDebug("[ERROR]", err)
			return
		}
		//Activates
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
	if msg.Author.ID != botID {
		return
	}
	//Check if the person wants to autolevel
	if !config.AutoLeveling {
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
	//Verifies that it's the answer to the user's message
	msgs, _ := s.ChannelMessages(msg.ChannelID, 5, msg.ID, "", "")
	ok := false
	for i := range msgs {
		if msgs[i].Author.ID == s.State.User.ID {
			ok = true
			break
		}
	}
	if !ok {
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
	currentPokemonLevel.SetText("Level : " + NewLevel)

	lvl, err := strconv.Atoi(NewLevel)
	if err != nil {
		return
	}
	maxlvl, err := strconv.Atoi(config.AutoLevelMax)
	if err != nil {
		return
	}
	if lvl >= maxlvl {
		time.Sleep(2 * time.Second)
		currentPokemonLevel.SetText("Level : Switching..")
		//Will prioritize priority queue over randomness
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
		//Activates
		infoMenu.ChannelID = m.ChannelID
		infoMenu.MessageID = m.ID
		infoMenu.Activated = true
	}
}
