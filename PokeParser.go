package main

import (
	"errors"
	"image/png"
	"regexp"
	"strconv"
	"strings"
	"image"

	"github.com/bwmarrin/discordgo"
	"github.com/nfnt/resize"
)

type pokeInfoParsed struct {
	Name       string
	Image	   image.Image
	Level      string
	Number     string
	LastNumber string
	Last       bool
	isInList   bool
	ListNumber string
	IVs        []int
	TotalIV    float64
	isShiny    bool
	isGalarian bool
	isAlolan   bool
}

func parsePokemonInfo(msg *discordgo.MessageCreate) (pokeInfoParsed, error) {
	Infos := pokeInfoParsed{}
	//Check if there is an embed
	if len(msg.Embeds) == 0 {
		err := errors.New("infoparser : message doesn't contain any embeds")
		return Infos, err
	}
	//Check if there is a title in the embed
	if msg.Embeds[0].Title == "" {
		err := errors.New("infoparser : embed doesn't contain any title")
		return Infos, err
	}
	//Check if there is a title in the embed
	if msg.Embeds[0].Author == nil {
		err := errors.New("infoparser : embed doesn't contain any title")
		return Infos, err
	}
	//Check if there it is a right title
	if msg.Embeds[0].Author.Name != "Professor Oak" {
		err := errors.New("infoparser : embed doesn't contain the right author name (" + msg.Embeds[0].Author.Name + ")")
		return Infos, err
	}
	//Check if there is a footer in the embed
	if msg.Embeds[0].Footer == nil {
		err := errors.New("infoparser : embed doesn't contain any footer")
		return Infos, err
	}
	//Check if there is an image in the embed
	if msg.Embeds[0].Image == nil {
		err := errors.New("infoparser : embed doesn't contain any image")
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

	for i := 0; i < pokemonListInfo.Realmax; i++ {
		n := strconv.Itoa(i)
		if Infos.Number == pokemonList[n].NewNumber {
			Infos.isInList = true
			Infos.ListNumber = n
		}
	}

	//Gets the name of the pokémon
	ImageURL := msg.Embeds[0].Image.URL
	ImageString := imageToString(ImageURL)
	Infos.Name = ""
	ImageDecoded, err := loadImg(ImageString)
	if err != nil {
		return Infos, err
	}
	ImageResized := resize.Resize(64, 64, ImageDecoded, resize.Bicubic)
	Buffer := &buf{}
	err = png.Encode(Buffer, ImageResized)
	if err != nil {
		return Infos, err
	}
	ImageResized, err = png.Decode(Buffer)
	if err != nil {
		return Infos, err
	}
	List := box.List()
	
	logDebug("1")
	for i := range List {
		if strings.Contains(List[i], "img") {
			//Gets rid of the path debris
			Name := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(List[i], "img/", ""), "img\\", ""), ".png", "")

			ScanImage := decodedImages[Name]
			Infos.Image = ScanImage
			Accuracy := compareIMG(ScanImage, ImageResized)
			if Accuracy < 0.35 {
				Infos.Name = strings.ReplaceAll(strings.ReplaceAll(Name, "♀", ""), "♂", "")
				break
			}
		}
	}
	logDebug("2")
	//Gets every IV and stores them
	if !strings.Contains(msg.Embeds[0].Description, "IV") {
		err := errors.New("infoparser : iv not enabled")
		return Infos, err
	}
	reg, err = regexp.Compile("[^0-9.]")
	if err != nil {
		return Infos, err
	}
	content := strings.Split(msg.Embeds[0].Description, "\n")
	for i := 3; i < len(content)-1; i++ {
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

	//Sets the galarian and alona field
	if strings.Contains(Infos.Name, "Alolan") {
		Infos.isAlolan = true
	}
	if strings.Contains(Infos.Name, "Galarian") {
		Infos.isAlolan = true
	}
	return Infos, nil
}
