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
	"github.com/smallpepperz/sachibotgo/commands"
)

var Session *discordgo.Session

func main() {
	Session, _ = discordgo.New(api.Config.Discord.Token)
	defer Session.Close()

	setApplicationId()
	commands.AddAllCommands()
	api.LoadCommands(Session, []string{"all"})

	Session.Identify.Intents = discordgo.IntentsAll

	err := Session.Open()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	setStatus()
	setUptime()

	fmt.Println(`Now running. Press CTRL-C to exit.`)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	fmt.Println("Shutting down")
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
			"Authorization": {api.Config.Discord.Token},
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
