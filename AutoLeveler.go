package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type InfoActivated struct {
	ChannelID string
	Activated bool
	MessageID string
}

//InfoActivator will activate the verification
func InfoActivator(s *discordgo.Session, msg *discordgo.MessageCreate) {
	//Check if the person is allowed
	if !Config.IsAllowedToUse {
		return
	}
	//Check if the user is the one sending the command
	if msg.Author.ID != s.State.User.ID {
		return
	}
	//Check if the command is p!info
	if msg.Content != Config.PrefixPokecord+"info" && !strings.Contains(msg.Content, Config.PrefixPokecord+"select")  {
		return
	}

	InfoMenu.ChannelID = msg.ChannelID
	InfoMenu.Activated = true
	InfoMenu.MessageID = msg.ID

}

//SelectVerifier will verify the level of the selected pokemon
func SelectVerifier(s *discordgo.Session, msg *discordgo.MessageCreate) {
	//Check if the person is allowed
	if !Config.IsAllowedToUse {
		return
	}
	//Check if the author is pokecord
	if msg.Author.ID != "365975655608745985" {
		return
	}
	if !InfoMenu.Activated {
		return
	}
	if msg.ChannelID != InfoMenu.ChannelID {
		return
	}
	msgs, _ := s.ChannelMessages(msg.ChannelID, 5, msg.ID, "", "")
	ok := false
	for i := range msgs {
		if msgs[i].ID == InfoMenu.MessageID {
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
		time.Sleep(3 * time.Second)
		//Will verify the next pokemon's level
		_, err = DiscordSession.ChannelMessageSend(Config.ChannelID, Config.PrefixPokecord+"info")
		if err != nil {
			if Config.Debug {
				fmt.Println(err)
			}
		}
	} else {
		LogCyanLn(Logs, "Autoleveler selected a new Pokemon.")
	}
	
}

//InfoVerifier will verify the infos of the current pokemon
func InfoVerifier(s *discordgo.Session, msg *discordgo.MessageCreate) {
	//Check if the person is allowed
	if !Config.IsAllowedToUse {
		return
	}
	//Check if the author is pokecord
	if msg.Author.ID != "365975655608745985" {
		return
	}
	if !InfoMenu.Activated {
		return
	}
	if msg.ChannelID != InfoMenu.ChannelID {
		return
	}
	msgs, _ := s.ChannelMessages(msg.ChannelID, 5, msg.ID, "", "")
	ok := false
	for i := range msgs {
		if msgs[i].ID == InfoMenu.MessageID {
			ok = true
			break
		}
	}
	if !ok {
		return
	}
	//Check if it's a p!info response
	Infos, err := ParsePokemonInfo(msg)
	if err != nil {
		return
	}

	InfoMenu.Activated = false

	if Infos.isInList {
		Current := Pokemon_List[Infos.ListNumber]
		Current.Level = Infos.Level
		Pokemon_List[Infos.ListNumber] = Current
		SavePokemonList()
	}

	if Infos.Level == "100" {
		Number := 0
		if !Infos.Last {
			Number, _ = strconv.Atoi(Infos.Number)
			Number++
		}
		time.Sleep(2 * time.Second)
		//Select the next pokemon
		n := strconv.Itoa(Number)
		_, err := DiscordSession.ChannelMessageSend(Config.ChannelID, Config.PrefixPokecord+"select "+n)
		if err != nil {
			if Config.Debug {
				fmt.Println(err)
			}
		}
	}
}

//AutoLeveler allows the bot to automatically level every pokemon possible.
func AutoLeveler(s *discordgo.Session, msg *discordgo.MessageCreate) {
	//Check if the person is allowed
	if !Config.IsAllowedToUse {
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
		_, err := DiscordSession.ChannelMessageSend(Config.ChannelID, Config.PrefixPokecord+"info")
		if err != nil {
			if Config.Debug {
				fmt.Println(err)
			}
		}
	}
}
