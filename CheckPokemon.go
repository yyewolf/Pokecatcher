package main

import (
	"fmt"
	"image/png"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne"

	"fyne.io/fyne/widget"
	"github.com/bwmarrin/discordgo"
	"github.com/nfnt/resize"
)

func imageToString(URL string) string {
	response, e := http.Get(URL)
	//response => image
	check(e)
	defer response.Body.Close()
	//Closes the web page when it's done
	image, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logDebug("[ERROR]", err)
		return "nothing"
	}
	return string(image)
}

func logPokemonSpawn(PokemonName string, GuildName string, ChannelName string, Accuracy float64, AliasUsed string) {
	wgPokeSpawn.Wait()
	wgPokeSpawn.Add(1)

	if len(logs.Children)+6 > 150 {
		logs.Children = []fyne.CanvasObject{}
		logs.Refresh()
		logBlueLn(logs, "The console has been cleared automatically.")
	}

	logGreenLn(logs, "-------------------------------------------------------")
	logs.Append(widget.NewHBox(greenTXT("A"), blueTXT(PokemonName), greenTXT("has spawned on :")))
	f := fmt.Sprintf("%f", Accuracy)
	logGreenLn(logs, "Guild Name : "+GuildName)
	logGreenLn(logs, "Channel Name : #"+ChannelName)
	logGreenLn(logs, "Accuracy : "+f+"%")
	logGreenLn(logs, "Alias used : "+AliasUsed)

	wgPokeSpawn.Done()
}

func fakeTalk(s *discordgo.Session, ChannelID string, Letters int) {
	//Fakes user typing
	for start := time.Now(); time.Since(start) < time.Duration(config.Delay)*time.Millisecond; {
		_ = s.ChannelTyping(ChannelID)
		time.Sleep(time.Duration(config.Delay/Letters) * time.Millisecond)
	}
}

func hasAliases(Pokemon string) bool {
	if _, ok := aliases[Pokemon]; ok {
		return true
	}
	return false
}

func checkForPokemon(s *discordgo.Session, msg *discordgo.MessageCreate) {
	//Check if the person is allowed
	if !config.IsAllowedToUse {
		return
	}
	//Check if the server is whitelisted
	if !serverWhitelist[msg.GuildID] {
		return
	}
	//Check if the author is pokecord
	if msg.Author.ID != botID {
		return
	}
	discordSession = s
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
	//STARTS DETECTING HERE

	ImageURL := msg.Embeds[0].Image.URL
	ImageString := imageToString(ImageURL)
	SpawnedPokemonName := ""
	ImageDecoded, err := loadImg(ImageString)
	if err != nil {
		logDebug("[ERROR]", err)
		return
	}
	ImageResized := resize.Resize(64, 64, ImageDecoded, resize.Bicubic)
	Buffer := &buf{}
	_ = png.Encode(Buffer, ImageResized)
	ImageResized, _ = png.Decode(Buffer)
	List := box.List()
	Accuracy := 0.0
	isInWhitelist := false

	for i := range List {
		if strings.Contains(List[i], "img") {
			//Gets rid of the path debris
			Name := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(List[i], "img/", ""), "img\\", ""), ".png", "")

			ScanImage := decodedImages[Name]
			Accuracy = compareIMG(ScanImage, ImageResized)
			if Accuracy < 0.35 {
				//Check if the Pokémon is in whitelist (now because of Nidoran)
				if pokemonWhitelist[Name] {
					isInWhitelist = true
				}
				SpawnedPokemonName = strings.ReplaceAll(strings.ReplaceAll(Name, "♀", ""), "♂", "")

				lastPokemonImg.Image = ScanImage
				lastPokemonLabel.SetText(SpawnedPokemonName)
				lastPokemonImg.Refresh()
				break
			}
		}
	}

	//STOPS DETECTING HERE

	Accuracy = 100.0 - Accuracy

	if SpawnedPokemonName == "" {
		return
	}
	GuildSpawn, err := s.Guild(msg.GuildID)
	if err != nil {
		logDebug("[ERROR]", err)
		return
	}
	ChannelSpawn, err := s.Channel(msg.ChannelID)
	if err != nil {
		logDebug("[ERROR]", err)
		return
	}
	//Logs info into the console and sends a notification to the website.
	OriginalName := SpawnedPokemonName

	CommandToCatch := ""

	logPokemonSpawn(OriginalName, GuildSpawn.Name, ChannelSpawn.Name, Accuracy, OriginalName)

	notifPokeSpawn(OriginalName, GuildSpawn.Name, CommandToCatch, ChannelSpawn.Name, ChannelSpawn.ID)
	logCyanLn(logs, "Command : "+OriginalName)

	latestPokemon = latestPokemonType{
		ChannelID: msg.ChannelID,
		Command:   strings.ToLower(OriginalName),
	}

	if config.AutoCatching && isInWhitelist {
		//Verifies that it isn't a duplicate if it's ON
		if config.Duplicate {
			if !strings.Contains(pokemonListInfo.Names, OriginalName) {
				return
			}
		}
		//Closes spammer
		if spamState {
			spamChannel <- 1
		}
		fakeTalk(s, msg.ChannelID, len(CommandToCatch+" "+strings.ToLower(OriginalName)))

		rand.Seed(time.Now().UnixNano())
		RandomNess := rand.Intn(422) - rand.Intn(221)

		logDebug("[DEBUG]", "Waiting to catch a", OriginalName)

		time.Sleep(time.Duration(config.Delay+RandomNess) * time.Millisecond)
		logBlueLn(logs, "Tried to catch your : "+OriginalName)

		_, err := s.ChannelMessageSend(msg.ChannelID, CommandToCatch+" "+strings.ToLower(OriginalName))

		logDebug("[DEBUG]", "Sent message to catch a", OriginalName)
		if err != nil {
			logDebug("[ERROR]", err)
			notifCatchingErr(OriginalName, GuildSpawn.Name, ChannelSpawn.Name)
			return
		}
		//Restart spammer
		if spamState {
			go spamFunc(discordSession, config.ChannelID, spamInterval, spamMessage)
		}
	}
}

