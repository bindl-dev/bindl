package config

import "go.xargs.dev/bindl/program"

type Config struct {
	Output string `json:"output"`

	Platforms map[string][]string `json:"platform"`

	Programs []*program.Config `json:"programs"`
}
