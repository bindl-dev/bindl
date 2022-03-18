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
	"io"
	"path/filepath"
	"sync"
	"text/template"

	"go.xargs.dev/bindl/download"
	"go.xargs.dev/bindl/internal"
)

type URLProgram struct {
	Base

	URLTemplate string                      `json:"url"`
	ChecksumSrc string                      `json:"checksum,omitempty"`
	Checksums   map[string]*ArchiveChecksum `json:"checksums,omitempty"`
}

func NewURLProgram(c *Config) (*URLProgram, error) {
	p := &URLProgram{
		Base: Base{
			PName:   c.PName,
			Version: c.Version,
			Overlay: c.Overlay,
		},
		URLTemplate: c.Path,
		Checksums:   map[string]*ArchiveChecksum{},
	}
	for f, cs := range c.Checksums {
		if f == "_src" {
			p.ChecksumSrc = cs
		} else {
			p.Checksums[f] = &ArchiveChecksum{Archive: cs, Binaries: map[string]string{}}
		}
	}
	return p, nil
}

// TOFU: Trust on first use -- should only be run first time a program was added to
// the lockfile. Collecting binary checksums by extracting archives.
func (p *URLProgram) collectBinaryChecksum(ctx context.Context, platforms map[string][]string) error {
	var wg sync.WaitGroup

	hasError := false

	for os, archs := range platforms {
		for _, arch := range archs {
			wg.Add(1)
			go func(os, arch string) {
				defer wg.Done()

				a, err := p.DownloadArchive(ctx, &download.HTTP{}, os, arch)
				if err != nil {
					internal.ErrorMsg(fmt.Errorf("downloading archive for '%s' in %s/%s: %w", p.PName, os, arch, err))
					return
				}

				b, err := a.BinaryChecksum(p.PName)
				if err != nil {
					internal.ErrorMsg(fmt.Errorf("calculating binary checksum for '%s' in %s/%s: %w", p.PName, os, arch, err))
					return
				}
				p.Checksums[a.Name].Binaries[p.PName] = string(b)
			}(os, arch)
		}
	}

	wg.Wait()
	if hasError {
		return fmt.Errorf("failed to collect all checksums")
	}
	return nil
}

func (p *URLProgram) URL(goOS, goArch string) (string, error) {
	t, err := template.New("url").Parse(p.URLTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, p.Vars(goOS, goArch)); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (p *URLProgram) DownloadArchive(ctx context.Context, d download.Downloader, goOS, goArch string) (*Archive, error) {
	url, err := p.URL(goOS, goArch)
	if err != nil {
		return nil, fmt.Errorf("generating URL for download: %w", err)
	}
	a := &Archive{Name: filepath.Base(url)}
	if a.Name == "" || a.Name == "." {
		return nil, fmt.Errorf("failed to determine archive name from url: %s", url)
	}

	if p.Checksums[a.Name] == nil {
		return nil, fmt.Errorf("no checksum reference found for '%s': aborting", a.Name)
	}
	expHash := p.Checksums[a.Name].Archive
	if expHash == "" {
		return nil, fmt.Errorf("expected hash for '%s' is invalid: cannot be empty string", a.Name)
	}

	body, err := d.Get(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("downloading program: %w", err)
	}

	c := &checksumCalculator{}
	var buf bytes.Buffer
	w := c.SHA256(&buf)

	n, err := io.Copy(w, body)
	internal.Log().Debug().Err(err).Int64("bytes", n).Msg("copying response body")
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}
	d.Close()

	if err := c.Error([]byte(expHash)); err != nil {
		return nil, fmt.Errorf("validating checksum for archive '%s': %w", a.Name, err)
	}

	a.Data = buf.Bytes()
	a.checksums = p.Checksums[a.Name]
	return a, nil
}
