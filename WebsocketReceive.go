package main

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

//////////////////////////////////////////////
//////////Funcs related to Refreshs///////////
//////////////////////////////////////////////

func refreshPokemonList(Request receiveRequest) {
	// Active requests variables :
	// #######

	//Resets everything
	pokemonList = make(map[string]pokemon) //Where the Pokemon List of the user will be stored.
	noPokemonList()
	loadPokemonList()

	logYellowLn(logs, "Refreshing your pokemon list...")
	if config.ChannelID != "" {
		_, err := discordSession.ChannelMessageSend(config.ChannelID, config.PrefixBot+"list")
		if err != nil {
			logDebug("[ERROR]", err)
			logRedLn(logs, "Couldn't send the message to start the reading of the list. (Try to register a new channel ?)")
		}
	} else {
		logRedLn(logs, "No channel are registered, register one before trying again.")
	}
}

func refreshPokemonMovesList(Request receiveRequest) {
	// Active requests variables :
	// #######

	logYellowLn(logs, "Refreshing your pokemon's moves list...")
	if config.ChannelID != "" {
		_, err := discordSession.ChannelMessageSend(config.ChannelID, config.PrefixPokecord+"moves")
		if err != nil {
			logDebug("[ERROR]", err)
			logRedLn(logs, "Couldn't send the message to start the reading of the list. (Try to register a new channel ?)")
		} else {
			refreshingMoves = true
		}
	} else {
		logRedLn(logs, "No channel are registered, register one before trying again.")
	}
}

func parsePriorityQueue(Request receiveRequest) {
	// Active requests variables :
	// Request.Change

	priorityQueue = strings.Split(Request.Change, ";")
	//Will now verify the Queue
	Keep := true
	for i := range priorityQueue {
		reg, _ := regexp.Compile("[^0-9]")
		if reg.MatchString(priorityQueue[i]) {
			Keep = false
			break
		}
	}

	if !Keep {
		priorityQueue = []string{}
		logDebug("[ERROR] PriorityQueue string malformed.")
		logRedLn(logs, "To make a Priority Queue you must enter numbers separated by semicolons.")
		logRedLn(logs, "Example : '25;14;36'")
		return
	}
	logDebug("[DEBUG] PriorityQueue has been parsed successfully.")
	logCyanLn(logs, "Your priority queue has been taken into account.")
}

//////////////////////////////////////////////
//////////Funcs related to Settings///////////
//////////////////////////////////////////////

func updateSpammerSettings(Request receiveRequest) {
	// Active requests variables :
	// Request.State
	// Request.Message
	// Request.SpamInterval

	//Check if the person is allowed
	if !config.IsAllowedToUse {
		return
	}

	spamState = Request.State
	if spamState {
		if config.ChannelID != "" {
			spamChannel = make(chan (int))
			logYellowLn(logs, "The spammer has been started.")
			go spamFunc(discordSession, config.ChannelID, Request.SpamInterval, Request.Message)
			//Better than the JS equivalent c:
		} else {
			logRedLn(logs, "No channel are registered, register one before trying again.")
		}
	} else {
		//This will cause the SpamFunc to return
		spamChannel <- 1
	}
}

func updateToken(Request receiveRequest) {
	// Active requests variables :
	// Request.Token

	noConfig()
	loadConfig()
	config.Token = Request.Token
	saveConfig()
	discordLogin()
	//Will login instead of prompting the user to relaunch the program
}

func updateServerWhitelist(Request receiveRequest) {
	// Active requests variables :
	// Request.GuildID
	// Request.GuildState

	serverWhitelist[Request.GuildID] = Request.GuildState
	saveWhitelist()
	logYellowLn(logs, "Your server whitelist has been successfully updated.")
}

func updatePokemonWhitelist(Request receiveRequest) {
	// Active requests variables :
	// Request.Name
	// Request.State

	pokemonWhitelist[Request.Name] = Request.State
	savePokemonWhitelist()
	logYellowLn(logs, "Your Pokémon whitelist has been successfully updated.")
}

func changePrefixes(Request receiveRequest) {
	// Active requests variables :
	// Request.Type
	// Request.Prefix
	// sb => self-bot's prefix ; pc => pokécord's prefix

	if Request.Type == "sb" {
		config.PrefixBot = Request.Prefix
	} else if Request.Type == "pc" {
		config.PrefixPokecord = Request.Prefix
	}
	saveConfig()
	logYellowLn(logs, "The prefix has been successfully updated.")
}

func changeDelay(Request receiveRequest) {
	// Active requests variables :
	// Request.Delay

	config.Delay = Request.AutoCatchDelay
	saveConfig()
	logYellowLn(logs, "The delay for the autocatcher has been succesfully updated.")
}

