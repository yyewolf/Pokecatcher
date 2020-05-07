package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

func noWhitelist() {
	//Creates the default whitelist file, is empty
	DefaultWhitelist := make(map[string]bool)
	path, _ := filepath.Abs("./saves/server_whitelist.json")
	file, _ := json.MarshalIndent(DefaultWhitelist, "", " ")
	_ = ioutil.WriteFile(path, file, 0644)
}

func loadWhitelist() {
	path, _ := filepath.Abs("./saves/server_whitelist.json")
	jsonFile, err := os.Open(path)
	if err != nil {
		logDebug("[ERROR] ", err)
		noWhitelist()
		return
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &serverWhitelist)
	//Loads server_whitelist.json into Whitelist var
}

func saveWhitelist() {
	path, _ := filepath.Abs("./saves/server_whitelist.json")
	file, _ := json.MarshalIndent(serverWhitelist, "", " ")
	_ = ioutil.WriteFile(path, file, 0644)
	//Save Whitelist var into server_whitelist.json
}
