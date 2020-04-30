package main

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type PokeInfoParsed struct {
	Level      string
	Number     string
	LastNumber string
	Last       bool
	isInList   bool
	ListNumber string
}

func ParsePokemonInfo(msg *discordgo.MessageCreate) (PokeInfoParsed, error) {
	Infos := PokeInfoParsed{}
	//Check if there is an embed
	if len(msg.Embeds) == 0 {
		err := errors.New("Autoleveler : Message doesn't contain any embeds.")
		return Infos, err
	}
	//Check if there is a title in the embed
	if msg.Embeds[0].Title == "" {
		err := errors.New("Autoleveler : Embed doesn't contain any title.")
		return Infos, err
	}
	//Check if there is a title in the embed
	if msg.Embeds[0].Author == nil {
		err := errors.New("Autoleveler : Embed doesn't contain any title.")
		return Infos, err
	}
	//Check if there it is a right title
	if msg.Embeds[0].Author.Name != "Professor Oak" {
		err := errors.New("Autoleveler : Embed doesn't contain the right Author name. (" + msg.Embeds[0].Author.Name + ")")
		return Infos, err
	}
	//Check if there is a footer in the embed
	if msg.Embeds[0].Footer == nil {
		err := errors.New("Autoleveler : Embed doesn't contain any footer.")
		return Infos, err
	}

	//Gets the level from the embed's title : "Level 99 Pokemon"
	//Steps :
	//99
	reg, err := regexp.Compile("[^0-9/]")
	if err != nil {
		return Infos, err
	}
	Infos.Level = reg.ReplaceAllString(msg.Embeds[0].Title, "")

	//Gets the number from the embed's footer : "Selected Pokémon: 50/81 - Use p!back and p!next to cycle through your pokémon!"
	//Steps :
	//50/81
	//50
	Infos.Number = reg.ReplaceAllString(msg.Embeds[0].Footer.Text, "")
	Infos.Number, Infos.LastNumber = strings.Split(Infos.Number, "/")[0], strings.Split(Infos.Number, "/")[1]

	//Stores if the pokemon is the last in the list
	if Infos.Number == Infos.LastNumber {
		Infos.Last = true
	}

	for i := 0; i < Pokemon_List_Info.Realmax; i++ {
		n := strconv.Itoa(i)
		if Infos.Number == Pokemon_List[n].NewNumber {
			Infos.isInList = true
			Infos.ListNumber = n
		}
	}

	if _, ok := Pokemon_List[Infos.Number]; ok {
		Infos.isInList = true
	}

	return Infos, nil
}
