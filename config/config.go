package config

import "go.xargs.dev/bindl/program"

type Config struct {
	Platforms map[string][]string `json:"platforms"`

	Programs []*program.Config `json:"programs"`
}
