// Copyright 2022 Bindl Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package program

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

	"github.com/bindl-dev/bindl/download"
	"github.com/bindl-dev/bindl/internal"
)

// Base is a minimal structure which exists in every program variations
type Base struct {
	Overlay map[string]map[string]string `json:"overlay,omitempty"`
	Name    string                       `json:"name"`
	Version string                       `json:"version"`
}

// Vars returns a map of variables to be used in templates,
// with overlays applied.
func (b *Base) Vars(goOS, goArch string) map[string]string {
	vars := map[string]string{
		"Name":    b.Name,
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

// Config is a program declaration in configuration file (default: bindl.yaml)
type Config struct {
	Base

	Checksums map[string]string `json:"checksums"`
	Provider  string            `json:"provider"`
	Path      string            `json:"path"`
}

// Lock converts current configuration to Lock, which is the format used by lockfile.
func (c *Config) Lock(ctx context.Context, platforms map[string][]string, useCache bool) (*Lock, error) {
	if err := c.loadChecksum(ctx, platforms, useCache); err != nil {
		return nil, fmt.Errorf("loading checksums: %w", err)
	}
	var p *Lock
	var err error
	switch c.Provider {
	case "url":
		p, err = NewLock(c)
		if err != nil {
			return nil, err
		}
		if err = p.collectBinaryChecksum(ctx, platforms, useCache); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknown program config provider: %s", c.Provider)
	}

	return p, nil
}

func (c *Config) loadChecksum(ctx context.Context, platforms map[string][]string, useCache bool) error {
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
			internal.Log().Debug().
				Str("program", c.Name).
				Str("platform", os+"/"+arch).
				Str("url", buf.String()).
				Msg("generate checksum url")
			checksumSrc[buf.String()] = struct{}{}
		}
	}

	checksums := map[string]string{}

	for url := range checksumSrc {
		d := &download.HTTP{UseCache: useCache}
		r, err := d.Get(ctx, url)
		if err != nil {
			return fmt.Errorf("retrieving checksums from '%s': %w", url, err)
		}
		data, err := readChecksumRef(r)
		d.Close()
		if err != nil {
			return fmt.Errorf("reading checksums: %w", err)
		}
		for f, cs := range data {
			internal.Log().Debug().Str("program", c.Name).Str(f, cs).Msg("retrieved checksum")
			checksums[f] = cs
		}
	}

	// Override the downloaded result with any explicitly specified checksum
	for f, cs := range c.Checksums {
		if f != "_src" {
			internal.Log().Warn().Str("program", c.Name).Str(f, cs).Msg("overwrite retrieved checksum")
		}
		checksums[f] = cs
	}
	c.Checksums = checksums

	return nil
}
