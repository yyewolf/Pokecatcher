package main

import (
	"fmt"
	"strconv"
	"time"
)

//////////////////////////////////////////////
//////////Funcs related to Refreshs///////////
//////////////////////////////////////////////

func RefreshPokemonList(Request Receive_Request) {
	// Active requests variables :
	// #######

	PrintYellowln("Refreshing your pokemon list...")
	if Config.ChannelID != "" {
		_, err := DiscordSession.ChannelMessageSend(Config.ChannelID, Config.PrefixBot+"list")
		if err != nil {
			if Config.Debug {
				fmt.Println(err)
			}
			PrintRedln("Couldn't send the message to start the reading of the list. (Try to register a new channel ?)")
		}
	} else {
		PrintRedln("No channel are registered, register one before trying again.")
	}
}

func RefreshPokemonMovesList(Request Receive_Request) {
	// Active requests variables :
	// #######

	PrintYellowln("Refreshing your pokemon's moves list...")
	if Config.ChannelID != "" {
		_, err := DiscordSession.ChannelMessageSend(Config.ChannelID, Config.PrefixPokecord+"moves")
		if err != nil {
			if Config.Debug {
				fmt.Println(err)
			}
			PrintRedln("Couldn't send the message to start the reading of the list. (Try to register a new channel ?)")
		} else {
			RefreshingMoves = true
		}
	} else {
		PrintRedln("No channel are registered, register one before trying again.")
	}
}

//////////////////////////////////////////////
//////////Funcs related to Settings///////////
//////////////////////////////////////////////

func UpdateSpammerSettings(Request Receive_Request) {
	// Active requests variables :
	// Request.State
	// Request.Message
	// Request.SpamInterval

	//Check if the person is allowed
	if !Config.IsAllowedToUse {
		return
	}

	SpamState = Request.State
	if SpamState {
		if Config.ChannelID != "" {
			SpamChannel = make(chan (int))
			PrintYellowln("The spammer has been started.")
			go SpamFunc(DiscordSession, Config.ChannelID, Request.SpamInterval, Request.Message)
			//Better than the JS equivalent c:
		} else {
			PrintRedln("No channel are registered, register one before trying again.")
		}
	} else {
		//This will cause the SpamFunc to return
		SpamChannel <- 1
	}
}

func UpdateToken(Request Receive_Request) {
	// Active requests variables :
	// Request.Token

	NoConfig()
	LoadConfig()
	Config.Token = Request.Token
	SaveConfig()
	PrintRedln("Restart the bot to apply your changes now. The bot may not work properly from now on.")
	//Not sure if it's true, don't have time to test it though.
}

func UpdateServerWhitelist(Request Receive_Request) {
	// Active requests variables :
	// Request.GuildID
	// Request.GuildState

	ServerWhitelist[Request.GuildID] = Request.GuildState
	SaveWhitelist()
	PrintYellowln("Your server whitelist has been successfully updated.")
}

func ChangePrefixes(Request Receive_Request) {
	// Active requests variables :
	// Request.Type
	// Request.Prefix
	// sb => self-bot ; pc => pokécord

	if Request.Type == "sb" {
		Config.PrefixBot = Request.Prefix
	} else if Request.Type == "pc" {
		Config.PrefixPokecord = Request.Prefix
	}
	SaveConfig()
	PrintYellowln("The prefix has been successfully updated.")
}

func ChangeDelay(Request Receive_Request) {
	// Active requests variables :
	// Request.Delay

	Config.Delay = Request.AutoCatchDelay
	SaveConfig()
	PrintYellowln("The delay for the autocatcher has been succesfully updated.")
}

func AutoCatcherOnOff(Request Receive_Request) {
	// Active requests variables :
	// Request.State

	Config.AutoCatching = Request.State
	PrintYellowln("Autocatching : " + strconv.FormatBool(Config.AutoCatching))
}

func DuplicatesOnOff(Request Receive_Request) {
	// Active requests variables :
	// Request.State

	Config.Duplicate = Request.State
	SaveConfig()
	PrintYellowln("Catching duplicates : " + strconv.FormatBool(Config.Duplicate))
}

func AliasesOnOff(Request Receive_Request) {
	// Active requests variables :
	// Request.State

	Config.Aliases = Request.State
	SaveConfig()
	PrintYellowln("Catching pokemons with aliases : " + strconv.FormatBool(Config.Aliases))
}

//////////////////////////////////////////////
//////////Funcs related to Pokemons///////////
//////////////////////////////////////////////

func Release(Request Receive_Request) {
	// Active requests variables :
	// Request.PokemonNumber

	//Check if the person is allowed
	if !Config.IsAllowedToUse {
		return
	}

	_, err := DiscordSession.ChannelMessageSend(Config.ChannelID, Config.PrefixPokecord+"release "+strconv.Itoa(Request.PokemonNumber))
	if err == nil {
		time.Sleep(3 * time.Second)
		_, err := DiscordSession.ChannelMessageSend(Config.ChannelID, Config.PrefixPokecord+"confirm")
		if err == nil {
			RemovePokemonFromList(Request)
		}
	} else {
		if Config.Debug {
			fmt.Println(err)
		}
		PrintRedln("Couldn't release your pokemon, check that you've registered a channel and try again.")
	}
}

func LearnNewMove(Request Receive_Request) {
	// Active requests variables :
	// Request.PokemonName
	// Request.MoveName
	// Request.MovePosition
	// Request.ChannelSet

	_, err := DiscordSession.Channel(Config.ChannelID)
	if err == nil {
		DiscordSession.ChannelMessageSend(Config.ChannelID, Config.PrefixPokecord+"learn "+Request.MoveName)
		time.Sleep(3 * time.Second)
		DiscordSession.ChannelMessageSend(Config.ChannelID, Config.PrefixPokecord+"replace "+strconv.Itoa(Request.MovePosition))
	} else {
		PrintRedln("You didn't register any channel ! Register one by using : '" + Config.PrefixBot + "register' in a channel.")
	}
}

