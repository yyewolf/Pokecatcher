package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func PokeNames() string {
	List := box.List()
	Names := []string{}
	for i := range List {
		if strings.Contains(List[i], "img") {
			Name := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(List[i], "img/", ""), "img\\", ""), ".png", "")
			Names = append(Names, Name)
		}
	}
	s, _ := json.Marshal(Names)
	return string(s)
}

func NoPokemonWhitelist() {
	//Creates the default list file
	Pokemon_Whitelist = make(map[string]bool)
	
	List := box.List()
	for i := range List {
		if strings.Contains(List[i], "img") {
			Name := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(List[i], "img/", ""), "img\\", ""), ".png", "")
			Pokemon_Whitelist[Name] = true
		}
	}
	
	path, _ := filepath.Abs("./saves/pokemon_whitelist.json")
	file, _ := json.MarshalIndent(Pokemon_Whitelist, "", " ")
	_ = ioutil.WriteFile(path, file, 0644)
}

func LoadPokemonWhitelist() {
	path, _ := filepath.Abs("./saves/pokemon_whitelist.json")
	jsonFile, err := os.Open(path)
	if err != nil {
		Debug("[ERROR] ", err)
		NoPokemonWhitelist()
		return
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &Pokemon_Whitelist)
	//Loads pokemon_whitelist.json into Pokemon_Whitelist var
}

func SavePokemonWhitelist() {
	path, _ := filepath.Abs("./saves/pokemon_whitelist.json")
	file, _ := json.MarshalIndent(Pokemon_Whitelist, "", " ")
	_ = ioutil.WriteFile(path, file, 0644)
	//Save Pokemon_Whitelist var into pokemon_whitelist.json
}
