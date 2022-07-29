package invite

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/smallpepperz/sachibotgo/api"
	"github.com/smallpepperz/sachibotgo/api/config"
	"github.com/smallpepperz/sachibotgo/api/database"
	"github.com/smallpepperz/sachibotgo/api/errors"
	"github.com/smallpepperz/sachibotgo/api/logger"
)

type Command struct {
	api.Command
}

var commandUpdate = &InviteCommandUpdate{}
var commandAdd = &InviteCommandAdd{}
var commandLink = &InviteCommandLink{}

func (*Command) Load(ds *discordgo.Session) {
	options := []*discordgo.ApplicationCommandOption{
		commandUpdate.GetOptions(ds),
		commandAdd.GetOptions(),
		commandLink.GetOptions(ds),
	}
	appCmd := &discordgo.ApplicationCommand{
		Name:        "invite",
		Description: "The MDSP invite system",
		Type:        discordgo.ChatApplicationCommand,
		Options:     options,
	}
	// api.CreateCommand(ds, []string{"764981968579461130"}, appCmd, dispatchCommand)
	api.CreateCommand(ds, []string{"797308956162392094"}, appCmd, dispatchCommand)
}

func dispatchCommand(ds *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	// As you can see, names of subcommands (nested, top-level)
	// and subcommand groups are provided through the arguments.
	if i.Member.Permissions&discordgo.PermissionManageRoles != discordgo.PermissionManageRoles {
		switch options[0].Name {
		case "add":
		default:
			errors.HandleError(ds, i, errors.NewErrorMissingPermission("manage roles"))
			return
		}
	}
	switch options[0].Name {
	case "update":
		commandUpdate.RunCommand(ds, i)
	case "add":
		commandAdd.RunCommand(ds, i)
	case "link":
		commandLink.RunCommand(ds, i)
	default:
		runCommand(ds, i)
	}
}

func runCommand(ds *discordgo.Session, i *discordgo.InteractionCreate) {
	err := ds.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   uint64(discordgo.MessageFlagsEphemeral),
			Content: "Please select a subcommand",
		},
	})
	if err != nil {
		fmt.Printf("Failed to send message\n%s", err)
	}
}

func updateEmbed(ds *discordgo.Session, invite *database.PotentialInvite) {
	var user, updater, inviter *discordgo.User
	var err error
	if user, err = invite.User(ds); err != nil {
		logger.Err().Printf("Failed to load user with id '%s: %v\n", invite.UserID, err)
	}
	if inviter, err = invite.Inviter(ds); err != nil {
		logger.Err().Printf("Failed to load user with id '%s: %v\n", invite.InviterID, err)
	}
	if updater, err = invite.Updater(ds); err != nil {
		logger.Err().Printf("Failed to load user with id '%s: %v\n", invite.UpdaterID, err)
	}
	embed := generateEmbed(user, updater, inviter, invite.InviteStatus())
	ds.ChannelMessageEditEmbed(config.InviteChannel, invite.InviteMessageID, embed)
}

func generateEmbed(user *discordgo.User, updater *discordgo.User, inviter *discordgo.User, status database.InviteStatus) *discordgo.MessageEmbed {
	fields := []string{
		fmt.Sprintf("**Mention** %s", user.Mention()),
		fmt.Sprintf("**User ID** `%s`", user.ID),
		fmt.Sprintf("**Status** %s", status.TermStatus(updater.Mention())),
		fmt.Sprintf("**Last Updated** <t:%d>", time.Now().Unix()),
	}
	embed := &discordgo.MessageEmbed{
		Title: user.Username + "#" + user.Discriminator,
		Color: status.Color,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    inviter.Username,
			IconURL: inviter.AvatarURL(""),
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: user.AvatarURL(""),
		},
		Description: strings.Join(fields, "\n"),
	}
	return embed
}

func getUsers(ds *discordgo.Session) []*discordgo.ApplicationCommandOptionChoice {
	users := make([]*database.PotentialInvite, 0)
	database.Get().Not(database.PotentialInvite{InviteStatusName: "Accepted"}).Not(database.PotentialInvite{InviteStatusName: "Declined"}).Find(&users)

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
	return userOptions
}
func (Command) Name() string {
	return "invite"
}
