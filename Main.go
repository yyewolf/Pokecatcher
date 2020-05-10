package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strconv"
	"sync"
	"syscall"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/widget"
	"github.com/bwmarrin/discordgo"
	"github.com/gobuffalo/packr"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

//Window
var logs *widget.Box
var logScroll *widget.ScrollContainer
var lastPokemonImg *canvas.Image
var lastPokemonLabel *widget.Label
var currentPokemonImg *canvas.Image
var currentPokemonLevel *widget.Label
var progressBar *widget.ProgressBar
var appli fyne.App
var windowIcon fyne.Resource

//Important
var config ConfigStruct
var aliases map[string][]string
var pokemonList map[string]pokemon
var pokemonWhitelist map[string]bool
var pokemonListInfo PokeListInfoStruct
var connections []*websocket.Conn
var websocketReceiveFunctions map[string]func(request receiveRequest)
var discordSession *discordgo.Session
var latestPokemon latestPokemonType

// Refreshes
var refreshingMoves bool
var refreshingMovesChannelID string
var refreshingList bool
var refreshingListChannelID string

var selectedPokemon SelectedPoke

var wgPokeSpawn sync.WaitGroup

//serverWhitelist : map[guildID]State, meaning it will catch if State is true
var serverWhitelist map[string]bool

// Spammer Variables
var spamInterval int
var spamMessage string
var spamState bool
var spamChannel chan (int)

// Readyness
var isReady bool
var isHosted bool

//Stdout
var stdoutKeeped *os.File

//AutoLeveler
var infoMenu infoActivated
var priorityQueue []string

func check(e error) {
	if e != nil {
		logDebug("[ERROR] ", e)
		return
	}
}

func openBrowser(url string) {
	var err error
	time.Sleep(2 * time.Second)
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		// Not important :/
	}
}

func guildCreate(s *discordgo.Session, event *discordgo.GuildUpdate) {
	if !isReady {
		isReady = true
		logGreenLn(logs, "The bot is ready to be used !")
	}

	discordSession = s
	return
}

func usefulVariables() {
	loadConfig()  // Will load config.json file into the program.
	startLogger() //Will log crashes if it happens.
	loadAliases() // Will load aliases.json file.
	logYellowLn(logs, "Your config file has been successfully imported !")
	pokemonList = make(map[string]pokemon)  //Where the Pokemon List of the user will be stored.
	loadPokemonList()                       // Will load the Users Pokémon list.
	loadPokemonWhitelist()                  // Will load the Users Pokémon Whitelist
	serverWhitelist = make(map[string]bool) //Where the Whitelist of the servers will be stored.
	loadWhitelist()                         // Will load server_whitelist into ServerWhitelist.
	logYellowLn(logs, "The server whitelist has been successfully imported !")
	websocketReceiveFunctions = make(map[string]func(request receiveRequest))
	websocketReceiveAllFunctions()
	discordLogin() //Logins to discord
}

func main() {
	box = packr.NewBox("./www")
	isReady = false
	isHosted = false
	_ = os.Chdir("/storage/emulated/0/Android/data/org.golang.todo.Pokecatcher/files")
	//Launches UI
	startUI()
}

func discordLogin() {
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New(config.Token)
	check(err)
	dg.LogLevel = -1

	//Ready event
	dg.AddHandler(botReady)

	//Recognize pokemons
	dg.AddHandler(checkForPokemon)

	//Adds pokemon to list
	dg.AddHandler(successfulCatch)

	//Recognize commands
	dg.AddHandler(checkForCommand)

	//Updates Servers
	dg.AddHandler(guildCreate)

	//AutoLeveling Feature
	dg.AddHandler(InfoVerifier)
	dg.AddHandler(SelectVerifier)
	dg.AddHandler(AutoLeveler)

	//AutoReleaser feature
	dg.AddHandler(AutoRelease)

	err = dg.Open()
	if err != nil {
		logDebug("[ERROR] ", err)
		logRedLn(logs, "Cannot connect to discord, check your token !")
		askUserForToken()
	}
	// Wait here until CTRL-C or other term signal is received.
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func botReady(session *discordgo.Session, evt *discordgo.Ready) {
	logGreenLn(logs, "Successfully connected to discord !")
	checkLicences(session)

	if !isHosted {
		isHosted = true
		logYellowLn(logs, "The website is being hosted you can connect to it on : http://localhost:"+strconv.Itoa(config.WebPort))
		logScroll.SetMinSize(logScroll.MinSize())
		hostWebsite() // Starts hosting the website.
	}
}
