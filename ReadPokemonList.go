package main

import (
	"encoding/json"
	"fmt"
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

func PokeListForWebsite() string {
	List, _ := json.Marshal(Pokemon_List)
	ListString := string(List)
	ListString = strings.TrimSuffix(ListString, "}")
	ListString += ","
	Infos, _ := json.Marshal(Pokemon_List_Info)
	InfoString := string(Infos)
	InfoString = strings.TrimPrefix(InfoString, "{")
	ListString += InfoString
	return ListString
}

func NoPokemonList() {
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

	path, _ = filepath.Abs("./saves/your_pokemon_list_info.json")
	jsonFile, err = os.Open(path)
	if err != nil {
		if Config.Debug {
			fmt.Println(err)
		}
		NoPokemonList()
		return
	}
	defer jsonFile.Close()
	byteValue, _ = ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &Pokemon_List_Info)
	//Loads your_pokemon_list_info.json into Pokemon_List_Info var
}

func SavePokemonList() {
	path, _ := filepath.Abs("./saves/your_pokemon_list.json")
	file, _ := json.MarshalIndent(Pokemon_List, "", " ")
	_ = ioutil.WriteFile(path, file, 0644)
	//Save Pokemon_List var into your_pokemon_list.json

	path, _ = filepath.Abs("./saves/your_pokemon_list_info.json")
	file, _ = json.MarshalIndent(Pokemon_List_Info, "", " ")
	_ = ioutil.WriteFile(path, file, 0644)
	//Save Pokemon_List_Info var into your_pokemon_list_info.json
}
