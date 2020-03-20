package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"fmt"
	"path/filepath"
)

func NoWhitelist() {
	//Creates the default whitelist file, is empty
	DefaultWhitelist := make(map[string]bool)
	path, _ := filepath.Abs("./saves/server_whitelist.json")
	file, _ := json.MarshalIndent(DefaultWhitelist, "", " ")
	_ = ioutil.WriteFile(path, file, 0644)
}

func LoadWhitelist() {
	path, _ := filepath.Abs("./saves/server_whitelist.json")
	jsonFile, err := os.Open(path)
	if err != nil {
		if Config.Debug {
			fmt.Println(err)
		}
		NoWhitelist()
		return
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &ServerWhitelist)
	//Loads server_whitelist.json into Whitelist var
}

func SaveWhitelist() {
	path, _ := filepath.Abs("./saves/server_whitelist.json")
	file, _ := json.MarshalIndent(ServerWhitelist, "", " ")
	_ = ioutil.WriteFile(path, file, 0644)
	//Save Whitelist var into server_whitelist.json
}
