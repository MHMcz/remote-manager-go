/*
	Copyright © 2016 Jan Markup <mhmcze@gmail.com>
	This work is free. You can redistribute it and/or modify it under the
	terms of the Do What The Fuck You Want To Public License, Version 2,
	as published by Sam Hocevar. See the COPYING file for more details.
*/

package config

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
)

type RemoteConfig struct {
	Name               string
	Alias              string
	Username           string
	Host               string
	SshParams          string
	McDefaultDirLocal  string
	McDefaultDirRemote string
}

type GroupConfig struct {
	Name    string
	Alias   string
	AliasMc string
	Remotes []RemoteConfig
}

type Configuration struct {
	Groups []GroupConfig
}

func Config() *Configuration {
	c := new(Configuration)
	c.ReloadConfig()
	if c.Groups == nil {
		c.Groups = make([]GroupConfig, 0)
	}
	return c
}

func (c *Configuration) ReloadConfig() {
	file, err := os.Open(getConfigPath())
	if err == nil {
		decoder := json.NewDecoder(file)
		err := decoder.Decode(&c)
		if err != nil {
			fmt.Println("error:", err)
		}
	}
}

func (c *Configuration) SaveConfig() {
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
