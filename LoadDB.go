package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

func LoadHashes() {
	path, _ := filepath.Abs("./saves/hashes.json")
	jsonFile, err := os.Open(path)
	if err != nil {
		return
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &Hashes_Database)
	//Loads hashes.json into Hashes_Database var
}
