package main

import "github.com/BurntSushi/toml"

type Config struct {
	Users    map[string]string `toml:"users"`
	GitHub   GithubConfig      `toml:"github"`
	Slack    SlackConfig       `toml:"slack"`
	Unassign UnassignConfig    `toml:"unassign"`
}

type GithubConfig struct {
	Organization   string   `toml:"organization"`
	Repos          []string `toml:"repo"`
	PrivateKey     string   `toml:"privateKey"`
	AppID          int64    `toml:"appId"`
	InstallationID int64    `toml:"installationId"`
	IgnoreLabels   []string `toml:"ignoreLabels"`
}

type SlackConfig struct {
	Token         string `toml:"token"`
	ErrlogChannel string `toml:"errlogChannel"`
}

type UnassignConfig struct {
	DaysUntilWarning  int `toml:"daysUntilWarning"`
	DaysUntilUnassign int `toml:"daysUntilUnassign"`
}

func configure(config *Config) {
	_, err := toml.DecodeFile("./config.toml", &config)
	if err != nil {
		panic(err)
	}
}