func autoCatcherOnOff(Request receiveRequest) {
	// Active requests variables :
	// Request.State

	config.AutoCatching = Request.State
	logYellowLn(logs, "Autocatching : "+strconv.FormatBool(config.AutoCatching))
}

func duplicatesOnOff(Request receiveRequest) {
	// Active requests variables :
	// Request.State

	config.Duplicate = Request.State
	saveConfig()
	logYellowLn(logs, "Ignoring duplicates : "+strconv.FormatBool(config.Duplicate))
}

func aliasesOnOff(Request receiveRequest) {
	// Active requests variables :
	// Request.State

	config.Aliases = Request.State
	saveConfig()
	logYellowLn(logs, "Catching pokemons with aliases : "+strconv.FormatBool(config.Aliases))
}

func filterOnOff(Request receiveRequest) {
	// Active requests variables :
	// Request.State

	config.GoodFilter = Request.State
	saveConfig()
	logYellowLn(logs, "Filter active : "+strconv.FormatBool(config.GoodFilter))
}

func updateFilters(Request receiveRequest) {
	// Active requests variables :
	// Request.Filters

	config.EveryFilters = Request.Filters
	saveConfig()
	logYellowLn(logs, "Custom filters were successfully registered !")
	//Will save the configuration to file
}

func customFilterOnOff(Request receiveRequest) {
	// Active requests variables :
	// Request.State

	config.CustomFilters = Request.State
	saveConfig()
	logYellowLn(logs, "Custom filters active : "+strconv.FormatBool(config.CustomFilters))
}

//////////////////////////////////////////////
//////////Funcs related to Pokemons///////////
//////////////////////////////////////////////

func release(Request receiveRequest) {
	// Active requests variables :
	// Request.PokemonNumber

	//Check if the person is allowed
	if !config.IsAllowedToUse {
		return
	}

	_, err := discordSession.ChannelMessageSend(config.ChannelID, config.PrefixPokecord+"release "+strconv.Itoa(Request.PokemonNumber))
	if err == nil {
		time.Sleep(3 * time.Second)
		_, _ = discordSession.ChannelMessageSend(config.ChannelID, config.PrefixPokecord+"confirm")
		removePokemonFromList(Request)
	} else {
		logDebug("[ERROR]", err)
		logRedLn(logs, "Couldn't release your pokemon, check that you've registered a channel and try again.")
	}
}

func learnNewMove(Request receiveRequest) {
	// Active requests variables :
	// Request.PokemonName
	// Request.MoveName
	// Request.MovePosition
	// Request.ChannelSet

	_, err := discordSession.Channel(config.ChannelID)
	if err == nil {
		discordSession.ChannelMessageSend(config.ChannelID, config.PrefixPokecord+"learn "+Request.MoveName)
		time.Sleep(3 * time.Second)
		discordSession.ChannelMessageSend(config.ChannelID, config.PrefixPokecord+"replace "+strconv.Itoa(Request.MovePosition))
	} else {
		logRedLn(logs, "You didn't register any channel ! Register one by using : '"+config.PrefixBot+"register' in a channel.")
	}
}

func catchAPokemon(Request receiveRequest) {
	// Active requests variables :
	// Request.ChannelID
	// Request.Command
	// Request.Name

	//Check if the person is allowed
	if !config.IsAllowedToUse {
		return
	}

	ChannelSpawn, err := discordSession.Channel(Request.ChannelID)
	check(err)
	GuildSpawn, err := discordSession.Guild(ChannelSpawn.GuildID)
	check(err)

	_, err = discordSession.ChannelMessageSend(Request.ChannelID, Request.Command+" "+Request.Name)
	if err != nil {
		notifCatchingErr(Request.Name, GuildSpawn.Name, ChannelSpawn.Name)
		logRedLn(logs, "There was a problem when trying to catch that pokemon, try again next time maybe ?")
	} else {
		logBlueLn(logs, "Tried to catch your : "+Request.Name)
	}

}

func renamePokemon(Request receiveRequest) {
	// Active requests variables :
	// Request.Nickname

	_, err := discordSession.ChannelMessageSend(config.ChannelID, config.PrefixPokecord+"nickname "+Request.Nickname)
	if err != nil {
		notifRenameFailed(Request.Nickname)
		logRedLn(logs, "There was a problem when trying to rename your pokemon, try with another one maybe ?")
	} else {
		notifRenameSuccess(Request.Nickname)
		logBlueLn(logs, "Successfully renamed your selected pokemon into "+Request.Nickname+" !")
	}
}

