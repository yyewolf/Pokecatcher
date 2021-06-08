package main

import (
	"encoding/json"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"fyne.io/fyne"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"

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

func logPokemonSpawn(PokemonName string, GuildName string, ChannelName string, AliasUsed string) {
	wgPokeSpawn.Wait()
	wgPokeSpawn.Add(1)

	if len(logs.Children)+6 > 150 {
		logs.Children = []fyne.CanvasObject{}
		logs.Refresh()
		logBlueLn(logs, "The console has been cleared automatically.")
	}

	logGreenLn(logs, "-------------------------------------------------------")
	logs.Append(widget.NewHBox(greenTXT("A"), blueTXT(PokemonName), greenTXT("has spawned on :")))
	logGreenLn(logs, "Guild Name : "+GuildName)
	logGreenLn(logs, "Channel Name : #"+ChannelName)
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

type responsePred struct {
	Name       string `json:"name"`
	ID         string `json:"id"`
	Confidence string `json:"confidence"`
}

type responseAPI struct {
	Predictions []responsePred `json:"predictions"`
	URL         string         `json:"image url"`
}

func removeAccents(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	output, _, e := transform.String(t, s)
	if e != nil {
		panic(e)
	}
	return output
}

func checkAPI(imageURL string, model string) responseAPI {
	resp, err := http.PostForm("https://aipokedex.com/getPoke", url.Values{"url": {imageURL}, "token": {config.AIToken}, "model": {model}})
	var response responseAPI
	if err != nil {
		return response
	}
	defer resp.Body.Close()
	logDebug("[REQ]", resp.Body)
	body, err := io.ReadAll(resp.Body)
	json.Unmarshal(body, &response)
	return response
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
	//Check if the author is a bot
	if !msg.Author.Bot {
		return
	}
	//Check if there is an embed
	if len(msg.Embeds) == 0 {
		return
	}
	//Check if the author name has poke
	if !strings.Contains(removeAccents(strings.ToLower(msg.Author.String())), "poke") && msg.Author.ID != "519850436899897346" {
		return
	}
	discordSession = s
	if msg.Embeds[0].Image == nil {
		return
	}
	//Check if it's a pokemon spawn
	if !strings.Contains(msg.Embeds[0].Title, "wild") {
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
	ImageResized := resize.Resize(uint(128.0*float64(ImageDecoded.Bounds().Dx())/float64(ImageDecoded.Bounds().Dy())), 128, ImageDecoded, resize.Bicubic)
	Buffer := &buf{}
	_ = png.Encode(Buffer, ImageResized)
	ImageResized, _ = png.Decode(Buffer)
	isInWhitelist := false

	model := "classic"
	//For Pokedi and Poketwo
	if msg.Author.ID == "716293342740348948" || msg.Author.ID == "716390085896962058" {
		model = "background"
	}

	APIResp := checkAPI(ImageURL, model)
	if len(APIResp.Predictions) == 0 {
		return
	}
	SpawnedPokemonName = APIResp.Predictions[0].Name
	SpawnedPokemonName = strings.Split(SpawnedPokemonName, " (")[0]
	//Check if the Pokémon is in whitelist (now because of Nidoran)
	if pokemonWhitelist[SpawnedPokemonName] {
		isInWhitelist = true
	}
	SpawnedPokemonName = strings.ReplaceAll(strings.ReplaceAll(SpawnedPokemonName, "♀", ""), "♂", "")
	lastPokemonImg.Image = ImageResized
	lastPokemonLabel.SetText(SpawnedPokemonName)
	lastPokemonImg.Refresh()
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
	CatchName := SpawnedPokemonName
	if hasAliases(SpawnedPokemonName) {
		Names := aliases[OriginalName]
		if len(Names) == 0 {
			CatchName = SpawnedPokemonName
		} else {
			CatchName = Names[0]
		}
		if config.Aliases {
			if len(Names) == 1 {
				CatchName = Names[0]
			} else {
				CatchName = Names[rand.Intn(len(Names)-1)]
			}
		}
	}

	logPokemonSpawn(OriginalName, GuildSpawn.Name, ChannelSpawn.Name, CatchName)
	//Gets the command from the message : "Guess the pokemon and type p!catch <pokémon> to catch it !"
	CommandToCatch := ""
	Split1 := strings.Split(msg.Embeds[0].Description, "type ")
	Split2 := strings.Split(msg.Embeds[0].Description, "Do ")
	if len(Split1) > 1 {
		CommandToCatch = strings.Split(Split1[1], " ")[0]
	} else if len(Split2) > 1 {
		CommandToCatch = strings.Split(Split2[1], " ")[0]
	}
	CommandToCatch = strings.ReplaceAll(CommandToCatch, "а", "a")
	CommandToCatch = strings.ReplaceAll(CommandToCatch, "`", "")
	//Pokécord patched this

	notifPokeSpawn(OriginalName, GuildSpawn.Name, CommandToCatch, ChannelSpawn.Name, ChannelSpawn.ID)
	logCyanLn(logs, "Command : "+CommandToCatch+" "+OriginalName)

	latestPokemon = latestPokemonType{
		ChannelID: msg.ChannelID,
		Command:   CommandToCatch + " " + strings.ToLower(CatchName),
	}

	if config.AutoCatching && isInWhitelist {
		//Verifies that it isn't a duplicate if it's ON
		OriginalName = strings.Split(OriginalName, " (")[0]
		if config.Duplicate {
			if !strings.Contains(pokemonListInfo.Names, OriginalName) {
				return
			}
		}
		//Closes spammer
		if spamState {
			spamChannel <- 1
		}

		rand.Seed(time.Now().UnixNano())
		RandomNess := rand.Intn(422) - rand.Intn(221)

		logDebug("[DEBUG]", "Waiting to catch a", OriginalName)

		time.Sleep(time.Duration(config.Delay+RandomNess) * time.Millisecond)
		logBlueLn(logs, "Tried to catch your : "+OriginalName)

		_, err := s.ChannelMessageSend(msg.ChannelID, CommandToCatch+" "+strings.ToLower(CatchName))

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

func successfulCatch(s *discordgo.Session, msg *discordgo.MessageCreate) {
	//Check if the person is allowed
	if !config.IsAllowedToUse {
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
	PokemonName := strings.Split(msg.Content, "level ")[1]
	PokemonName = strings.Split(PokemonName, "!")[0]
	PokemonName = strings.ReplaceAll(PokemonName, " ", "")
	PokemonName = reg.ReplaceAllString(PokemonName, "")
	//Do the same to detect its level
	PokemonLevel := strings.Split(msg.Content, "level ")[1]
	PokemonLevel = strings.Split(PokemonLevel, " "+PokemonName)[0]
	//Sends Debug infos
	logDebug("[DEBUG] Caught a ", PokemonName)
	//Adds the pokemon to the list

	pokemonListInfo.Names += PokemonName + ","
	pokemonListInfo.Realmax++
	pokemonListInfo.Array++

	PokemonNumber := strconv.Itoa(pokemonListInfo.Realmax)

	pokemonList[PokemonNumber] = pokemon{
		Name:      PokemonName,
		Level:     PokemonLevel,
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
