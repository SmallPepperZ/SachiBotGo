package api

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestMain(m *testing.M) {
	m.Run()
}
func TestSave(t *testing.T) {
	// Avoid overriding the actual config file
	t.Setenv("SACHIBOTGO_CONFIGPATH", path.Join(t.TempDir(), "config.json"))
	initEnv()

	compareConfig := &ConfigType{}
	testConfig := &ConfigType{
		Discord: DiscordConfigType{
			Token:  "abc",
			Guilds: []string{"def", "xyz"},
		},
	}
	testConfig.Save()

	// Read the file
	file, err := os.Open(configPath)
	if err != nil {
		t.Error("Cannot open config.json:", err)
	}

	// Get the file contents
	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		t.Error("Cannot read config.json:", err)
	}

	// Unmarshal the file contents
	err = json.Unmarshal(fileContents, compareConfig)
	if err != nil {
		t.Log("Contents:", string(fileContents))
		t.Error("Cannot unmarshal config.json:", err)
	}

	// Ensure the file matches the expected contents
	if !compareConfig.Equal(testConfig) {
		t.Error("Config not saved")
	}
}
