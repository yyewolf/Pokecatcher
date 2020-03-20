package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/fatih/color"
)

//SelectedPoke = Struct of the selected pokemon.
type SelectedPoke struct {
	Number int    `json:"number"`
	Name   string `json:"name"`
}

//ConfigStruct = JSON file where the config will be saved
type ConfigStruct struct {
	Token          string `json:"Token"`
	ChannelID      string `json:"Channel_Registered_ID"`
	Delay          int    `json:"Delay_For_Spammer"`
	Duplicate      bool   `json:"Do_I_Catch_Duplicate"`
	AutoCatching   bool   `json:"-"`
	WebPort        int    `json:"Port_For_Webserver"`
	PrefixPokecord string `json:"Pokecord_Prefix_On_Your_Server"`
	PrefixBot      string `json:"Prefix_For_This_Bot"`
	Debug          bool   `json:"Debug"`
	IsAllowedToUse bool   `json:"-"`
}

func NoConfig() {
	//Creates the default config file
	DefaultConfig := ConfigStruct{
		Token:          "Put your token here (in case of problem add 'Bot' before your token)",
		ChannelID:      "Put a channel ID here",
		Delay:          3000,
		Duplicate:      true,
		WebPort:        3000,
		PrefixPokecord: "p!",
		PrefixBot:      "p^",
		Debug:          false,
	}
	path, _ := filepath.Abs("./saves/config.json")
	file, _ := json.MarshalIndent(DefaultConfig, "", " ")
	_ = ioutil.WriteFile(path, file, 0644)
	color.Red("Enter your informations in config.json !")
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

func SaveConfig() {
	path, _ := filepath.Abs("./saves/config.json")
	file, _ := json.MarshalIndent(Config, "", " ")
	_ = ioutil.WriteFile(path, file, 0644)
	//Save Config var into config.json
}
