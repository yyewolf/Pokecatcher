package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

//SelectedPoke = Struct of the selected pokemon.
type SelectedPoke struct {
	Number int    `json:"number"`
	Name   string `json:"name"`
}

//ConfigStruct = JSON file where the config will be saved
type ConfigStruct struct {
	Token          string               `json:"Token"`
	ChannelID      string               `json:"Channel_Registered_ID"`
	Delay          int                  `json:"Delay_For_Autocatcher"`
	Duplicate      bool                 `json:"Do_I_Ignore_Duplicate"`
	Aliases        bool                 `json:"Do_I_Catch_With_Aliases"`
	Alerts         bool                 `json:"Do_I_Alert_Someone"`
	AlertChannelID string               `json:"Where_Do_I_Alert_ID"`
	GoodFilter     bool                 `json:"Do_I_Release_Bad_Pokemons"`
	EveryFilters   []customFilterStruct `json:"Custom_Filters"`
	CustomFilters  bool                 `json:"Do_I_Use_Custom_Filters"`
	AutoLevelMax   string               `json:"What_Is_The_Max_Level_AutoLeveler"`
	AutoCatching   bool                 `json:"Do_I_Autocatch_Pokemon"`
	AutoLeveling   bool                 `json:"Do_I_Autolevel_Pokemon"`
	WebPort        int                  `json:"Port_For_Webserver"`
	PrefixPokecord string               `json:"Pokecord_Prefix_On_Your_Server"`
	PrefixBot      string               `json:"Prefix_For_This_Bot"`
	Debug          bool                 `json:"Debug"`
	IsAllowedToUse bool                 `json:"-"`
}

func noConfig() {
	//Creates the default config file
	config = ConfigStruct{
		Token:          "Put your token here (in case of problem add 'Bot' before your token)",
		ChannelID:      "Put a channel ID here",
		Delay:          3000,
		Duplicate:      false,
		Aliases:        true,
		CustomFilters:  false,
		EveryFilters:   []customFilterStruct{},
		WebPort:        3000,
		PrefixPokecord: "p!",
		PrefixBot:      "p^",
		AutoLevelMax:   "100",
		Debug:          false,
		AutoLeveling:	true,
	}
	_, err := os.Stat("saves")
	if os.IsNotExist(err) {
		errDir := os.MkdirAll("saves", 0755)
		if errDir != nil {
			//
		}
	}
	path, _ := filepath.Abs("./saves/config.json")
	file, _ := json.MarshalIndent(config, "", " ")
	_ = ioutil.WriteFile(path, file, 0644)
	logRedLn(logs, "Enter your token in the window!")
}

func loadConfig() {
	path, _ := filepath.Abs("./saves/config.json")
	jsonFile, err := os.Open(path)
	if err != nil {
		noConfig()
		return
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &config)
	//Loads config.json into Config var
}

func noAliases() {
	logRedLn(logs, "Pokemon's aliases are missing, this might result in a blacklist!")
}

func loadAliases() {
	path, _ := filepath.Abs("./saves/aliases.json")
	jsonFile, err := os.Open(path)
	if err != nil {
		noAliases()
		return
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &aliases)
	//Loads aliases.json into Aliases var
}

func saveConfig() {
	path, _ := filepath.Abs("./saves/config.json")
	file, _ := json.MarshalIndent(config, "", " ")
	_ = ioutil.WriteFile(path, file, 0644)
	//Save Config var into config.json
}
