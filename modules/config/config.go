// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"sync"

	"code.gitea.io/tea/modules/utils"

	"github.com/adrg/xdg"
	"gopkg.in/yaml.v2"
)

// Preferences that are stored in and read from the config file
type Preferences struct {
	// Prefer using an external text editor over inline multiline prompts
	Editor bool `yaml:"editor"`
}

// LocalConfig represents local configurations
type LocalConfig struct {
	Logins []Login     `yaml:"logins"`
	Prefs  Preferences `yaml:"preferences"`
}

var (
	// config contain if loaded local tea config
	config         LocalConfig
	loadConfigOnce sync.Once
)

// GetConfigPath return path to tea config file
func GetConfigPath() string {
	configFilePath, err := xdg.ConfigFile("tea/config.yml")

	var exists bool
	if err != nil {
		exists = false
	} else {
		exists, _ = utils.PathExists(configFilePath)
	}

	// fallback to old config if no new one exists
	if !exists {
		file := filepath.Join(xdg.Home, ".tea", "tea.yml")
		exists, _ = utils.PathExists(file)
		if exists {
			return file
		}
	}

	if err != nil {
		log.Fatal("unable to get or create config file")
	}

	return configFilePath
}

// GetPreferences returns preferences based on the config file
func GetPreferences() Preferences {
	return config.Prefs
}

// loadConfig load config from file
func loadConfig() (err error) {
	loadConfigOnce.Do(func() {
		ymlPath := GetConfigPath()
		exist, _ := utils.FileExist(ymlPath)
		if exist {
			bs, err := ioutil.ReadFile(ymlPath)
			if err != nil {
				err = fmt.Errorf("Failed to read config file: %s", ymlPath)
			}

			err = yaml.Unmarshal(bs, &config)
			if err != nil {
				err = fmt.Errorf("Failed to parse contents of config file: %s", ymlPath)
			}
		}
	})
	return
}

// saveConfig save config to file
func saveConfig() error {
	ymlPath := GetConfigPath()
	bs, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(ymlPath, bs, 0660)
}
