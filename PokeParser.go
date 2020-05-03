package main

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/nfnt/resize"
	"image/png"
)

type PokeInfoParsed struct {
	Name       string
	Level      string
	Number     string
	LastNumber string
	Last       bool
	isInList   bool
	ListNumber string
	IVs        []int
	TotalIV    float64
	isShiny	   bool
}

func ParsePokemonInfo(msg *discordgo.MessageCreate) (PokeInfoParsed, error) {
	Infos := PokeInfoParsed{}
	//Check if there is an embed
	if len(msg.Embeds) == 0 {
		err := errors.New("InfoParser : Message doesn't contain any embeds.")
		return Infos, err
	}
	//Check if there is a title in the embed
	if msg.Embeds[0].Title == "" {
		err := errors.New("InfoParser : Embed doesn't contain any title.")
		return Infos, err
	}
	//Check if there is a title in the embed
	if msg.Embeds[0].Author == nil {
		err := errors.New("InfoParser : Embed doesn't contain any title.")
		return Infos, err
	}
	//Check if there it is a right title
	if msg.Embeds[0].Author.Name != "Professor Oak" {
		err := errors.New("InfoParser : Embed doesn't contain the right Author name. (" + msg.Embeds[0].Author.Name + ")")
		return Infos, err
	}
	//Check if there is a footer in the embed
	if msg.Embeds[0].Footer == nil {
		err := errors.New("InfoParser : Embed doesn't contain any footer.")
		return Infos, err
	}
	//Check if there is an image in the embed
	if msg.Embeds[0].Image == nil {
		err := errors.New("InfoParser : Embed doesn't contain any image.")
		return Infos, err
	}
	
	//Check if the pokemon is a shiny
	if strings.Contains(msg.Embeds[0].Title, "star") {
		Infos.isShiny = true
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
	
	//Gets the name of the pokémon
	ImageURL := msg.Embeds[0].Image.URL
	ImageString := ImageToString(ImageURL)
	Infos.Name = ""
	ImageDecoded, err := loadImg(ImageString)
	if err != nil {
		return Infos, err
	}
	ImageResized := resize.Resize(64, 64, ImageDecoded, resize.Bicubic)
	Buffer := &buf{}
	_ = png.Encode(Buffer, ImageResized)
	ImageResized, _ = png.Decode(Buffer)
	List := box.List()

	for i := range List {
		if strings.Contains(List[i], "img") {
			//Gets rid of the path debris
			Name := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(List[i], "img/", ""), "img\\", ""), ".png", "")
			
			ScanImage := DecodedImages[Name]
			Accuracy := CompareIMG(ScanImage, ImageResized)
			if Accuracy < 0.35 {
				Infos.Name = strings.ReplaceAll(strings.ReplaceAll(Name, "♀", ""), "♂", "")
				break
			}
		}
	}

	//Gets every IV and stores them
	if !strings.Contains(msg.Embeds[0].Description, "IV") {
		err := errors.New("InfoParser : IV not enabled.")
		return Infos, err
	}
	reg, err = regexp.Compile("[^0-9.]")
	if err != nil {
		return Infos, err
	}
	content := strings.Split(msg.Embeds[0].Description, "\n")
	for i := 3; i < len(content)-1; i++ {
		fmt.Println(content[i])
		iv := strings.Split(content[i], "-")[1]
		iv = strings.Split(iv, "/")[0]
		iv = reg.ReplaceAllString(iv, "")
		save, _ := strconv.Atoi(iv)
		Infos.IVs = append(Infos.IVs, save)
	}
	//Gets the total :
	Total := content[len(content)-1]
	Total = reg.ReplaceAllString(Total, "")
	Infos.TotalIV, _ = strconv.ParseFloat(Total, 8)
	
	return Infos, nil
}
