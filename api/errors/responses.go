package errors

import (
	"os"
	"strings"
	"unicode"

	"github.com/bwmarrin/discordgo"
	"github.com/go-errors/errors"
	"github.com/smallpepperz/sachibotgo/api/logger"
)

func HandleError(ds *discordgo.Session, i *discordgo.InteractionCreate, err error) {
	logger.Err().Println("Command '"+i.ApplicationCommandData().Name+"' errored!\n", errors.Wrap(err, 2).ErrorStack())
	error_text := strings.ReplaceAll(err.Error(), os.Getenv("USER"), "user")
	r := []rune(error_text)
	r[0] = unicode.ToUpper(r[0])
	error_text = string(r)

	respondToUser(ds, i, error_text)
}

func respondToUser(ds *discordgo.Session, i *discordgo.InteractionCreate, error_text string) {
	err := ds.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Something went wrong!",
					Description: error_text,
					Color:       0xFF0000,
				},
			},
			Flags: uint64(discordgo.MessageFlagsEphemeral),
		}})
	if err != nil {
		logger.Err().Println("Error reporting error:", err)
	}
}
