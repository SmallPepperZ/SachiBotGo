package invite

import (
	"github.com/bwmarrin/discordgo"
	"github.com/smallpepperz/sachibotgo/api"
	"github.com/smallpepperz/sachibotgo/api/database"
)

type InviteCommandLink struct{}

func (*InviteCommandLink) GetOptions() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Name:        "link",
		Description: "Gets an invite link for an approved potential invite",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "The user to get a link for",
				Required:    true,
			},
		},
		Type: discordgo.ApplicationCommandOptionSubCommand,
	}
}

func (*InviteCommandLink) RunCommand(ds *discordgo.Session, i *discordgo.InteractionCreate) {
	user := i.ApplicationCommandData().Options[0].Options[0].UserValue(ds)
	potentialInvite, err := database.GetPotentialInvite(user.ID)
	if err != nil {
		api.RespondWithError(ds, i, err)
		return
	}
	potentialInvite.InviteStatus = database.InviteStatuses.Invited
	potentialInvite.Save()
	updateEmbed(ds, potentialInvite)
}
