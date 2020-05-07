package main

import (
	"encoding/json"
	"fmt"
)

func websocketSendMoveList(PokemonName string, PokemonMoves string, LearnChannelID string) {
	// Message in websocket :
	// action = refreshmovelist

	Message, err := json.Marshal(sendRequest{
		Action:         "refreshmovelist",
		PokemonName:    PokemonName,
		PokemonMoves:   PokemonMoves,
		LearnChannelID: LearnChannelID,
	})
	if err != nil {
		if config.Debug {
			fmt.Print("[ERROR] ")
			fmt.Println(err)
		}
		return
	}
	websocketBroadcast(string(Message))
}

func websocketSendPokemonList() {
	// Message in websocket :
	// action = refreshlist

	Message, err := json.Marshal(sendRequest{
		Action:      "refreshlist",
		PokemonList: pokeListForWebsite(),
	})
	if err != nil {
		if config.Debug {
			fmt.Print("[ERROR] ")
			fmt.Println(err)
		}
		return
	}
	websocketBroadcast(string(Message))
}

func websocketRemovedFromList(Number int) {
	// Message in websocket :
	// action = removefromlist

	Message, err := json.Marshal(sendRequest{
		Action:        "removefromlist",
		PokemonList:   pokeListForWebsite(),
		PokemonNumber: Number,
	})
	if err != nil {
		if config.Debug {
			fmt.Print("[ERROR] ")
			fmt.Println(err)
		}
		return
	}
	websocketBroadcast(string(Message))
}

func websocketSelected(PokemonName string, Number int) {
	// Message in websocket :
	// action = selected

	Message, err := json.Marshal(sendRequest{
		Action:        "selected",
		Name:          PokemonName,
		PokemonNumber: Number,
	})
	if err != nil {
		if config.Debug {
			fmt.Print("[ERROR] ")
			fmt.Println(err)
		}
		return
	}
	websocketBroadcast(string(Message))
}
