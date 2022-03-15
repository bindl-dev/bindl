package program

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"text/template"

	"go.xargs.dev/bindl/download"
	"go.xargs.dev/bindl/internal"
)

type Program interface {
	Name() string
	URL(goOS, goArch string) (string, error)
	DownloadArchive(ctx context.Context, d download.Downloader, goOS, goArch string) (*Archive, error)
}

type Base struct {
	PName   string                       `json:"name"`
	Version string                       `json:"version"`
	Overlay map[string]map[string]string `json:"overlay"`
}

func (b *Base) Name() string {
	return b.PName
}

func (b *Base) Vars(goOS, goArch string) map[string]string {
	vars := map[string]string{
		"Name":    b.PName,
		"Version": b.Version,

		"OS":   goOS,
		"Arch": goArch,
	}

	for label, originals := range b.Overlay {
		for original, replacement := range originals {
			if vars[label] == original {
				vars[label] = replacement
			}
		}
	}

	return vars
}

type Config struct {
	Base

	Provider  string            `json:"provider"`
	Path      string            `json:"path"`
	Checksums map[string]string `json:"checksums"`
}

func (c *Config) Program(platforms map[string][]string) (Program, error) {
	if err := c.loadChecksum(platforms); err != nil {
		return nil, fmt.Errorf("loading checksums: %w", err)
	}
	switch c.Provider {
	case "url":
		return NewURLProgram(c)
	default:
		return nil, fmt.Errorf("unknown program config provider: %s", c.Provider)
	}
}

func (c *Config) loadChecksum(platforms map[string][]string) error {
	src := c.Checksums["_src"]
	if src == "" {
		internal.Log().Debug().Msg("no checksum source was provided")
		return nil
	}

	t, err := template.New("url").Parse(src)
	if err != nil {
		return fmt.Errorf("parsing checksum url template: %w", err)
	}

	checksumSrc := map[string]struct{}{}

	for os, archs := range platforms {
		for _, arch := range archs {
			var buf bytes.Buffer
			if err := t.Execute(&buf, c.Vars(os, arch)); err != nil {
				return fmt.Errorf("retrieving checksum for %s/%s: %w", os, arch, err)
			}
			checksumSrc[buf.String()] = struct{}{}
		}
	}

	checksums := map[string]string{}

	for url, _ := range checksumSrc {
		resp, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("retrieving checksums from '%s': %w", url, err)
		}
		data, err := readChecksum(resp.Body)
		resp.Body.Close()
		if err != nil {
			return fmt.Errorf("reading checksums: %w", err)
		}
		for f, cs := range data {
			checksums[f] = cs
		}
	}

	// Override the downloaded result with any explicitly specified checksum
	for f, cs := range c.Checksums {
		checksums[f] = cs
	}
	c.Checksums = checksums

	return nil
}
