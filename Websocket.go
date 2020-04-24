package main

import (
	"fmt"
	"net/http"
)

type Send_Request struct {
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

type Receive_Request struct {
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

	//When changing farmer settings.
	Change   string `json:"change"`
	Settings string `json:"setting"`

	//When adding/removing from Whitelist
	GuildID    string `json:"serverid"`
	GuildState bool   `json:"serverstate"`
}

func Websocket_Broadcast(msg string) {
	for i := range Connections {
		if i < len(Connections) {
			err := Connections[i].WriteMessage(1, []byte(msg))
			if err != nil {
				if Config.Debug {
					fmt.Println(err)
				}
				Connections = append(Connections[:i], Connections[i+1:]...)
				//Removes the connection if it's unable to send a message.
			}
		}
	}
}

func Websocket_Connection(w http.ResponseWriter, r *http.Request) {
	//Check if the person is allowed
	if !Config.IsAllowedToUse {
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity
	if err != nil {
		if Config.Debug {
			fmt.Println(err)
		}
		return
	}
	conn.SetCloseHandler(func(code int, text string) error {
		for i := range Connections {
			if Connections[i] == conn {
				Connections = append(Connections[:i], Connections[i+1:]...)
			}
		}
		return nil
	})
	Connections = append(Connections, conn)

	for {
		// Read message from browser
		rec := Receive_Request{}
		err := conn.ReadJSON(&rec)
		if err != nil {
			if Config.Debug {
				fmt.Println(err)
			}
			return
		}
		if _, ok := Websocket_Receive_Functions[rec.Action]; !ok {
			return
		}
		Websocket_Receive_Functions[rec.Action](rec)
	}
}
