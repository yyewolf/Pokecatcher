package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

//IsAGoodPokemon will verify the pokémon using the default filter
func IsAGoodPokemon(Pokemons pokeInfoParsed) bool {
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
	for i := range legendaries {
		if Pokemons.Name == legendaries[i] {
			return true
		}
	}
	return false
}

func filterPokemon(Pkmn pokeInfoParsed) bool {
	r := true
	if config.GoodFilter {
		r = IsAGoodPokemon(Pkmn)
	}
	//This is a logic OR
	if config.CustomFilters {
		r = allFilters(Pkmn)
	}
	logDebug("[DEBUG] Finished checking r = ", r)
	return r
}

//AutoRelease will verify the infos of the current pokemon and release it if it's bad
func AutoRelease(s *discordgo.Session, msg *discordgo.MessageCreate) {
	//Check if the person is allowed
	if !config.IsAllowedToUse {
		return
	}
	//Check if the author is pokecord
	if msg.Author.ID != "365975655608745985" {
		return
	}
	//Check that the function is activated
	if !infoMenu.Activated {
		return
	}
	//Check that the channel is the right one
	if msg.ChannelID != infoMenu.ChannelID {
		return
	}
	//Verify that it's an answer to the user
	msgs, err := s.ChannelMessages(msg.ChannelID, 5, msg.ID, "", "")
	if err != nil {
		logDebug("[ERROR]", err)
		return
	}
	ok := false
	for i := range msgs {
		if msgs[i].ID == infoMenu.MessageID {
			ok = true
			break
		}
	}
	if !ok {
		return
	}
	logDebug("[Debug] Entered Auto Releaser.")
	//Check if it's a p!info response
	Infos, err := parsePokemonInfo(msg)
	if err != nil {
		logDebug("[ERROR]", err)
		return
	}

	//Sends it to website
	if Infos.isInList {
		Current := pokemonList[Infos.ListNumber]
		Current.Level = Infos.Level
		Current.IV = "IV: " + fmt.Sprintf("%.2f", Infos.TotalIV)
		pokemonList[Infos.ListNumber] = Current
		savePokemonList()
		websocketSendPokemonList()
	}

	//Verify that it should release it
	if !infoMenu.AutoRelease {
		return
	}

	infoMenu.Activated = false
	infoMenu.AutoRelease = false

	//Won't release it if it's a good pokémon (decided by every filter)
	if filterPokemon(Infos) {
		logCyanLn(logs, "You caught a good Pokémon !")
		return
	}

	logDebug("[DEBUG] Will release a ", Infos.Name)

	time.Sleep(3 * time.Second)

	//Release the pokémon and removes it from list
	n, _ := strconv.Atoi(Infos.Number)
	release(receiveRequest{
		PokemonNumber: n,
	})
}
