package ping

import (
	"fmt"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/smallpepperz/sachibotgo/api"
)

type Command struct {
	api.Command
}

func (*Command) Load(ds *discordgo.Session) {
	appCmd := &discordgo.ApplicationCommand{
		Name:        "ping",
		Description: "Checks the latency of the bot",
		Type:        discordgo.ChatApplicationCommand,
	}
	api.CreateGlobalCommand(ds, appCmd, runCommand)
}

func runCommand(ds *discordgo.Session, i *discordgo.InteractionCreate) {
	latency := ds.LastHeartbeatAck.UnixMilli() - ds.LastHeartbeatSent.UnixMilli()

	// set up the embed
	var fields []*discordgo.MessageEmbedField
	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "Latency",
		Value:  strconv.Itoa(int(latency)) + "ms",
		Inline: true,
	})
	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "Uptime",
		Value:  "Started <t:" + strconv.Itoa(int(api.Globals.StartTime.Unix())) + ":R>",
		Inline: true,
	})

	// add image to embed
	embed := &discordgo.MessageEmbed{
		Title:  "Ping",
		Fields: fields,
	}

	err := ds.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:  uint64(discordgo.MessageFlagsEphemeral),
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
	if err != nil {
		fmt.Printf("Failed to send message\n%s", err)
	}
}

func (Command) Name() string {
	return "utility/ping"
}
