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

	"github.com/bwmarrin/discordgo"
	"github.com/fatih/color"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

var Config ConfigStruct
var Pokemon_List map[string]interface{}
var Connections []*websocket.Conn
var Websocket_Receive_Functions map[string]func(request Receive_Request)
var DiscordSession *discordgo.Session

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

func check(e error) {
	if e != nil {
		fmt.Println(e)
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

func Useful_Variables() {
	StartLogger() //Will log crashes if it happens.
	LoadConfig()  // Will load config.json file into the program.
	color.Yellow("Your config file has been successfully imported !")
	Pokemon_List = make(map[string]interface{}) //Where the Pokemon List of the user will be stored.
	LoadPokemonList()                           // Will load the Users Pokémons list.
	ServerWhitelist = make(map[string]bool)     //Where the Whitelist of the servers will be stored.
	LoadWhitelist()                             // Will load server_whitelist into ServerWhitelist.
	color.Yellow("The server whitelist has been successfully imported !")
	Websocket_Receive_Functions = make(map[string]func(request Receive_Request))
	Websocket_Receive_AllFunctions()
}

func main() {
	// Create a new Discord session using the provided bot token.
	Useful_Variables()

	dg, err := discordgo.New(Config.Token)
	if err != nil {
		color.Red("Cannot connect to discord, check your token !")
	}
	color.Yellow("The website is being hosted you can connect to it on : http://localhost:" + strconv.Itoa(Config.WebPort))
	dg.AddHandler(botReady)
	dg.AddHandler(CheckForPokemon)
	dg.AddHandler(SuccessfulCatch)
	dg.AddHandler(CheckForCommand)
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
	color.Green("Successfully connected to discord !")

	CheckLicences(session)
	go OpenBrowser("http://localhost:" + strconv.Itoa(Config.WebPort))
	Host_Website() // Starts hosting the website.
}
