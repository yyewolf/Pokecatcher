package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

//IsAGoodPokemon will verify the pokémon using the default filter
func IsAGoodPokemon(Pokemons PokeInfoParsed) bool {
	if len(Pokemons.IVs) == 0 {
		return false
	}
	//Shiny = good
	if Pokemons.isShiny {
		return true
	}
	c1 := 0 // IV = 31/31
	c2 := 0 // IV = 0/31
	for i := range Pokemons.IVs {
		c := Pokemons.IVs[i]
		if c == 31 {
			c1++
		} else if c == 0 {
			c2++
		}
	}
	//triple max IV stats, triple 0 IV stats
	if c1 >= 3 || c2 >= 3 {
		return true
	}
	//IV below 10%
	if Pokemons.TotalIV < 10 {
		return true
	}
	//IV above 80%
	if Pokemons.TotalIV > 80 {
		return true
	}
	//Pokemon is legendary
	for i := range Legendaries {
		if Pokemons.Name == Legendaries[i] {
			return true
		}
	}
	return false
}

func filterPokemon(Pkmn PokeInfoParsed) bool {
	r := true
	if Config.GoodFilter {
		r = IsAGoodPokemon(Pkmn)
	}
	if Config.CustomFilters {
		r = allFilters(Pkmn)
	}
	Debug("[DEBUG] Finished checking r = ", r)
	return r
}

//AutoRelease will verify the infos of the current pokemon and release it if it's bad
func AutoRelease(s *discordgo.Session, msg *discordgo.MessageCreate) {
	//Check if the person is allowed
	if !Config.IsAllowedToUse {
		return
	}
	//Check if the author is pokecord
	if msg.Author.ID != "365975655608745985" {
		return
	}
	if !InfoMenu.Activated {
		return
	}
	if msg.ChannelID != InfoMenu.ChannelID {
		return
	}
	msgs, err := s.ChannelMessages(msg.ChannelID, 5, msg.ID, "", "")
	if err != nil {
		Debug("[ERROR]", err)
		return
	}
	ok := false
	for i := range msgs {
		if msgs[i].ID == InfoMenu.MessageID {
			ok = true
			break
		}
	}
	if !ok {
		return
	}
	Debug("[Debug] Entered Auto Releaser.")
	//Check if it's a p!info response
	Infos, err := ParsePokemonInfo(msg)
	if err != nil {
		Debug("[ERROR]", err)
		return
	}

	if Infos.isInList {
		Current := Pokemon_List[Infos.ListNumber]
		Current.Level = Infos.Level
		Current.IV = "IV: " + fmt.Sprintf("%.2f", Infos.TotalIV)
		Pokemon_List[Infos.ListNumber] = Current
		SavePokemonList()
	}

	if !InfoMenu.AutoRelease {
		return
	}

	InfoMenu.Activated = false
	InfoMenu.AutoRelease = false

	//Won't release it if it's a good pokémon (decided by every filter)
	if filterPokemon(Infos) {
		LogCyanLn(Logs, "You caught a good Pokémon !")
		return
	}

	Debug("[DEBUG] Will release a ", Infos.Name)

	time.Sleep(3 * time.Second)

	//Release the pokémon and removes it from list
	n, _ := strconv.Atoi(Infos.Number)
	Release(Receive_Request{
		PokemonNumber: n,
	})
}
