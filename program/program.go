package program

import (
	"context"

	"go.xargs.dev/bindl/download"
)

type Program interface {
	Name() string
	URL(goOS, goArch string) (string, error)
	DownloadArchive(ctx context.Context, d download.Downloader, goOS, goArch string) (*Archive, error)
}

type Config struct {
	Name     string                       `json:"name"`
	Provider string                       `json:"provider"`
	Path     string                       `json:"path"`
	Version  string                       `json:"version"`
	Checksum string                       `json:"checksum"`
	Overlay  map[string]map[string]string `json:"overlay"`
}

//func (c *Config) Program() (Program, error) {
//	switch c.Provider{
//	case "url":
//
//	}
//}
