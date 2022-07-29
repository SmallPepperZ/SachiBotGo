package invite

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/smallpepperz/sachibotgo/api"
	"github.com/smallpepperz/sachibotgo/api/database"
)

type InviteCommandUpdate struct{}

var InviteStatusesMap = map[string]database.InviteStatus{
	"Active":   database.InviteStatuses.Active,
	"Approved": database.InviteStatuses.Approved,
	"Accepted": database.InviteStatuses.Accepted,
	"Rejected": database.InviteStatuses.Rejected,
	"Declined": database.InviteStatuses.Declined,
	"Paused":   database.InviteStatuses.Paused,
}

func (*InviteCommandUpdate) GetOptions(ds *discordgo.Session) *discordgo.ApplicationCommandOption {
	statusOptions := make([]*discordgo.ApplicationCommandOptionChoice, 0, 6)
	users := make([]*database.PotentialInvite, 0)

	for key := range InviteStatusesMap {
		statusOptions = append(statusOptions, &discordgo.ApplicationCommandOptionChoice{
			Name:  key,
			Value: key,
		})
	}
	database.Get().Find(&users)

	userOptions := make([]*discordgo.ApplicationCommandOptionChoice, 0, len(users))
	for _, user := range users {
		discordUser, err := user.User(ds)
		var username string
		if err == nil {
			username = discordUser.Username
		} else {
			username = user.UserID
		}
		userOptions = append(userOptions, &discordgo.ApplicationCommandOptionChoice{
			Name:  username,
			Value: user.UserID,
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
				Choices:     userOptions,
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
				api.RespondWithError(ds, i, fmt.Errorf("could not find user '%s'", value.StringValue()))
				return
			}
		case "action":
			action = value.StringValue()
		case "force":
			force = value.BoolValue()
		}
	}
	potentialInvite, err := database.GetPotentialInvite(user.ID)
	if potentialInvite == nil || err != nil {
		api.RespondWithError(ds, i, fmt.Errorf("could not find user '%s' in the invite system", user.Username))
		return
	}

	switch {
	case force:
	case potentialInvite.InviteStatus == database.InviteStatuses.Accepted:
		sendResponse(ds, i, "User is marked as having accepted their invitation. Use the force flag to update this user")
		return
	case potentialInvite.InviteStatus == database.InviteStatuses.Declined:
		sendResponse(ds, i, "User is marked as having declined their invitation. Use the force flag to update this user")
		return
	}
	potentialInvite.InviteStatus = InviteStatusesMap[action]

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
		api.RespondWithError(ds, i, err)
	}
}
