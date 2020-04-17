package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/fatih/color"
	"github.com/nfnt/resize"
)

func ImageToString(URL string) string {
	response, e := http.Get(URL)
	//response => image
	check(e)
	defer response.Body.Close()
	//Closes the web page when it's done
	image, err := ioutil.ReadAll(response.Body)
	if err != nil {
		if Config.Debug {
			fmt.Println(err)
		}
		return "nothing"
	}
	return string(image)
}

func LogPokemonSpawn(PokemonName string, GuildName string, ChannelName string, Accuracy float64, AliasUsed string) {
	wgPokeSpawn.Wait()
	wgPokeSpawn.Add(1)
	color.Green("-------------------------------------------------------")
	PrintGreen := color.New(color.FgHiGreen).PrintfFunc()
	PrintBlue := color.New(color.FgHiBlue).PrintfFunc()
	PrintGreen("A ")
	PrintBlue(PokemonName)
	PrintGreen(" has spawned on : \nGuild Name : "+GuildName+"\nChannel Name : #"+ChannelName+".\nAccuracy : %f%%\nAlias used : "+AliasUsed+"\n", Accuracy)
	wgPokeSpawn.Done()
}

func FakeTalk(s *discordgo.Session, ChannelID string, Letters int) {
	//Fakes user typing
	for start := time.Now(); time.Since(start) < time.Duration(Config.Delay)*time.Millisecond; {
		_ = s.ChannelTyping(ChannelID)
		time.Sleep(time.Duration(Config.Delay/Letters) * time.Millisecond)
	}
}

func HasAliases(Pokemon string) bool {
	if _, ok := Aliases[Pokemon]; ok {
		return true
	}
	return false
}

func CheckForPokemon(s *discordgo.Session, msg *discordgo.MessageCreate) {
	//Check if the person is allowed
	if !Config.IsAllowedToUse {
		return
	}
	//Check if the server is whitelisted
	if !ServerWhitelist[msg.GuildID] {
		return
	}
	//Check if the author is pokecord
	if msg.Author.ID != "365975655608745985" {
		return
	}
	DiscordSession = s
	//Check if there is an embed
	if len(msg.Embeds) == 0 {
		return
	}
	if msg.Embeds[0].Image == nil {
		return
	}
	//Check if it's a pokemon spawn
	if !strings.Contains(msg.Embeds[0].Title, "A wild") {
		return
	}
	ImageURL := msg.Embeds[0].Image.URL
	ImageString := ImageToString(ImageURL)
	Spawned_Pokemon_Name := ""
	ImageDecoded, err := loadImg(ImageString)
	if err != nil {
		if Config.Debug {
			fmt.Println(err)
		}
		return
	}
	ImageResized := resize.Resize(64, 64, ImageDecoded, resize.Bicubic)
	List := box.List()
	Accuracy := 0.0

	for i := range List {
		if strings.Contains(List[i], "img") {
			Name := strings.ReplaceAll(strings.ReplaceAll(List[i], "img/", ""), ".png", "")
			ScanImage := DecodedImages[Name]
			Accuracy = CompareIMG(ScanImage, ImageResized)
			if Accuracy < 0.35 {
				Spawned_Pokemon_Name = strings.ReplaceAll(strings.ReplaceAll(Name, "♀", ""), "♂", "")
				break
			}
		}
	}

	Accuracy = 100.0 - Accuracy

	if Spawned_Pokemon_Name == "" {
		return
	}
	Guild_Spawn, err := s.Guild(msg.GuildID)
	if err != nil {
		if Config.Debug {
			fmt.Println(err)
		}
		return
	}
	Channel_Spawn, err := s.Channel(msg.ChannelID)
	if err != nil {
		if Config.Debug {
			fmt.Println(err)
		}
		return
	}
	//Logs info into the console and sends a notification to the website.
	OriginalName := Spawned_Pokemon_Name
	CatchName := Spawned_Pokemon_Name
	if HasAliases(Spawned_Pokemon_Name) {
		Names := Aliases[OriginalName]
		CatchName = Names[0]
		if Config.Aliases {
			if len(Names) == 1 {
				CatchName = Names[0]
			} else {
				CatchName = Names[rand.Intn(len(Names)-1)]
			}
		}
	}

	LogPokemonSpawn(OriginalName, Guild_Spawn.Name, Channel_Spawn.Name, Accuracy, CatchName)
	//Gets the command from the message : "Guess the pokemon and type p!catch <pokémon> to catch it !"
	Command_To_Catch := strings.Split(strings.Split(msg.Embeds[0].Description, "type ")[1], " <po")[0]
	Notif_PokeSpawn(OriginalName, Guild_Spawn.Name, Command_To_Catch, Channel_Spawn.Name, Channel_Spawn.ID)
	color.HiBlue("Command : " + Command_To_Catch + " " + OriginalName)
	if Config.AutoCatching {
		//Closes spammer
		if SpamState {
			SpamChannel <- 1
		}
		FakeTalk(s, msg.ChannelID, len(Command_To_Catch+" "+strings.ToLower(CatchName)))

		rand.Seed(time.Now().UnixNano())
		RandomNess := rand.Intn(250) - rand.Intn(250)

		time.Sleep(time.Duration(Config.Delay+RandomNess) * time.Millisecond)
		color.Blue("Tried to catch your : " + OriginalName)

		Command_To_Catch = strings.ReplaceAll(Command_To_Catch, "а", "a")
		_, err := s.ChannelMessageSend(msg.ChannelID, Command_To_Catch+" "+strings.ToLower(CatchName))
		if err != nil {
			if Config.Debug {
				fmt.Println(err)
			}
			Notif_CatchingErr(OriginalName, Guild_Spawn.Name, Channel_Spawn.Name)
			return
		}
		//Restart spammer
		if SpamState {
			go SpamFunc(DiscordSession, Config.ChannelID, SpamInterval, SpamMessage)
		}
	}
}

