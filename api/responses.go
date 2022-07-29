package api

import (
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/go-errors/errors"
	"github.com/smallpepperz/sachibotgo/api/logger"
)

func RespondWithError(ds *discordgo.Session, i *discordgo.InteractionCreate, err error) {
	logger.Err().Println("Command '"+i.ApplicationCommandData().Name+"' errored!\n", errors.Wrap(err, 2).ErrorStack())
	error_text := strings.ReplaceAll(err.Error(), os.Getenv("USER"), "user")

	err = ds.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "An error occurred",
			Embeds: []*discordgo.MessageEmbed{
				{
					Description: "`" + error_text + "`",
					Color:       0xF00,
				},
			},
			Flags: uint64(discordgo.MessageFlagsEphemeral),
		}})
	if err != nil {
		logger.Err().Println("Error reporting error:", err)
	}
}
