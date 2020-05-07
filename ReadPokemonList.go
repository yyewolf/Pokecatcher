package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

//PokeListInfoStruct = JSON file where the PokemonslistInfos will be saved
type PokeListInfoStruct struct {
	Array   int    `json:"array"`
	Realmax int    `json:"realmax"`
	Names   string `json:"names"`
}

func pokeListForWebsite() string {
	List, _ := json.Marshal(pokemonList)
	if pokemonListInfo.Realmax == 0 {
		return ""
	}
	ListString := string(List)
	ListString = strings.TrimSuffix(ListString, "}")
	ListString += ","
	Infos, _ := json.Marshal(pokemonListInfo)
	InfoString := string(Infos)
	InfoString = strings.TrimPrefix(InfoString, "{")
	ListString += InfoString
	return ListString
}

func noPokemonList() {
	//Creates the default list file
	DefaultList := make(map[string]interface{})
	path, _ := filepath.Abs("./saves/your_pokemon_list.json")
	file, _ := json.MarshalIndent(DefaultList, "", " ")
	_ = ioutil.WriteFile(path, file, 0644)

	//Creates the default list infos file
	DefaultListInfo := PokeListInfoStruct{
		Array:   0,
		Realmax: 0,
		Names:   "",
	}
	path, _ = filepath.Abs("./saves/your_pokemon_list_info.json")
	file, _ = json.MarshalIndent(DefaultListInfo, "", " ")
	_ = ioutil.WriteFile(path, file, 0644)
}

func loadPokemonList() {
	path, _ := filepath.Abs("./saves/your_pokemon_list.json")
	jsonFile, err := os.Open(path)
	if err != nil {
		logDebug("[ERROR] ", err)
		noPokemonList()
		return
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &pokemonList)
	//Loads your_pokemon_list.json into Pokemon_List var

	path, _ = filepath.Abs("./saves/your_pokemon_list_info.json")
	jsonFile, err = os.Open(path)
	if err != nil {
		logDebug("[ERROR] ", err)
		noPokemonList()
		return
	}
	defer jsonFile.Close()
	byteValue, _ = ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &pokemonListInfo)
	//Loads your_pokemon_list_info.json into Pokemon_List_Info var
}

func savePokemonList() {
	path, _ := filepath.Abs("./saves/your_pokemon_list.json")
	file, _ := json.MarshalIndent(pokemonList, "", " ")
	_ = ioutil.WriteFile(path, file, 0644)
	//Save Pokemon_List var into your_pokemon_list.json

	path, _ = filepath.Abs("./saves/your_pokemon_list_info.json")
	file, _ = json.MarshalIndent(pokemonListInfo, "", " ")
	_ = ioutil.WriteFile(path, file, 0644)
	//Save Pokemon_List_Info var into your_pokemon_list_info.json
}