func SuccessfulCatch(s *discordgo.Session, msg *discordgo.MessageCreate) {
	//Check if the person is allowed
	if !Config.IsAllowedToUse {
		return
	}
	//Check if the author is pokecord
	if msg.Author.ID != "365975655608745985" {
		return
	}
	//Check if it's a pokemon catch
	if !strings.Contains(msg.Content, "Congratulations") {
		return
	}
	//Check if the author is mentioned
	hasMention := false
	for i := range msg.Mentions {
		if msg.Mentions[i].ID == s.State.User.ID {
			hasMention = true
			break
		}
	}
	if !hasMention {
		return
	}
	reg, err := regexp.Compile("[0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	//Gets the command from the message : "Congratulations @User! You caught a level 99 Pokemon! Added to Pokédex."
	//Steps :
	//99 Pokemon! Added to Pokédex."
	//99 Pokemon
	//99Pokemon
	//Pokemon
	PokemonName := reg.ReplaceAllString(strings.ReplaceAll(strings.Split(strings.Split(msg.Content, "level ")[1], "! Added")[0], " ", ""), "")
	PokemonLevel := strings.Split(strings.Split(msg.Content, "level ")[1], " "+PokemonName)[0]

	if len(Pokemon_List) != 0 {

		PokemonNumber := strconv.Itoa(Pokemon_List_Info.Realmax)

		Pokemon_List[PokemonNumber] = Pokemon{
			Name:      PokemonName,
			Level:     PokemonLevel,
			IV:        "-",
			NewNumber: PokemonNumber,
		}

		Pokemon_List_Info.Names += PokemonName + ","
		Pokemon_List_Info.Realmax += 1
		Pokemon_List_Info.Array += 1
		SavePokemonList()
		Websocket_SendPokemonList()
	}

	Guild_Spawn, err := s.Guild(msg.GuildID)
	if err != nil {
		if Config.Debug {
			fmt.Println(err)
		}
		return
	}
	Channel_Spawn, err := s.Channel(msg.ChannelID)
	if err != nil {
		if Config.Debug {
			fmt.Println(err)
		}
		return
	}
	color.HiBlue("You caught a " + PokemonName + " !")
	Notif_PokeCaught(PokemonName, Guild_Spawn.Name, Channel_Spawn.Name)
}