func successfulCatchUpdate(s *discordgo.Session, msg *discordgo.MessageUpdate) {
	defer func() {
		if r := recover(); r != nil {
			//
		}
	}()
	m := discordgo.MessageCreate{}

	m.Embeds = msg.Embeds
	m.Author = msg.Author
	successfulCatch(s, &m)
}

func successfulCatch(s *discordgo.Session, msg *discordgo.MessageCreate) {
	//Check if the person is allowed
	if !config.IsAllowedToUse {
		return
	}
	//Check if the author is pokecord
	if msg.Author.ID != botID {
		return
	}
	//Check if there is an embed
	if len(msg.Embeds) == 0 {
		return
	}
	//Check if it's a pokemon catch
	if !strings.Contains(msg.Embeds[0].Description, "Congratulations") {
		return
	}
	//Check if the author is mentioned
	if !strings.Contains(msg.Embeds[0].Description, s.State.User.ID) {
		return
	}
	//Gets the command from the message : "Congratulations @User, you have caught a Foongus!"
	//Steps :
	//99 Pokemon! Added to Pokédex."
	//99 Pokemon
	//99Pokemon
	//Pokemon
	PokemonName := strings.Split(msg.Embeds[0].Description, "a ")[1]
	PokemonName = strings.Split(PokemonName, "!")[0]
	//Sends Debug infos
	logDebug("[DEBUG] Caught a ", PokemonName)
	//Adds the pokemon to the list

	pokemonListInfo.Names += PokemonName + ","
	pokemonListInfo.Realmax++
	pokemonListInfo.Array++

	PokemonNumber := strconv.Itoa(pokemonListInfo.Realmax)

	pokemonList[PokemonNumber] = pokemon{
		Name:      PokemonName,
		Level:     "-",
		IV:        "-",
		NewNumber: PokemonNumber,
	}

	savePokemonList()
	websocketSendPokemonList()

	GuildSpawn, err := s.Guild(msg.GuildID)
	if err != nil {
		logDebug("[ERROR]", err)
		return
	}
	ChannelSpawn, err := s.Channel(msg.ChannelID)
	if err != nil {
		logDebug("[ERROR]", err)
		return
	}
	logCyanLn(logs, "You caught a "+PokemonName+" !")
	notifPokeCaught(PokemonName, GuildSpawn.Name, ChannelSpawn.Name)

	if !config.GoodFilter && !config.CustomFilters {
		//Will alert the user if necessary
		if config.Alerts && config.AlertChannelID != "" {
			text := ""
			//Pokemon is legendary
			for i := range legendaries {
				if PokemonName == legendaries[i] {
					text += " legendary"
				}
			}
			//Pokemon is shiny
			if strings.Contains(msg.Content, ":star:") {
				text += " shiny"
			}
			_, err = s.ChannelMessageSend(config.AlertChannelID, "You caught a"+text+" "+PokemonName+" !")
			if err != nil {
				logDebug("[ERROR]", err)
				return
			}
		}
		return
	}

	logDebug("[DEBUG] Will verify a ", PokemonName)

	time.Sleep(3 * time.Second)
	//Will release the pokémon if it is bad
	infoMenu.Activated = true
	infoMenu.AutoRelease = true
	infoMenu.ChannelID = config.ChannelID
	m, err := s.ChannelMessageSend(config.ChannelID, config.PrefixPokecord+"info latest")
	if err != nil {
		logDebug("[ERROR]", err)
		return
	}
	infoMenu.MessageID = m.ID
}
