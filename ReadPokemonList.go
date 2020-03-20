package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"fmt"
	"path/filepath"
)

func NoPokemonList() {
	//Creates the default list file
	DefaultList := make(map[string]interface{})
	path, _ := filepath.Abs("./saves/your_pokemon_list.json")
	file, _ := json.MarshalIndent(DefaultList, "", " ")
	_ = ioutil.WriteFile(path, file, 0644)
}

func LoadPokemonList() {
	path, _ := filepath.Abs("./saves/your_pokemon_list.json")
	jsonFile, err := os.Open(path)
	if err != nil {
		if Config.Debug {
			fmt.Println(err)
		}
		NoPokemonList()
		return
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &Pokemon_List)
	//Loads your_pokemon_list.json into Pokemon_List var
}

func SavePokemonList() {
	path, _ := filepath.Abs("./saves/your_pokemon_list.json")
	file, _ := json.MarshalIndent(Pokemon_List, "", " ")
	_ = ioutil.WriteFile(path, file, 0644)
	//Save Pokemon_List var into your_pokemon_list.json
}
