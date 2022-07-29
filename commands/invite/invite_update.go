package invite

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/smallpepperz/sachibotgo/api/database"
	"github.com/smallpepperz/sachibotgo/api/errors"
)

type InviteCommandUpdate struct{}

func (*InviteCommandUpdate) GetOptions(ds *discordgo.Session) *discordgo.ApplicationCommandOption {
	statusOptions := make([]*discordgo.ApplicationCommandOptionChoice, 0, 6)

	for key := range database.InviteStatusesMap {
		statusOptions = append(statusOptions, &discordgo.ApplicationCommandOptionChoice{
			Name:  key,
			Value: key,
		})
	}

	return &discordgo.ApplicationCommandOption{
		Name:        "update",
		Description: "Updates a potential invite",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "user",
				Description: "The user to update",
				Required:    true,
				Choices:     getUsers(ds),
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "status",
				Description: "The status to set the user to. If unset, it will refresh their user information",
				Required:    false,
				Choices:     statusOptions,
			},
			{
				Type:        discordgo.ApplicationCommandOptionBoolean,
				Name:        "force",
				Description: "Set this flag to bypass checks",
				Required:    false,
			},
		},
		Type: discordgo.ApplicationCommandOptionSubCommand,
	}
}

func (*InviteCommandUpdate) RunCommand(ds *discordgo.Session, i *discordgo.InteractionCreate) {
	var user *discordgo.User
	var action string
	var force bool
	var err error
	for _, value := range i.ApplicationCommandData().Options[0].Options {
		switch value.Name {
		case "user":
			user, err = ds.User(value.StringValue())
			if err != nil {
				errors.HandleException(ds, i, fmt.Errorf("could not find user '%s'", value.StringValue()))
				return
			}
		case "status":
			action = value.StringValue()
		case "force":
			force = value.BoolValue()
		}
	}
	potentialInvite, err := database.GetPotentialInvite(user.ID)
	if potentialInvite == nil || err != nil {
		errors.HandleError(ds, i, fmt.Errorf("could not find user '%s' in the invite system", user.Username))
		return
	}

	switch {
	case force:
	case potentialInvite.InviteStatus() == database.InviteStatuses.Accepted:
		sendResponse(ds, i, "User is marked as having accepted their invitation. Use the force flag to update this user")
		return
	case potentialInvite.InviteStatus() == database.InviteStatuses.Declined:
		sendResponse(ds, i, "User is marked as having declined their invitation. Use the force flag to update this user")
		return
	}
	if action != "" {
		potentialInvite.InviteStatusName = action
	}

	potentialInvite.Save()
	updateEmbed(ds, potentialInvite)
	sendResponse(ds, i, fmt.Sprintf("Sucessfully updated status for %s", user.Username))
}

func sendResponse(ds *discordgo.Session, i *discordgo.InteractionCreate, response string) {
	err := ds.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
		},
	})
	if err != nil {
		errors.HandleException(ds, i, err)
	}
}
