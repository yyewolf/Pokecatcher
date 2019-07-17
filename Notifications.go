package main

import "encoding/json"

func Notif_DuplicateErr(PokemonName string, GuildName string, ChannelName string) {
	// Message in websocket :
	// action = warn
	// message = couldnt-duplicate

	Notification, err := json.Marshal(Send_Request{
		Action:      "warn",
		Message:     "couldnt-duplicate",
		Name:        PokemonName,
		GuildName:   GuildName,
		ChannelName: ChannelName,
	})
	if err != nil {
		return
	}
	Websocket_Broadcast(string(Notification))
}

func Notif_CatchingErr(PokemonName string, GuildName string, ChannelName string) {
	// Message in websocket :
	// action = warn
	// message = couldnt-normal

	Notification, err := json.Marshal(Send_Request{
		Action:      "notification",
		Message:     "couldnt-normal",
		Name:        PokemonName,
		GuildName:   GuildName,
		ChannelName: ChannelName,
	})
	if err != nil {
		return
	}
	Websocket_Broadcast(string(Notification))
}

func Notif_PokeSpawn(PokemonName string, GuildName string, Command string, ChannelName string, ChannelID string) {
	// Message in websocket :
	// action = notification

	Notification, err := json.Marshal(Send_Request{
		Action:      "notification",
		Command:     Command,
		Name:        PokemonName,
		GuildName:   GuildName,
		ChannelName: ChannelName,
		ChannelID:   ChannelID,
	})
	if err != nil {
		return
	}
	Websocket_Broadcast(string(Notification))
}

func Notif_PokeCaught(PokemonName string, GuildName string, ChannelName string) {
	// Message in websocket :
	// action = warn
	// message = could

	Notification, err := json.Marshal(Send_Request{
		Action:      "warn",
		Message:     "could",
		Name:        PokemonName,
		GuildName:   GuildName,
		ChannelName: ChannelName,
	})
	if err != nil {
		return
	}
	Websocket_Broadcast(string(Notification))
}

func Notif_ChannelRegistered() {
	// Message in websocket :
	// action = registered

	Notification, err := json.Marshal(Send_Request{
		Action: "registered",
	})
	if err != nil {
		return
	}
	Websocket_Broadcast(string(Notification))
}

func Notif_RenameSuccess(Nickname string) {
	// Message in websocket :
	// action = warn
	// message = nickname-success

	Notification, err := json.Marshal(Send_Request{
		Action:   "warn",
		Message:  "nickname-success",
		Nickname: Nickname,
	})
	if err != nil {
		return
	}
	Websocket_Broadcast(string(Notification))
}

func Notif_RenameFailed(Nickname string) {
	// Message in websocket :
	// action = warn
	// message = nickname-failed

	Notification, err := json.Marshal(Send_Request{
		Action:   "warn",
		Message:  "nickname-failed",
		Nickname: Nickname,
	})
	if err != nil {
		return
	}
	Websocket_Broadcast(string(Notification))
}
