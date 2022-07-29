package invite

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/smallpepperz/sachibotgo/api/database"
	"github.com/smallpepperz/sachibotgo/api/errors"
)

type InviteCommandLink struct{}

func (*InviteCommandLink) GetOptions(ds *discordgo.Session) *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Name:        "link",
		Description: "Gets an invite link for an approved potential invite",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "user",
				Description: "The user to get a link for",
				Required:    true,
				Choices:     getUsers(ds),
			},
		},
		Type: discordgo.ApplicationCommandOptionSubCommand,
	}
}

func (*InviteCommandLink) RunCommand(ds *discordgo.Session, i *discordgo.InteractionCreate) {
	user, err := ds.User(i.ApplicationCommandData().Options[0].Options[0].StringValue())
	if err != nil {
		errors.HandleException(ds, i, err)
		return
	}
	potentialInvite, err := database.GetPotentialInvite(user.ID)
	if err != nil {
		errors.HandleException(ds, i, err)
		return
	}
	var invite *discordgo.Invite
	if potentialInvite.InviteCode != "" {
		invite, err = ds.Invite(potentialInvite.InviteCode)
		if err != nil {
			invite, err = createInvite(ds, i)
		}
	} else {
		invite, err = createInvite(ds, i)
	}
	if err != nil {
		errors.HandleException(ds, i, err)
		return
	}

	potentialInvite.InviteStatusName = database.InviteStatuses.Invited.Name
	potentialInvite.InviteCode = invite.Code
	potentialInvite.Save()

	ds.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Invite for %s: https://discord.gg/%s", user.Mention(), invite.Code),
			Flags:   uint64(discordgo.MessageFlagsEphemeral),
		},
	})
	updateEmbed(ds, potentialInvite)
}

func createInvite(ds *discordgo.Session, i *discordgo.InteractionCreate) (invite *discordgo.Invite, err error) {
	guild, _ := ds.Guild(i.GuildID)
	invite, err = ds.ChannelInviteCreate(guild.RulesChannelID, discordgo.Invite{
		MaxUses: 1,
		MaxAge:  0,
	})
	return
}
