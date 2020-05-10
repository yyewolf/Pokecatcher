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

type prefixes struct {
	Pokecord string `json:"pokecord"`
	Pokebot  string `json:"pokebot"`
}

func hostWebsite() {
	box = packr.NewBox("./www")
	logYellowLn(logs, "Pokemon's decoding started!")
	decodeKnown() // Will decode every resized pokemon images for comparisions.
	logYellowLn(logs, "Pokemon's decoding done !")

	http.HandleFunc("/ws", websocketConnection)
	http.HandleFunc("/", WebsiteHandler)
	http.Handle("/img/", http.FileServer(box))
	//Prevents glitch in case of reconnect

	go openBrowser("http://localhost:" + strconv.Itoa(config.WebPort))
	http.ListenAndServe(":"+strconv.Itoa(config.WebPort), nil)
}

func findServers() (string, string) {
	Names := ""
	IDs := ""
	if discordSession == nil || discordSession.State == nil || discordSession.State.Guilds == nil {
		return Names, IDs
	}
	for i := range discordSession.State.Guilds {
		Names += discordSession.State.Guilds[i].Name + ";"
		IDs += discordSession.State.Guilds[i].ID + ";"
	}
	return Names, IDs
}

//WebsiteHandler : Copies the old handler from server.js, optimize it also.
func WebsiteHandler(w http.ResponseWriter, r *http.Request) {
	//Check if the person is allowed
	if !config.IsAllowedToUse {
		return
	}
	Path := r.URL.Path[1:]
	if Path == "" {
		dat, err := box.FindString("index.html")
		check(err)
		Filters, _ := json.Marshal(config.EveryFilters)
		if string(Filters) == "null" {
			Filters = []byte("[]")
		}
		Legend, _ := json.Marshal(legendaries)
		Queue := strings.Join(priorityQueue, ";")
		Whitelist, _ := json.Marshal(serverWhitelist)
		Selected, _ := json.Marshal(selectedPokemon)
		PokeWhitelist, _ := json.Marshal(pokemonWhitelist)
		ServerNames, ServerIDs := findServers()
		Prefix := prefixes{
			Pokebot:  config.PrefixBot,
			Pokecord: config.PrefixPokecord,
		}
		SendPrefix, _ := json.Marshal(Prefix)
		add := "<script> var websocket = '" + strconv.Itoa(config.WebPort) + "/ws';</script>\n"
		add += "<script> var autocatchdelay = '" + strconv.Itoa(config.Delay) + "';</script>\n"
		add += "<script> var customfilters = " + strconv.FormatBool(config.CustomFilters) + ";</script>\n"
		add += "<script> var whitelist = '" + string(Whitelist) + "';</script>\n"
		add += "<script> var legendaries = " + string(Legend) + ";</script>\n"
		add += "<script> var serverid = '" + ServerIDs + "';</script>\n"
		add += "<script> var servernames = '" + strings.ReplaceAll(ServerNames, "'", "\\'") + "';</script>\n"
		add += "<script> var spamactive = " + strconv.FormatBool(spamState) + ";</script>\n"
		add += "<script> var token = '" + config.Token + "';</script>\n"
		add += "<script> var selected = " + string(Selected) + ";</script>\n"
		add += "<script> var textchannel = '" + config.ChannelID + "';</script>\n"
		add += "<script> var autocatcher = " + strconv.FormatBool(config.AutoCatching) + "</script>\n"
		add += "<script> var duplicate = " + strconv.FormatBool(config.Duplicate) + "</script>\n"
		add += "<script> var aliases = " + strconv.FormatBool(config.Aliases) + "</script>\n"
		add += "<script> var filter = " + strconv.FormatBool(config.GoodFilter) + "</script>\n"
		add += "<script> var prefixes = " + string(SendPrefix) + "</script>\n"
		add += "<script> var pokewhitelist = " + string(PokeWhitelist) + "</script>\n"
		add += "<script> var queue = '" + Queue + "'</script>\n"
		add += "<script> var filters = " + string(Filters) + "</script>\n"
		add += "<script> var autolevelmax = '" + config.AutoLevelMax + "'</script>\n"

		MaxPoke := strconv.Itoa(pokemonListInfo.Array)
		if MaxPoke != "0" {
			add += "<script> var listobj = " + pokeListForWebsite() + "</script>\n"
		}
		fmt.Fprint(w, add+dat)
	} else {
		dat, err := box.FindString(Path)
		check(err)
		fmt.Fprint(w, dat)
	}
}
