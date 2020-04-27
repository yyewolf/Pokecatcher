//go:generate goversioninfo
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
	"fyne.io/fyne/widget"
	"fyne.io/fyne/canvas"
	"github.com/bwmarrin/discordgo"
	"github.com/gorilla/websocket"
	"github.com/gobuffalo/packr"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

//Terminal

var Logs *widget.Box
var LogScroll *widget.ScrollContainer
var LastPokemonImg *canvas.Image
var LastPokemonLabel *widget.Label
var ProgressBar *widget.ProgressBar
var App fyne.App

//Important

var Config ConfigStruct
var Aliases map[string][]string
var Pokemon_List map[string]interface{}
var Pokemon_List_Info PokeListInfoStruct
var Connections []*websocket.Conn
var Websocket_Receive_Functions map[string]func(request Receive_Request)
var DiscordSession *discordgo.Session
var LatestPokemon LatestPokemonType

// Refreshes

var RefreshingMoves bool
var RefreshingMovesChannelID string
var RefreshingList bool
var RefreshingListChannelID string

var SelectedPokemon SelectedPoke

var wgPokeSpawn sync.WaitGroup

//ServerWhitelist : map[guildID]State, meaning it will catch if State is true
var ServerWhitelist map[string]bool

// Spammer Variables
var SpamInterval int
var SpamMessage string
var SpamState bool
var SpamChannel chan (int)

// Readyness
var Ready bool
var isHosted bool

//Stdout
var OStdout *os.File

func check(e error) {
	if e != nil {
		if Config.Debug {
			fmt.Println(e)
		}
		return
	}
}

func OpenBrowser(url string) {
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

func GuildCreate(s *discordgo.Session, event *discordgo.GuildUpdate) {
	if !Ready {
		Ready = true
		LogGreenLn(Logs, "The bot is ready to be used !")
	}

	DiscordSession = s
	return
}

func Useful_Variables() {
	StartLogger() //Will log crashes if it happens.
	LoadConfig()  // Will load config.json file into the program.
	LoadAliases() // Will load aliases.json file.
	LogYellowLn(Logs, "Your config file has been successfully imported !")
	Pokemon_List = make(map[string]interface{}) //Where the Pokemon List of the user will be stored.
	LoadPokemonList()                           // Will load the Users Pokémons list.
	ServerWhitelist = make(map[string]bool)     //Where the Whitelist of the servers will be stored.
	LoadWhitelist()                             // Will load server_whitelist into ServerWhitelist.
	LogYellowLn(Logs, "The server whitelist has been successfully imported !")
	Websocket_Receive_Functions = make(map[string]func(request Receive_Request))
	Websocket_Receive_AllFunctions()
	Login() //Logins to discord
}

func main() {
	box = packr.NewBox("./www")
	Ready = false
	isHosted = false
	//Launches UI
	UI()
}

func Login() {
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New(Config.Token)
	if err != nil {
		if Config.Debug {
			fmt.Println(err)
		}
		LogRedLn(Logs, "Cannot connect to discord, check your token !")
	}
	LogYellowLn(Logs, "The website is being hosted you can connect to it on : http://localhost:"+strconv.Itoa(Config.WebPort))
	dg.LogLevel = -1
	dg.AddHandler(botReady)
	dg.AddHandler(CheckForPokemon)
	dg.AddHandler(SuccessfulCatch)
	dg.AddHandler(CheckForCommand)
	dg.AddHandler(GuildCreate)
	err = dg.Open()
	check(err)
	// Wait here until CTRL-C or other term signal is received.
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func botReady(session *discordgo.Session, evt *discordgo.Ready) {
	LogGreenLn(Logs, "Successfully connected to discord !")
	CheckLicences(session)

	if !isHosted {
		Host_Website() // Starts hosting the website.
		go OpenBrowser("http://localhost:" + strconv.Itoa(Config.WebPort))
	}
}
