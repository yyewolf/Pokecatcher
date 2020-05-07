package main

import (
	"net/http"
)

type sendRequest struct {
	Action         string `json:"action"`
	Command        string `json:"command"`
	Name           string `json:"name"`
	Nickname       string `json:"nickname"`
	GuildName      string `json:"server"`
	ChannelName    string `json:"channelname"`
	ChannelID      string `json:"channelid"`
	Message        string `json:"message"`
	PokemonList    string `json:"listobj"`
	PokemonNumber  int    `json:"number"`
	PokemonName    string `json:"pokemonname"`
	PokemonMoves   string `json:"moves"`
	LearnChannelID string `json:"channelmovesid"`
}

type receiveRequest struct {
	Action         string `json:"action"`
	State          bool   `json:"state"`
	Nickname       string `json:"nickname"`
	Name           string `json:"name"`
	GuildName      string `json:"server"`
	Token          string `json:"token"`
	SpamInterval   int    `json:"interval"`
	Message        string `json:"message"`
	Command        string `json:"command"`
	ChannelID      string `json:"channelid"`
	PokemonNumber  int    `json:"number"`
	AutoCatchDelay int    `json:"delay"`

	//When learning new move.
	PokemonName  string `json:"pokemonname"`
	MovePosition int    `json:"movenumber"`
	MoveName     string `json:"movename"`
	ChannelSet   string `json:"channelset"`

	//When changing prefixes.
	Type   string `json:"type"`
	Prefix string `json:"prefix"`

	//When changing priority queue.
	Change string `json:"change"`

	//When adding/removing from Whitelist
	GuildID    string `json:"serverid"`
	GuildState bool   `json:"serverstate"`

	//When updating custom filters
	Filters []customFilterStruct `json:"filters"`
}

func websocketBroadcast(msg string) {
	for i := range connections {
		if i < len(connections) {
			err := connections[i].WriteMessage(1, []byte(msg))
			if err != nil {
				logDebug("[ERROR]", err)
				connections = append(connections[:i], connections[i+1:]...)
				//Removes the connection if it's unable to send a message.
			}
		}
	}
}

func websocketConnection(w http.ResponseWriter, r *http.Request) {
	//Check if the person is allowed
	if !config.IsAllowedToUse {
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity
	if err != nil {
		logDebug("[ERROR]", err)
		return
	}
	conn.SetCloseHandler(func(code int, text string) error {
		for i := range connections {
			if connections[i] == conn {
				connections = append(connections[:i], connections[i+1:]...)
			}
		}
		return nil
	})
	connections = append(connections, conn)

	for {
		// Read message from browser
		rec := receiveRequest{}
		err := conn.ReadJSON(&rec)
		if err != nil {
			logDebug("[ERROR]", err)
			return
		}
		if _, ok := websocketReceiveFunctions[rec.Action]; !ok {
			return
		}
		websocketReceiveFunctions[rec.Action](rec)
	}
}