func removePokemonFromList(Request receiveRequest) {
	// Active requests variables :
	// Request.PokemonNumber

	PokemonNumber := pokemonListInfo.Array
	// Loops through all pokémons to lower their numbers and keep things ordered.
	RemovedPoke := pokemonList[strconv.Itoa(Request.PokemonNumber)]
	pokemonListInfo.Names = ""
	for i := 0; i < PokemonNumber; i++ {
		Number := i + 1
		CurrentPoke := pokemonList[strconv.Itoa(Number)]
		PokemonNewNumber, _ := strconv.Atoi(RemovedPoke.NewNumber)
		if Number > PokemonNewNumber && CurrentPoke.Name != "" {
			CurrentPoke.NewNumber = strconv.Itoa(i)
			pokemonList[strconv.Itoa(Number)] = CurrentPoke
		}
		pokemonListInfo.Names = pokemonListInfo.Names + CurrentPoke.Name + ","
	}
	// realmax corresponds to the amount of pokemon present.
	pokemonListInfo.Realmax = pokemonListInfo.Realmax - 1
	delete(pokemonList, strconv.Itoa(Request.PokemonNumber))
	logBlueLn(logs, "Removed the pokemon from your Pokemon List.")
	// Sends info to the websocket to update the list.
	go savePokemonList()
	websocketRemovedFromList(Request.PokemonNumber)
}

func selectPokemon(Request receiveRequest) {
	// Active requests variables :
	// Request.Name
	// Request.Number
	// Request.Name

	//Check if the person is allowed
	if !config.IsAllowedToUse {
		return
	}

	PokemonNumberOnWebsite := strconv.Itoa(Request.PokemonNumber)
	PokemonNumber := pokemonList[PokemonNumberOnWebsite].NewNumber

	_, err := discordSession.ChannelMessageSend(config.ChannelID, config.PrefixPokecord+"select "+PokemonNumber) //Type insertion because it is an interface{} type
	if err != nil {
		logDebug("[ERROR]", err)
		logRedLn(logs, "There was a problem when trying to select the pokemon, try with another one maybe ?")
	} else {
		selectedPokemon.Number, _ = strconv.Atoi(PokemonNumber) //Type insertion (again) because it is an interface{} type
		selectedPokemon.Name = Request.Name
		websocketSelected(selectedPokemon.Name, selectedPokemon.Number)
	}
}

func whitelistAllChecked(Request receiveRequest) {
	// Active requests variables :
	//

	for current := range pokemonWhitelist {
		pokemonWhitelist[current] = true
	}
	savePokemonWhitelist()
}

func whitelistAllUnchecked(Request receiveRequest) {
	// Active requests variables :
	//

	for current := range pokemonWhitelist {
		pokemonWhitelist[current] = false
	}
	savePokemonWhitelist()
}

func whitelistLegendChecked(Request receiveRequest) {
	// Active requests variables :
	//

	for current := range pokemonWhitelist {
		//Pokemon is legendary
		for i := range legendaries {
			if current == legendaries[i] {
				pokemonWhitelist[current] = true
			}
		}
	}
	savePokemonWhitelist()
}

func websocketReceiveAllFunctions() {
	// Directly translated from the old code.
	websocketReceiveFunctions["aca"] = autoCatcherOnOff
	websocketReceiveFunctions["duplicate"] = duplicatesOnOff
	websocketReceiveFunctions["aliases"] = aliasesOnOff
	websocketReceiveFunctions["filter"] = filterOnOff
	websocketReceiveFunctions["prefixchange"] = changePrefixes
	websocketReceiveFunctions["autodelaychange"] = changeDelay
	websocketReceiveFunctions["tokenchange"] = updateToken
	websocketReceiveFunctions["queuelist"] = parsePriorityQueue
	websocketReceiveFunctions["customfilters"] = customFilterOnOff
	websocketReceiveFunctions["filterschange"] = updateFilters

	websocketReceiveFunctions["refresh"] = refreshPokemonList
	websocketReceiveFunctions["refreshmoves"] = refreshPokemonMovesList
	websocketReceiveFunctions["whitelist"] = updateServerWhitelist
	websocketReceiveFunctions["pokemonwhitelist"] = updatePokemonWhitelist
	websocketReceiveFunctions["spam"] = updateSpammerSettings

	websocketReceiveFunctions["release"] = release
	websocketReceiveFunctions["select"] = selectPokemon
	websocketReceiveFunctions["remove"] = removePokemonFromList
	websocketReceiveFunctions["nickname"] = renamePokemon
	websocketReceiveFunctions["catch"] = catchAPokemon
	websocketReceiveFunctions["learn"] = learnNewMove

	websocketReceiveFunctions["pkmnwhitelistallchecked"] = whitelistAllChecked
	websocketReceiveFunctions["pkmnwhitelistallunchecked"] = whitelistAllUnchecked
	websocketReceiveFunctions["pkmnwhitelistlegendchecked"] = whitelistLegendChecked
}
