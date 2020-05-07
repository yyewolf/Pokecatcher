package main

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type pokemon struct {
	Name      string `json:"name"`
	Level     string `json:"level"`
	IV        string `json:"iv"`
	NewNumber string `json:"newnumber"`
}

type latestPokemonType struct {
	ChannelID string
	Command   string
}

func catchLatest() {
	_, err := discordSession.ChannelMessageSend(latestPokemon.ChannelID, latestPokemon.Command)
	if err != nil {
		logRedLn(logs, "There was a problem when trying to catch that pokemon, try again next time maybe ?")
	} else {
		logBlueLn(logs, "Tried to catch latest Pokemon.")
	}
}

func checkForCommand(s *discordgo.Session, msg *discordgo.MessageCreate) {
	// Reload the session
	if discordSession != s {
		discordSession = s
	}
	//Check if the person is allowed
	if !config.IsAllowedToUse {
		return
	}
	//Check if the user is the one sending the command
	if msg.Author.ID != s.State.User.ID && msg.Author.ID != "365975655608745985" {
		return
	}

	if strings.HasPrefix(msg.Content, config.PrefixBot+"register") {

		config.ChannelID = msg.ChannelID
		ChannelRegistered, err := s.Channel(msg.ChannelID)
		if err != nil {
			logDebug("[ERROR] ", err)
			return
		}
		logYellowLn(logs, "Successfully registered channel : #"+ChannelRegistered.Name)
	}

	if strings.HasPrefix(msg.Content, config.PrefixBot+"list") {
		s.ChannelMessageSend(msg.ChannelID, config.PrefixPokecord+"pokemon")
		refreshingList = true
		refreshingListChannelID = msg.ChannelID
	}

	listLoader(s, msg)
	movesLoader(s, msg)
}

func movesLoader(s *discordgo.Session, msg *discordgo.MessageCreate) {
	//Check if the user is the one sending the command
	if msg.Author.ID != "365975655608745985" {
		return
	}
	//Verifies that there's an embed
	if len(msg.Embeds) == 0 {
		return
	}
	//Check that the refreshing is active.
	if !refreshingMoves {
		return
	}
	//Check if there is a footer in the embed
	if msg.Embeds[0].Footer == nil {
		return
	}
	// Looking for the right message
	if strings.Contains(msg.Embeds[0].Footer.Text, "Moves") {
		refreshingMoves = false
		//Gets the level from the embed's title : "Level 99 Pokemon"
		reg, _ := regexp.Compile("[^0-9/]")
		PokeLevel := reg.ReplaceAllString(msg.Embeds[0].Title, "")
		//Gets the name using the level
		PokemonName := strings.Split(msg.Embeds[0].Title, PokeLevel+" ")[1]
		Moves := strings.ReplaceAll(msg.Embeds[0].Fields[1].Value, "\n", ";")
		Moves = strings.ReplaceAll(Moves, " ;", ";")
		RefreshingMovesChannelID := msg.ChannelID

		websocketSendMoveList(PokemonName, Moves, RefreshingMovesChannelID)
	}
}

func listLoader(s *discordgo.Session, msg *discordgo.MessageCreate) {
	//Verifies that there's an embed
	if len(msg.Embeds) == 0 {
		return
	}
	//Check that the refreshing is active.
	if !refreshingList {
		return
	}
	//Check that it's the right channel.
	if refreshingListChannelID != msg.ChannelID {
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
		pokemonListInfo.Array = MaxPokemon
		pokemonListInfo.Realmax = MaxPokemon
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
			if len(InfosSlice) >= 4 {
				CurrentPokemonIV = InfosSlice[3]
			}

			pokemonList[CurrentPokemonNumber] = pokemon{
				Name:      CurrentPokemonName,
				Level:     CurrentPokemonLevel,
				IV:        CurrentPokemonIV,
				NewNumber: CurrentPokemonNumber,
			}

			pokemonListInfo.Names = pokemonListInfo.Names + CurrentPokemonName + ","
		}

		progressBar.Min, progressBar.Max = 0, MaxPage
		progressBar.SetValue(CurrentPage)
		progressBar.Refresh()
		if CurrentPage != MaxPage {
			//Goes to the next page
			time.Sleep(4500 * time.Millisecond)
			s.ChannelMessageSend(msg.ChannelID, config.PrefixPokecord+"pokemon "+fmt.Sprintf("%.0f", (CurrentPage+1)))
		} else {
			refreshingList = false
			savePokemonList()
			logYellowLn(logs, "Your pokemon list has been loaded !")
			websocketSendPokemonList()
			progressBar.Min, progressBar.Max = 0, 1
			progressBar.SetValue(0)
		}
	}
}
