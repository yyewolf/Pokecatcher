package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gobuffalo/packr"
)

var box packr.Box

type Prefixes struct {
	Pokecord string `json:"pokecord"`
	Pokebot  string `json:"pokebot"`
}

func Host_Website() {
	box = packr.NewBox("./www")
	LogYellowLn(Logs, "Pokemon's decoding started!")
	DecodeKnown() // Will decode every resized pokemon images for comparisions.
	LogYellowLn(Logs, "Pokemon's decoding done !")

	http.HandleFunc("/ws", Websocket_Connection)
	http.HandleFunc("/", Website_Handler)
	http.Handle("/img/", http.FileServer(box))
	//Prevents glitch in case of reconnect

	go OpenBrowser("http://localhost:" + strconv.Itoa(Config.WebPort))
	http.ListenAndServe(":"+strconv.Itoa(Config.WebPort), nil)
}

func FindServers() (string, string) {
	Names := ""
	IDs := ""
	if DiscordSession == nil || DiscordSession.State == nil || DiscordSession.State.Guilds == nil {
		return Names, IDs
	}
	for i := range DiscordSession.State.Guilds {
		Names += DiscordSession.State.Guilds[i].Name + ";"
		IDs += DiscordSession.State.Guilds[i].ID + ";"
	}
	return Names, IDs
}

//Website_Handler : Copies the old handler from server.js, optimize it also.
func Website_Handler(w http.ResponseWriter, r *http.Request) {
	//Check if the person is allowed
	if !Config.IsAllowedToUse {
		return
	}
	Path := r.URL.Path[1:]
	if Path == "" {
		dat, err := box.FindString("index.html")
		check(err)
		Legend, _ := json.Marshal(Legendaries)
		Whitelist, _ := json.Marshal(ServerWhitelist)
		Selected, _ := json.Marshal(SelectedPokemon)
		PokeWhitelist, _ := json.Marshal(Pokemon_Whitelist)
		ServerNames, ServerIDs := FindServers()
		Prefix := Prefixes{
			Pokebot:  Config.PrefixBot,
			Pokecord: Config.PrefixPokecord,
		}
		SendPrefix, _ := json.Marshal(Prefix)
		add := "<script> var websocket = '" + strconv.Itoa(Config.WebPort) + "/ws';</script>\n"
		add += "<script> var autocatchdelay = '" + strconv.Itoa(Config.Delay) + "';</script>\n"
		add += "<script> var whitelist = '" + string(Whitelist) + "';</script>\n"
		add += "<script> var legendaries = " + string(Legend) + ";</script>\n"
		add += "<script> var serverid = '" + ServerIDs + "';</script>\n"
		add += "<script> var servernames = '" + strings.ReplaceAll(ServerNames, "'", "\\'") + "';</script>\n"
		add += "<script> var spamactive = " + strconv.FormatBool(SpamState) + ";</script>\n"
		add += "<script> var token = '" + Config.Token + "';</script>\n"
		add += "<script> var selected = " + string(Selected) + ";</script>\n"
		add += "<script> var textchannel = '" + Config.ChannelID + "';</script>\n"
		add += "<script> var autocatcher = " + strconv.FormatBool(Config.AutoCatching) + "</script>\n"
		add += "<script> var duplicate = " + strconv.FormatBool(Config.Duplicate) + "</script>\n"
		add += "<script> var aliases = " + strconv.FormatBool(Config.Aliases) + "</script>\n"
		add += "<script> var filter = " + strconv.FormatBool(Config.GoodFilter) + "</script>\n"
		add += "<script> var prefixes = " + string(SendPrefix) + "</script>\n"
		add += "<script> var pokewhitelist = " + string(PokeWhitelist) + "</script>\n"

		MaxPoke := strconv.Itoa(Pokemon_List_Info.Array)
		if MaxPoke != "0" {
			add += "<script> var listobj = " + PokeListForWebsite() + "</script>\n"
		}
		fmt.Fprint(w, add+dat)
	} else {
		dat, err := box.FindString(Path)
		check(err)
		fmt.Fprint(w, dat)
	}
}
