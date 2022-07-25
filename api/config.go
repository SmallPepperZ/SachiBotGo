package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
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

type ConfigType struct {
	Discord DiscordConfigType `json:"discord"`
}

func (a *ConfigType) Equal(b *ConfigType) bool {
	return a.Discord.Equal(&b.Discord)
}

var Config ConfigType
var configPath string

func init() {
	initEnv()
}

func initEnv() {
	configPath = os.Getenv("SACHIBOTGO_CONFIGPATH")
	createConfigFile(configPath)
	err := cleanenv.ReadConfig(configPath, &Config)
	if err != nil {
		fmt.Println("Cannot open config.json:", err)
		panic(err)
	}
}
func (c *ConfigType) Save() error {
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

func createConfigFile(path string) (file *os.File, err error) {
	_, err = os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		file, err = os.Create(path)
		if err != nil {
			return file, fmt.Errorf("cannot create config.json: %w", err)
		}
	}
	err = json.NewEncoder(file).Encode(nil)
	return
}
