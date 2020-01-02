package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/fatih/color"
)

func ImageToSHA256(URL string) string {
	response, e := http.Get(URL)
	//response => image
	check(e)
	defer response.Body.Close()
	//Closes the web page when it's done
	image, err := ioutil.ReadAll(response.Body)
	check(err)
	h := sha256.New()
	h.Write(image)
	//Encodes the image to sha256 (database is sha256)
	sha256_hash := hex.EncodeToString(h.Sum(nil))
	return sha256_hash
}

func LogPokemonSpawn(PokemonName string, GuildName string, ChannelName string) {
	wgPokeSpawn.Wait()
	wgPokeSpawn.Add(1)
	fmt.Println("")
	color.Green("-------------------------------------------------------")
	PrintGreen := color.New(color.FgHiGreen).PrintfFunc()
	PrintBlue := color.New(color.FgHiBlue).PrintfFunc()
	PrintGreen("A ")
	PrintBlue(PokemonName)
	PrintGreen(" has spawned on : \nGuild Name : " + GuildName + "\nChannel Name : #" + ChannelName + ".")
	fmt.Println("")
	wgPokeSpawn.Done()
}

func FakeTalk(s *discordgo.Session, ChannelID string, Letters int) {
	//Fakes user typing
	for start := time.Now(); time.Since(start) < time.Duration(Config.Delay)*time.Millisecond; {
		_ = s.ChannelTyping(ChannelID)
		time.Sleep(time.Duration(Config.Delay/Letters) * time.Millisecond)
	}
}

func CheckForPokemon(s *discordgo.Session, msg *discordgo.MessageCreate) {
	//Check if the person is allowed
	if !Config.IsAllowedToUse {
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
	//Check if it's a pokemon spawn
	if !strings.Contains(msg.Embeds[0].Title, "A wild") {
		return
	}
	ImageURL := msg.Embeds[0].Image.URL
	ImageHash := ImageToSHA256(ImageURL)
	if _, ok := Hashes_Database[ImageHash]; !ok {
		return
	}
	Spawned_Pokemon_Name := Hashes_Database[ImageHash]
	Guild_Spawn, err := s.Guild(msg.GuildID)
	check(err)
	Channel_Spawn, err := s.Channel(msg.ChannelID)
	check(err)
	//Logs info into the console and sends a notification to the website.
	LogPokemonSpawn(Spawned_Pokemon_Name, Guild_Spawn.Name, Channel_Spawn.Name)
	//Gets the command from the message : "Guess the pokemon and type p!catch <pokémon> to catch it !"
	Command_To_Catch := strings.Split(strings.Split(msg.Embeds[0].Description, "type ")[1], " <po")[0]
	Notif_PokeSpawn(Spawned_Pokemon_Name, Guild_Spawn.Name, Command_To_Catch, Channel_Spawn.Name, Channel_Spawn.ID)
	color.HiBlue("Command : " + Command_To_Catch + " " + Spawned_Pokemon_Name)
	if Config.AutoCatching && ServerWhitelist[msg.GuildID] {
		//Closes spammer
		if SpamState {
			SpamChannel <- 1
		}
		FakeTalk(s, msg.ChannelID, len(Command_To_Catch+" "+strings.ToLower(Spawned_Pokemon_Name)))
		time.Sleep(time.Duration(Config.Delay) * time.Millisecond)
		color.Blue("Tried to catch your : " + Spawned_Pokemon_Name)

		Command_To_Catch = strings.ReplaceAll(Command_To_Catch, "а", "a")
		_, err := s.ChannelMessageSend(msg.ChannelID, Command_To_Catch+" "+strings.ToLower(Spawned_Pokemon_Name))
		if err != nil {
			Notif_CatchingErr(Spawned_Pokemon_Name, Guild_Spawn.Name, Channel_Spawn.Name)
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
		PokemonNumber := fmt.Sprintf("%.0f", Pokemon_List["realmax"].(float64))

		Pokemon_List[PokemonNumber] = Pokemon{
			Name:      PokemonName,
			Level:     PokemonLevel,
			IV:        "-",
			NewNumber: PokemonNumber,
		}

		Pokemon_List["names"] = Pokemon_List["names"].(string) + PokemonName + ","
		Pokemon_List["realmax"] = Pokemon_List["realmax"].(float64) + 1
		Pokemon_List["array"] = Pokemon_List["array"].(float64) + 1
	}

	SavePokemonList()
	Websocket_SendPokemonList()

	Guild_Spawn, err := s.Guild(msg.GuildID)
	check(err)
	Channel_Spawn, err := s.Channel(msg.ChannelID)
	check(err)
	color.HiBlue("You caught a " + PokemonName + " !")
	Notif_PokeCaught(PokemonName, Guild_Spawn.Name, Channel_Spawn.Name)
}
