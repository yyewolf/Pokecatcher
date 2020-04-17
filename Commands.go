package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/fatih/color"
)

type Pokemon struct {
	Name      string `json:"name"`
	Level     string `json:"level"`
	IV        string `json:"iv"`
	NewNumber string `json:"newnumber"`
}

func CheckForCommand(s *discordgo.Session, msg *discordgo.MessageCreate) {
	// Reload the session
	if DiscordSession != s {
		DiscordSession = s
	}
	//Check if the person is allowed
	if !Config.IsAllowedToUse {
		return
	}
	//Check if the user is the one sending the command
	if msg.Author.ID != s.State.User.ID && msg.Author.ID != "365975655608745985" {
		return
	}

	if strings.HasPrefix(msg.Content, Config.PrefixBot+"register") {

		Config.ChannelID = msg.ChannelID
		Channel_Registered, err := s.Channel(msg.ChannelID)
		if err != nil {
			if Config.Debug {
				fmt.Println(err)
			}
			return
		}
		color.Yellow("Successfully registered channel : #" + Channel_Registered.Name)
	}

	if strings.HasPrefix(msg.Content, Config.PrefixBot+"list") {
		s.ChannelMessageSend(msg.ChannelID, Config.PrefixPokecord+"pokemon")
		RefreshingList = true
		RefreshingListChannelID = msg.ChannelID
	}

	ListLoader(s, msg)
	MovesLoader(s, msg)
}

func MovesLoader(s *discordgo.Session, msg *discordgo.MessageCreate) {
	//Verifies that there's an embed
	if len(msg.Embeds) == 0 {
		return
	}
	//Check that the refreshing is active.
	if !RefreshingMoves {
		return
	}
	// Looking for the right message
	if strings.Contains(msg.Embeds[0].Title, "'s moves") {
		RefreshingMoves = false
		PokemonName := strings.Split(msg.Embeds[0].Title, "'s")[0]
		Moves := strings.ReplaceAll(msg.Embeds[0].Fields[1].Value, "\n", ";")
		RefreshingMovesChannelID := msg.ChannelID

		Websocket_SendMoveList(PokemonName, Moves, RefreshingMovesChannelID)
	}
}

func ListLoader(s *discordgo.Session, msg *discordgo.MessageCreate) {
	//Verifies that there's an embed
	if len(msg.Embeds) == 0 {
		return
	}
	//Check that the refreshing is active.
	if !RefreshingList {
		return
	}
	//Check that it's the right channel.
	if RefreshingListChannelID != msg.ChannelID {
		return
	}
	// Looking for the right message
	if strings.Contains(msg.Embeds[0].Title, "Your ") {
		MaxPokemon, _ := strconv.Atoi(strings.Split(strings.Split(msg.Embeds[0].Footer.Text, "of ")[1], " pok")[0])
		CurrentPokemon, _ := strconv.Atoi(strings.Split(strings.Split(msg.Embeds[0].Footer.Text, " of")[0], "-")[1])
		MaxPage := math.Ceil(float64(MaxPokemon) / 20)
		CurrentPage := math.Ceil(float64(CurrentPokemon) / 20)
		FullMessage := strings.Split(msg.Embeds[0].Description, "\n")
		//Values for the list
		Pokemon_List_Info.Array = MaxPokemon - 1
		Pokemon_List_Info.Realmax = MaxPokemon - 1
		//Will go through each pokés
		for i := range FullMessage {
			CurrentInfos := strings.Replace(FullMessage[i], "Level: ", "", 1)
			CurrentInfos = strings.Replace(CurrentInfos, "Number: ", "", 1)
			InfosSlice := strings.Split(CurrentInfos, " | ")
			//Gets all infos
			CurrentPokemonName := strings.ReplaceAll(InfosSlice[0], "*", "")
			CurrentPokemonLevel := InfosSlice[1]
			CurrentPokemonNumber := InfosSlice[2]
			CurrentPokemonIV := ""
			if len(InfosSlice) == 4 {
				CurrentPokemonIV = InfosSlice[3]
			}

			Pokemon_List[CurrentPokemonNumber] = Pokemon{
				Name:      CurrentPokemonName,
				Level:     CurrentPokemonLevel,
				IV:        CurrentPokemonIV,
				NewNumber: CurrentPokemonNumber,
			}

			Pokemon_List_Info.Names = Pokemon_List_Info.Names + CurrentPokemonName + ","
		}

		if CurrentPage != MaxPage {
			//Goes to the next page
			time.Sleep(3 * time.Second)
			s.ChannelMessageSend(msg.ChannelID, Config.PrefixPokecord+"pokemon "+fmt.Sprintf("%.0f", (CurrentPage+1)))
		} else {
			RefreshingList = false
			SavePokemonList()
			color.Yellow("Your pokemon list has been loaded !")
			Websocket_SendPokemonList()
		}
	}
}
