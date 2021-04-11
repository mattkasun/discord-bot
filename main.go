package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// Configuration
type Configuration struct {
	Token string
}

var configuration Configuration

func init() {

	//flag.StringVar(&Token, "t", "", "Bot Token")
	//flag.Parse()
	file, err := os.Open("/home/mkasun/discord-bot/discord-bot.conf")
	if err != nil {
		panic("could not read token file")
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&configuration)
	if err != nil {
		panic("could not read token")
	}

}

func main() {

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + configuration.Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	switch m.Content {
	// If the message is "ping" reply with "Pong!"
	case "ping":
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	// If the message is "pong" reply with "Ping!"
	case "pong":
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	case "help", "!help":
		s.ChannelMessageSend(m.ChannelID, "I understand the following commands:\nping\npong\nfortune\nweather")
	case "fortune":
		fortune, _ := exec.Command("/usr/games/fortune").Output()
		s.ChannelMessageSend(m.ChannelID, string(fortune))
	case "weather":
		ottWeather, _ := exec.Command("/usr/bin/curl", "wttr.in/ottawa?format=3").Output()
		vanWeather, _ := exec.Command("/usr/bin/curl", "wttr.in/vancouver?format=3").Output()
		weather := fmt.Sprintf("%s%s", ottWeather, vanWeather)
		s.ChannelMessageSend(m.ChannelID, weather)
	default:
		fmt.Println("Recieved: ", m.Content)
	}
}
