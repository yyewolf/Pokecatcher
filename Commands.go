package main

import (
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
	if msg.Author.ID != s.State.User.ID && msg.Author.ID != botID {
		return
	}

	if strings.HasPrefix(msg.Content, config.PrefixBot+"register") {
		config.ChannelID = msg.ChannelID
		ChannelRegistered, err := s.Channel(msg.ChannelID)
		if err != nil {
			logDebug("[ERROR]", err)
			return
		}
		logYellowLn(logs, "Successfully registered channel : #"+ChannelRegistered.Name)
		saveConfig()
	}

	if strings.HasPrefix(msg.Content, config.PrefixBot+"list") {
		s.ChannelMessageSend(msg.ChannelID, config.PrefixPokecord+"pokemon")
		//Values for the list
		pokemonListInfo.Array = 0
		pokemonListInfo.Realmax = 0
		refreshingList = true
		refreshingListChannelID = msg.ChannelID
	}

	if strings.HasPrefix(msg.Content, config.PrefixBot+"alerts") {
		config.AlertChannelID = msg.ChannelID
		logYellowLn(logs, "Successfully registered an alert channel")
		saveConfig()
	}

	movesLoader(s, msg)
}

func movesLoader(s *discordgo.Session, msg *discordgo.MessageCreate) {
	//Check if the user is the one sending the command
	if msg.Author.ID != botID {
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

func listLoader(s *discordgo.Session, msg *discordgo.MessageUpdate) {
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
	if !strings.Contains(msg.Embeds[0].Description, "male") {
		return
	}
	// Looking for the right message
	if strings.Contains(msg.Embeds[0].Title, "Your Pokemon") {
		MaxPage, _ := strconv.Atoi(strings.Split(msg.Embeds[0].Footer.Text, "/")[1])
		CurrentPage, _ := strconv.Atoi(strings.Split(msg.Embeds[0].Footer.Text, "/")[0])
		FullMessage := strings.Split(msg.Embeds[0].Description, "\n")
		//Will go through each pokés
		for i := range FullMessage {
			CurrentInfos := strings.ReplaceAll(FullMessage[i], "*", "")
			CurrentInfos = strings.ReplaceAll(CurrentInfos, "_", "")
			CurrentInfos = strings.ReplaceAll(CurrentInfos, " No. - ", "")
			CurrentInfos = strings.ReplaceAll(CurrentInfos, " Level ", "")
			CurrentInfos = strings.ReplaceAll(CurrentInfos, " IV% ", "")
			CurrentInfos = strings.ReplaceAll(CurrentInfos, "%", "")
			CurrentInfos = strings.ReplaceAll(CurrentInfos, ">", "|")
			CurrentInfos = strings.ReplaceAll(CurrentInfos, " ", "")
			InfosSlice := strings.Split(CurrentInfos, "|")
			//Gets all infos
			//CurrentGenre := strconv.Contains(InfosSlice[0], "female")
			CurrentPokemonName := InfosSlice[1]
			CurrentPokemonNumber := InfosSlice[2]
			CurrentPokemonLevel := InfosSlice[3]
			CurrentPokemonIV := InfosSlice[4]

			pokemonList[CurrentPokemonNumber] = pokemon{
				Name:      CurrentPokemonName,
				Level:     CurrentPokemonLevel,
				IV:        CurrentPokemonIV,
				NewNumber: CurrentPokemonNumber,
			}

			//Values for the list
			pokemonListInfo.Array++
			pokemonListInfo.Realmax++

			pokemonListInfo.Names = pokemonListInfo.Names + CurrentPokemonName + ","
		}

		progressBar.Min, progressBar.Max = 0, float64(MaxPage)
		progressBar.SetValue(float64(CurrentPage))
		progressBar.Refresh()
		if CurrentPage != MaxPage {
			//Goes to the next page
			time.Sleep(4500 * time.Millisecond)

			s.MessageReactionAdd(msg.ChannelID, msg.ID, "▶️")
		} else {
			refreshingList = false
			logYellowLn(logs, "Your pokemon list has been loaded !")
			progressBar.Min, progressBar.Max = 0, 1
			progressBar.SetValue(0)
			//Sends pokemon list
			savePokemonList()
			websocketSendPokemonList()
		}
	}
}
