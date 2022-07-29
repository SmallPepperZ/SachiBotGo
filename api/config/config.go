package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type DiscordConfigType struct {
	Token  string   `json:"token"`
	Guilds []string `json:"guilds"`
}

func (a *DiscordConfigType) Equal(b *DiscordConfigType) bool {
	equal := a.Token == b.Token
	if len(a.Guilds) != len(b.Guilds) {
		return false
	}
	for i, server := range a.Guilds {
		if server != b.Guilds[i] {
			return false
		}
	}
	return equal
}

type configType struct {
	Discord       DiscordConfigType `json:"discord"`
	DBPath        string            `json:"db_path"`
	InviteChannel string            `json:"invite_channel"`
}

var configData configType
var configPath string
var Discord DiscordConfigType
var DBPath string
var InviteChannel string

func init() {
	initEnv()
}

func initEnv() {
	configPath = os.Getenv("SACHIBOTGO_CONFIGPATH")
	createConfigFile(configPath)
	file, err := os.Open(configPath)
	if err != nil {
		fmt.Println("Cannot open config.json:", err)
		panic(err)
	}
	err = json.NewDecoder(file).Decode(&configData)
	if err != nil {
		fmt.Println("Cannot decode config.json:", err)
		panic(err)
	}
	Discord = configData.Discord
	DBPath = configData.DBPath
	InviteChannel = configData.InviteChannel
}
func (c *configType) Save() error {
	file, err := os.OpenFile(configPath, os.O_WRONLY, os.ModeAppend)
	file.Truncate(0)
	if err != nil {
		return fmt.Errorf("cannot open config.json: %v", err)
	}

	enc := json.NewEncoder(file)
	enc.SetIndent("", "    ")
	err = enc.Encode(c)

	if err != nil {
		err = fmt.Errorf("cannot write to config.json: %v", err)
	}
	return err
}

func createConfigFile(path string) (err error) {
	_, err = os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		file, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("cannot create config.json: %w", err)
		}
		enc := json.NewEncoder(file)
		enc.SetIndent("", "    ")
		err = enc.Encode(configType{
			Discord: DiscordConfigType{
				Token:  "Bot TOKEN",
				Guilds: []string{"123", "456"},
			},
			DBPath: "/path/to/sqlite.db",
		})
		if err != nil {
			return fmt.Errorf("cannot write default config.json: %w", err)
		}
	}
	return
}
