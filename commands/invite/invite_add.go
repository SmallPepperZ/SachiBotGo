package invite

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/smallpepperz/sachibotgo/api"
	"github.com/smallpepperz/sachibotgo/api/config"
	"github.com/smallpepperz/sachibotgo/api/database"
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

	if entry, _ := database.GetPotentialInvite(user.ID); entry != nil {
		api.RespondWithError(ds, i, fmt.Errorf("%s is already in the invite system", user.Username))
		return
	}

	embed := generateEmbed(user, i.Member.User, i.Member.User, database.InviteStatuses.Active)
	message, err := ds.ChannelMessageSendEmbed(config.InviteChannel, embed)

	if err != nil {
		api.RespondWithError(ds, i, err)
		return
	}

	potentialInvite := database.PotentialInvite{
		UserID:          user.ID,
		InviterID:       i.Member.User.ID,
		UpdaterID:       i.Member.User.ID,
		InviteMessageID: message.ID,
		InviteStatus:    database.InviteStatuses.Active,
	}

	tx := database.Get().Create(&potentialInvite)
	if tx.Error != nil {
		api.RespondWithError(ds, i, tx.Error)
	}

	ds.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprint("Added ", user.Username, " to the invite system"),
			Flags:   uint64(discordgo.MessageFlagsEphemeral),
		},
	})

}
