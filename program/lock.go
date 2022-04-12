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

	"github.com/bindl-dev/bindl/download"
	"github.com/bindl-dev/bindl/internal"
)

// Lock is a configuration used by lockfile to explicitly state the
// expected validations of each program.
type Lock struct {
	Base

	Checksums map[string]*ArchiveChecksum `json:"checksums,omitempty"`
	Paths     *RemotePath                 `json:"paths"`
	Cosign    []*CosignBundle             `json:"cosign,omitempty"`
}

func NewLock(ctx context.Context, c *Config, platforms map[string][]string, useCache bool) (*Lock, error) {
	p := &Lock{
		Base: Base{
			Name:    c.Name,
			Version: c.Version,
			Overlay: c.Overlay,
		},
		Checksums: map[string]*ArchiveChecksum{},
		Paths:     c.Paths,
	}
	if c.Paths.Cosign != nil {
		p.Cosign = c.Paths.Cosign
	}
	for f, cs := range c.Checksums {
		p.Checksums[f] = &ArchiveChecksum{Archive: cs}
	}
	if err := p.collectBinaryChecksum(ctx, platforms, useCache); err != nil {
		return nil, err
	}
	return p, nil
}

// TOFU: Trust on first use -- should only be run first time a program was added to
// the lockfile. Collecting binary checksums by extracting archives.
// TODO: Use the values presented in SBOM when available.
func (p *Lock) collectBinaryChecksum(ctx context.Context, platforms map[string][]string, useCache bool) error {
	var wg sync.WaitGroup

	hasError := false

	for os, archs := range platforms {
		for _, arch := range archs {
			wg.Add(1)
			go func(os, arch string) {
				defer wg.Done()

				a, err := p.DownloadArchive(ctx, &download.HTTP{UseCache: useCache}, os, arch)
				if err != nil {
					internal.ErrorMsg(fmt.Errorf("downloading archive for '%s' in %s/%s: %w", p.Name, os, arch, err))
					return
				}

				b, err := a.BinaryChecksum(p.Name)
				if err != nil {
					internal.ErrorMsg(fmt.Errorf("calculating binary checksum for '%s' in %s/%s: %w", p.Name, os, arch, err))
					return
				}
				p.Checksums[a.Name].Binary = string(b)
			}(os, arch)
		}
	}

	wg.Wait()
	if hasError {
		return fmt.Errorf("failed to collect all checksums")
	}

	// Not all checksums was necessarily used, remove the ones not caught by platform matrix
	toDelete := []string{}
	for archiveName, archiveCS := range p.Checksums {
		if archiveCS.Binary == "" {
			toDelete = append(toDelete, archiveName)
		}
	}
	for _, archiveName := range toDelete {
		delete(p.Checksums, archiveName)
	}

	return nil
}

// ArchiveName returns the archive name with OS and Arch interpolated
// if necessary, i.e. someprogram-linux-amd64.tar.gz.
// This reads from URL and assumes that contains the archive name.
func (p *Lock) ArchiveName(os, arch string) (string, error) {
	url, err := p.URL(os, arch)
	if err != nil {
		return "", err
	}

	return filepath.Base(url), nil
}

// URL returns the download URL with variables interpolated as necessary.
func (p *Lock) URL(goOS, goArch string) (string, error) {
	t, err := template.New("url").Parse(p.Paths.target())
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, p.Vars(goOS, goArch)); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// DownloadArchive returns Archive which has the archive data in-memory, with guarantees
// on archive checksum. That is, if checksum fails, no data will be made available to caller.
func (p *Lock) DownloadArchive(ctx context.Context, d download.Downloader, goOS, goArch string) (*Archive, error) {
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

	internal.Log().Info().Str("program", p.Name).Str("platform", goOS+"/"+goArch).Msg("downloading archive")
	body, err := d.Get(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("downloading program: %w", err)
	}

	c := &ChecksumCalculator{}
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
