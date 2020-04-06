package main

import "encoding/json"
import "fmt"

func Websocket_SendMoveList(PokemonName string, PokemonMoves string, LearnChannelID string) {
	// Message in websocket :
	// action = refreshmovelist

	Message, err := json.Marshal(Send_Request{
		Action:         "refreshmovelist",
		PokemonName:    PokemonName,
		PokemonMoves:   PokemonMoves,
		LearnChannelID: LearnChannelID,
	})
	if err != nil {
		if Config.Debug {
			fmt.Println(err)
		}
		return
	}
	Websocket_Broadcast(string(Message))
}

func Websocket_SendPokemonList() {
	// Message in websocket :
	// action = refreshlist

	Message, err := json.Marshal(Send_Request{
		Action:      "refreshlist",
		PokemonList: PokeListForWebsite(),
	})
	if err != nil {
		if Config.Debug {
			fmt.Println(err)
		}
		return
	}
	Websocket_Broadcast(string(Message))
}

func Websocket_RemovedFromList(Number int) {
	// Message in websocket :
	// action = removefromlist

	Message, err := json.Marshal(Send_Request{
		Action:        "removefromlist",
		PokemonList:   PokeListForWebsite(),
		PokemonNumber: Number,
	})
	if err != nil {
		if Config.Debug {
			fmt.Println(err)
		}
		return
	}
	Websocket_Broadcast(string(Message))
}

func Websocket_Selected(PokemonName string, Number int) {
	// Message in websocket :
	// action = selected

	Message, err := json.Marshal(Send_Request{
		Action:        "selected",
		Name:          PokemonName,
		PokemonNumber: Number,
	})
	if err != nil {
		if Config.Debug {
			fmt.Println(err)
		}
		return
	}
	Websocket_Broadcast(string(Message))
}
