package config

import (
	"os"
	"time"

	"go.xargs.dev/bindl/program"
	"sigs.k8s.io/yaml"
)

type Lock struct {
	Updated  time.Time             `json:"updated"`
	Programs []*program.URLProgram `json:"programs"`
}

func ParseLock(path string) (*Lock, error) {
	l := &Lock{}
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(raw, l); err != nil {
		return nil, err
	}
	return l, nil
}
