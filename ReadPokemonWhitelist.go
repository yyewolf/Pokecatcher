package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func pokeNames() string {
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

func noPokemonWhitelist() {
	//Creates the default list file
	pokemonWhitelist = make(map[string]bool)

	List := box.List()
	for i := range List {
		if strings.Contains(List[i], "img") {
			Name := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(List[i], "img/", ""), "img\\", ""), ".png", "")
			pokemonWhitelist[Name] = true
		}
	}

	path, _ := filepath.Abs("./saves/pokemon_whitelist.json")
	file, _ := json.MarshalIndent(pokemonWhitelist, "", " ")
	_ = ioutil.WriteFile(path, file, 0644)
}

func loadPokemonWhitelist() {
	path, _ := filepath.Abs("./saves/pokemon_whitelist.json")
	jsonFile, err := os.Open(path)
	if err != nil {
		logDebug("[ERROR] ", err)
		noPokemonWhitelist()
		return
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &pokemonWhitelist)
	//Loads pokemon_whitelist.json into Pokemon_Whitelist var
}

func savePokemonWhitelist() {
	path, _ := filepath.Abs("./saves/pokemon_whitelist.json")
	file, _ := json.MarshalIndent(pokemonWhitelist, "", " ")
	_ = ioutil.WriteFile(path, file, 0644)
	//Save Pokemon_Whitelist var into pokemon_whitelist.json
}
