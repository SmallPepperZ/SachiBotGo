package database

import (
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/smallpepperz/sachibotgo/api/config"
)

type PotentialInvite struct {
	UserID          string `gorm:"primarykey"`
	InviterID       string
	UpdaterID       string
	InviteMessageID string
	InviteStatus    InviteStatus `gorm:"embedded;embeddedPrefix:status_"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (invite *PotentialInvite) User(ds *discordgo.Session) (user *discordgo.User, err error) {
	user, err = ds.User(invite.UserID)
	return
}

func (invite *PotentialInvite) Message(ds *discordgo.Session) (user *discordgo.Message, err error) {
	user, err = ds.ChannelMessage(config.InviteChannel, invite.InviteMessageID)
	return
}

func (invite *PotentialInvite) Inviter(ds *discordgo.Session) (user *discordgo.User, err error) {
	user, err = ds.User(invite.InviterID)
	return
}

func (invite *PotentialInvite) Updater(ds *discordgo.Session) (user *discordgo.User, err error) {
	user, err = ds.User(invite.UpdaterID)
	return
}

func (invite *PotentialInvite) Save() {
	db := Get()
	db.Save(invite)
}

type InviteStatus struct {
	Name               string
	Color              int
	TermStatusTemplate string
}

func (status InviteStatus) TermStatus(user string) string {
	return strings.Replace(status.TermStatusTemplate, "{user}", user, 1)
}

var InviteStatuses = struct {
	Invited  InviteStatus
	Active   InviteStatus
	Approved InviteStatus
	Accepted InviteStatus
	Rejected InviteStatus
	Declined InviteStatus
	Paused   InviteStatus
}{
	Invited: InviteStatus{
		"invited",
		0xFFFF00,
		"Invited",
	},
	Active: InviteStatus{
		"active",
		0xFFFF00,
		"Suggested by {user}",
	},
	Approved: InviteStatus{
		"approved",
		0x17820e,
		"Approved by {user}",
	},
	Accepted: InviteStatus{
		"accepted",
		0x1bc912,
		"Accepted invitation",
	},
	Rejected: InviteStatus{
		"rejected",
		0xa01116,
		"Rejected by {user}",
	},
	Declined: InviteStatus{
		"declined",
		0xd81d1a,
		"Declined invitation",
	},
	Paused: InviteStatus{
		"paused",
		0x444444,
		"Paused by {user}",
	},
}
