package invite

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/smallpepperz/sachibotgo/api/config"
	"github.com/smallpepperz/sachibotgo/api/database"
	"github.com/smallpepperz/sachibotgo/api/errors"
	"github.com/smallpepperz/sachibotgo/api/logger"
)

type InviteCommandAdd struct{}

func (*InviteCommandAdd) GetOptions() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Name:        "add",
		Description: "Adds a user to the potential invite system",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "The user to add",
				Required:    true,
			},
		},
		Type: discordgo.ApplicationCommandOptionSubCommand,
	}
}

func (*InviteCommandAdd) RunCommand(ds *discordgo.Session, i *discordgo.InteractionCreate) {
	user := i.ApplicationCommandData().Options[0].Options[0].UserValue(ds)

	if _, err := database.GetPotentialInvite(user.ID); err == nil {
		errors.HandleException(ds, i, fmt.Errorf("%s is already in the invite system", user.Username))
		return
	} else {
		logger.Out().Println("Ignore error")
	}

	embed := generateEmbed(user, i.Member.User, i.Member.User, database.InviteStatuses.Active)
	message, err := ds.ChannelMessageSendEmbed(config.InviteChannel, embed)

	if err != nil {
		errors.HandleException(ds, i, err)
		return
	}

	potentialInvite := database.PotentialInvite{
		UserID:           user.ID,
		InviterID:        i.Member.User.ID,
		UpdaterID:        i.Member.User.ID,
		InviteMessageID:  message.ID,
		InviteStatusName: database.InviteStatuses.Active.Name,
	}

	tx := database.Get().Create(&potentialInvite)
	if tx.Error != nil {
		errors.HandleException(ds, i, tx.Error)
	}

	thread, err := ds.MessageThreadStart(message.ChannelID, message.ID, user.Username, 10080)

	if err != nil {
		errors.HandleException(ds, i, err)
		return
	}

	ds.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprint("Added ", user.Username, " to the invite system\nThread at ", thread.Mention()),
			Flags:   uint64(discordgo.MessageFlagsEphemeral),
		},
	})

}
