package commands

import (
	"github.com/smallpepperz/sachibotgo/api"
	"github.com/smallpepperz/sachibotgo/commands/invite"
	"github.com/smallpepperz/sachibotgo/commands/utility/list"
	"github.com/smallpepperz/sachibotgo/commands/utility/ping"
)

func AddAllCommands() {
	api.AddCommand(&list.Command{})
	api.AddCommand(&ping.Command{})
	api.AddCommand(&invite.Command{})
}