func CatchAPokemon(Request Receive_Request) {
	// Active requests variables :
	// Request.ChannelID
	// Request.Command
	// Request.Name

	//Check if the person is allowed
	if !Config.IsAllowedToUse {
		return
	}

	Channel_Spawn, err := DiscordSession.Channel(Request.ChannelID)
	check(err)
	Guild_Spawn, err := DiscordSession.Guild(Channel_Spawn.GuildID)
	check(err)

	_, err = DiscordSession.ChannelMessageSend(Request.ChannelID, Request.Command+" "+Request.Name)
	if err != nil {
		Notif_CatchingErr(Request.Name, Guild_Spawn.Name, Channel_Spawn.Name)
		PrintRedln("There was a problem when trying to catch that pokemon, try again next time maybe ?")
	} else {
		PrintBlueln("Tried to catch your : " + Request.Name)
	}

}

func RenamePokemon(Request Receive_Request) {
	// Active requests variables :
	// Request.Nickname

	_, err := DiscordSession.ChannelMessageSend(Config.ChannelID, Config.PrefixPokecord+"nickname "+Request.Nickname)
	if err != nil {
		Notif_RenameFailed(Request.Nickname)
		PrintRedln("There was a problem when trying to rename your pokemon, try with another one maybe ?")
	} else {
		Notif_RenameSuccess(Request.Nickname)
		PrintBlueln("Successfully renamed your selected pokemon into " + Request.Nickname + " !")
	}
}

func RemovePokemonFromList(Request Receive_Request) {
	// Active requests variables :
	// Request.PokemonNumber

	PokemonNumber := Pokemon_List_Info.Array
	// Loops through all pokémons to lower their numbers and keep things ordered.
	for i := 0; i < PokemonNumber; i++ {
		Number := i + 1
		CurrentPoke := Pokemon_List[strconv.Itoa(Request.PokemonNumber)].(Pokemon) // Type insertion cuz interface{}
		PokemonNewNumber, _ := strconv.Atoi(CurrentPoke.NewNumber)
		if Number > PokemonNewNumber {
			CurrentPoke.NewNumber = strconv.Itoa(i)
			Pokemon_List[strconv.Itoa(Request.PokemonNumber)] = CurrentPoke
		}
	}
	// realmax corresponds to the amount of pokemon present.
	Pokemon_List_Info.Realmax = PokemonNumber - 1
	delete(Pokemon_List, strconv.Itoa(Request.PokemonNumber))
	PrintBlueln("Removed the pokemon from your Pokemon List.")
	// Sends info to the websocket to update the list.
	Websocket_RemovedFromList(Request.PokemonNumber)
}

func SelectPokemon(Request Receive_Request) {
	// Active requests variables :
	// Request.Name
	// Request.Number
	// Request.Name

	//Check if the person is allowed
	if !Config.IsAllowedToUse {
		return
	}

	PokemonNumberOnWebsite := strconv.Itoa(Request.PokemonNumber)
	PokemonNumber := "0"
	switch Pokemon_List[PokemonNumberOnWebsite].(type) {
	case Pokemon:
		PokemonNumber = Pokemon_List[PokemonNumberOnWebsite].(Pokemon).NewNumber
	case map[string]interface{}:
		PokemonNumber = Pokemon_List[PokemonNumberOnWebsite].(map[string]interface{})["newnumber"].(string)
	default:
		return
	}

	_, err := DiscordSession.ChannelMessageSend(Config.ChannelID, Config.PrefixPokecord+"select "+PokemonNumber) //Type insertion because it is an interface{} type
	if err != nil {
		if Config.Debug {
			fmt.Println(err)
		}
		PrintRedln("There was a problem when trying to select the pokemon, try with another one maybe ?")
	} else {
		SelectedPokemon.Number, _ = strconv.Atoi(PokemonNumber) //Type insertion (again) because it is an interface{} type
		SelectedPokemon.Name = Request.Name
		Websocket_Selected(SelectedPokemon.Name, SelectedPokemon.Number)
	}
}

func Websocket_Receive_AllFunctions() {
	// Directly translated from the old code.
	Websocket_Receive_Functions["aca"] = AutoCatcherOnOff
	Websocket_Receive_Functions["duplicate"] = DuplicatesOnOff
	Websocket_Receive_Functions["aliases"] = AliasesOnOff
	Websocket_Receive_Functions["refresh"] = RefreshPokemonList
	Websocket_Receive_Functions["refreshmoves"] = RefreshPokemonMovesList
	Websocket_Receive_Functions["release"] = Release
	Websocket_Receive_Functions["select"] = SelectPokemon
	Websocket_Receive_Functions["remove"] = RemovePokemonFromList
	Websocket_Receive_Functions["nickname"] = RenamePokemon
	Websocket_Receive_Functions["tokenchange"] = UpdateToken
	Websocket_Receive_Functions["spam"] = UpdateSpammerSettings
	Websocket_Receive_Functions["catch"] = CatchAPokemon
	Websocket_Receive_Functions["learn"] = LearnNewMove
	Websocket_Receive_Functions["whitelist"] = UpdateServerWhitelist
	Websocket_Receive_Functions["autodelaychange"] = RefreshPokemonMovesList
	Websocket_Receive_Functions["prefixchange"] = ChangePrefixes
	Websocket_Receive_Functions["autodelaychange"] = ChangeDelay
}
