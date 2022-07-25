package api

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/smallpepperz/sachibotgo/api/logger"
)

var availableCommands = make(map[string]Command, 0)
var loadedCommands = make(map[string]Command, 0)

func LoadCommands(ds *discordgo.Session, commands []string) {
	loadedCommands = availableCommands

	for k, v := range loadedCommands {
		v.Load(ds)
		fmt.Printf("Loaded %s\n", k)
	}
}

func AddCommand(command Command) {
	availableCommands[command.Name()] = command
}

func GetLoadedCommands() map[string]Command {
	return loadedCommands
}
func GetAvailableCommands() map[string]Command {
	return availableCommands
}
/*
	A command to be registered with the bot.
*/
type Command interface {
	/*
		An entrypoint function for the command.
	*/
	Load(session *discordgo.Session)
	Name() string
}


func CreateGlobalCommand(ds *discordgo.Session, command *discordgo.ApplicationCommand, cmdFunction func(ds *discordgo.Session, i *discordgo.InteractionCreate)) []error {
	return CreateCommand(ds, Config.Discord.Guilds, command, cmdFunction)
}

func RegisterCommand(ds *discordgo.Session, guilds []string, command *discordgo.ApplicationCommand) (errs []error) {
	ds.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		for _, v := range guilds {
			logger.Out().Printf("Registering command '%s' for guild %s\n", command.Name, v)
			_, err := s.ApplicationCommandCreate(Globals.AppID, v, command)
			if err != nil {
				errs = append(errs, err)
				logger.Err().Printf("Cannot create slash command %q: %v", command.Name, err)
			}
		}
	})
	return
}

func CreateCommand(ds *discordgo.Session, guilds []string, command *discordgo.ApplicationCommand, cmdFunction func(ds *discordgo.Session, i *discordgo.InteractionCreate)) (errs []error) {
	errs = RegisterCommand(ds, guilds, command)
	ds.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			{
				if i.ApplicationCommandData().Name == command.Name {
					cmdFunction(s, i)
				}
			}
		}
	})
	return
}