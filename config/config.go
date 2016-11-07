package config

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
)

type remoteConfig struct {
	Name               string
	Alias              string
	Username           string
	Host               string
	SshParams          string
	McDefaultDirLocal  string
	McDefaultDirRemote string
}

type groupConfig struct {
	Name    string
	Alias   string
	AliasMc string
	Remotes []remoteConfig
}

type config struct {
	Groups []groupConfig
}

func Config() *config {
	c := new(config)
	c.ReloadConfig()
	if c.Groups == nil {
		c.Groups = make([]groupConfig, 0)
	}
	return c
}

func (c config) ReloadConfig() {
	file, err := os.Open(getConfigPath())
	if err == nil {
		decoder := json.NewDecoder(file)
		err := decoder.Decode(&c)
		if err != nil {
			fmt.Println("error:", err)
		}
	}
}

func (c config) SaveConfig() {
	file, _ := os.Create(getConfigPath())
	file.Chmod(0600)
	encoder := json.NewEncoder(file)
	err := encoder.Encode(c)
	if err != nil {
		fmt.Println("error:", err)
	}
}

func getConfigPath() string {
	usr, err := user.Current()
	if err != nil {
		fmt.Println("error:", err)
	}

	os.Mkdir(usr.HomeDir+"/.remote-manager", 0700)

	return usr.HomeDir + "/.remote-manager/conf.json"
}
