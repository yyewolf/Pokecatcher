package main

import "encoding/json"

func notifDuplicateErr(PokemonName string, GuildName string, ChannelName string) {
	// Message in websocket :
	// action = warn
	// message = couldnt-duplicate

	Notification, err := json.Marshal(sendRequest{
		Action:      "warn",
		Message:     "couldnt-duplicate",
		Name:        PokemonName,
		GuildName:   GuildName,
		ChannelName: ChannelName,
	})
	if err != nil {
		logDebug("[ERROR]", err)
		return
	}
	websocketBroadcast(string(Notification))
}

func notifCatchingErr(PokemonName string, GuildName string, ChannelName string) {
	// Message in websocket :
	// action = warn
	// message = couldnt-normal

	Notification, err := json.Marshal(sendRequest{
		Action:      "notification",
		Message:     "couldnt-normal",
		Name:        PokemonName,
		GuildName:   GuildName,
		ChannelName: ChannelName,
	})
	if err != nil {
		logDebug("[ERROR]", err)
		return
	}
	websocketBroadcast(string(Notification))
}

func notifPokeSpawn(PokemonName string, GuildName string, Command string, ChannelName string, ChannelID string) {
	// Message in websocket :
	// action = notification

	Notification, err := json.Marshal(sendRequest{
		Action:      "notification",
		Command:     Command,
		Name:        PokemonName,
		GuildName:   GuildName,
		ChannelName: ChannelName,
		ChannelID:   ChannelID,
	})
	if err != nil {
		logDebug("[ERROR]", err)
		return
	}
	websocketBroadcast(string(Notification))
}

func notifPokeCaught(PokemonName string, GuildName string, ChannelName string) {
	// Message in websocket :
	// action = warn
	// message = could

	Notification, err := json.Marshal(sendRequest{
		Action:      "warn",
		Message:     "could",
		Name:        PokemonName,
		GuildName:   GuildName,
		ChannelName: ChannelName,
	})
	if err != nil {
		logDebug("[ERROR]", err)
		return
	}
	websocketBroadcast(string(Notification))
}

func notifChannelRegistered() {
	// Message in websocket :
	// action = registered

	Notification, err := json.Marshal(sendRequest{
		Action: "registered",
	})
	if err != nil {
		logDebug("[ERROR]", err)
		return
	}
	websocketBroadcast(string(Notification))
}

func notifRenameSuccess(Nickname string) {
	// Message in websocket :
	// action = warn
	// message = nickname-success

	Notification, err := json.Marshal(sendRequest{
		Action:   "warn",
		Message:  "nickname-success",
		Nickname: Nickname,
	})
	if err != nil {
		logDebug("[ERROR]", err)
		return
	}
	websocketBroadcast(string(Notification))
}

func notifRenameFailed(Nickname string) {
	// Message in websocket :
	// action = warn
	// message = nickname-failed

	Notification, err := json.Marshal(sendRequest{
		Action:   "warn",
		Message:  "nickname-failed",
		Nickname: Nickname,
	})
	if err != nil {
		logDebug("[ERROR]", err)
		return
	}
	websocketBroadcast(string(Notification))
}
