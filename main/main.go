package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/smallpepperz/sachibotgo/api"
	"github.com/smallpepperz/sachibotgo/api/config"
	"github.com/smallpepperz/sachibotgo/api/logger"
	"github.com/smallpepperz/sachibotgo/commands"
)

var Session *discordgo.Session

func main() {
	Session, _ = discordgo.New(config.Discord.Token)
	defer Session.Close()

	addHandlers()
	setApplicationId()
	commands.AddAllCommands()
	api.LoadCommands(Session, []string{"all"})

	Session.Identify.Intents = discordgo.IntentsAll

	err := Session.Open()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	setUptime()

	fmt.Println(`Now running. Press CTRL-C to exit.`)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	fmt.Println("Removing commands")

	unregisterCommands()

	fmt.Println("Shutting down")
}

func unregisterCommands() {
	for _, guild := range config.Discord.Guilds {
		registeredCommands, err := Session.ApplicationCommands(Session.State.User.ID, guild)
		if err != nil {
			logger.Err().Printf("Could not fetch registered commands: %v", err)
		}
		for _, v := range registeredCommands {
			err := Session.ApplicationCommandDelete(Session.State.User.ID, guild, v.ID)
			if err != nil {
				logger.Err().Printf("Cannot delete '%v' command: %v", v.Name, err)
			}
		}
	}
}

func setStatus() {
	usd := discordgo.UpdateStatusData{
		Status: string(discordgo.StatusDoNotDisturb),
		Activities: []*discordgo.Activity{
			{
				Name: "time pass, never to return",
				Type: discordgo.ActivityTypeWatching,
			},
		},
	}
	Session.UpdateStatusComplex(usd)
}

func setApplicationId() {
	u, err := url.Parse("https://discord.com/api/oauth2/applications/@me")
	if err != nil {
		panic(err)
	}

	request := &http.Request{
		Method: "GET",
		URL:    u,
		Header: map[string][]string{
			"Authorization": {config.Discord.Token},
		},
	}

	response, err := Session.Client.Do(request)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	var app discordgo.Application
	json.NewDecoder(response.Body).Decode(&app)
	api.Globals.AppID = app.ID
}

func setUptime() {
	api.Globals.StartTime = time.Now()
}

func addHandlers() {
	Session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		logger.Out().Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
		setStatus()
	})
	Session.AddHandler(func(s *discordgo.Session, r *discordgo.Resumed) {
		setStatus()
	})
}
