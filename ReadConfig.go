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
	GoodFilter     bool                 `json:"Do_I_Release_Bad_Pokemons"`
	EveryFilters   []customFilterStruct `json:"Custom_Filters"`
	CustomFilters  bool                 `json:"Do_I_Use_Custom_Filters"`
	AutoCatching   bool                 `json:"-"`
	WebPort        int                  `json:"Port_For_Webserver"`
	PrefixPokecord string               `json:"Pokecord_Prefix_On_Your_Server"`
	PrefixBot      string               `json:"Prefix_For_This_Bot"`
	Debug          bool                 `json:"Debug"`
	IsAllowedToUse bool                 `json:"-"`
}

func NoConfig() {
	//Creates the default config file
	Config = ConfigStruct{
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
		Debug:          false,
	}
	_, err := os.Stat("saves")
	if os.IsNotExist(err) {
		errDir := os.MkdirAll("saves", 0755)
		if errDir != nil {
			//
		}
	}
	path, _ := filepath.Abs("./saves/config.json")
	file, _ := json.MarshalIndent(Config, "", " ")
	_ = ioutil.WriteFile(path, file, 0644)
	LogRedLn(Logs, "Enter your token in the window!")
}

func LoadConfig() {
	path, _ := filepath.Abs("./saves/config.json")
	jsonFile, err := os.Open(path)
	if err != nil {
		NoConfig()
		return
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &Config)
	//Loads config.json into Config var
}

func NoAliases() {
	LogRedLn(Logs, "Pokemon's aliases are missing, this might result in a blacklist!")
}

func LoadAliases() {
	path, _ := filepath.Abs("./saves/aliases.json")
	jsonFile, err := os.Open(path)
	if err != nil {
		NoAliases()
		return
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &Aliases)
	//Loads aliases.json into Aliases var
}

func SaveConfig() {
	path, _ := filepath.Abs("./saves/config.json")
	file, _ := json.MarshalIndent(Config, "", " ")
	_ = ioutil.WriteFile(path, file, 0644)
	//Save Config var into config.json
}
