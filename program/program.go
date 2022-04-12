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
	"encoding/base64"
	"fmt"
	"io"
	"strings"
	"sync"
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

	Paths     *RemotePath       `json:"paths"`
	Checksums map[string]string `json:"checksums"`
	Provider  string            `json:"provider"`
}

// Lock converts current configuration to Lock, which is the format used by lockfile.
func (c *Config) Lock(ctx context.Context, platforms map[string][]string, useCache bool) (*Lock, error) {
	switch c.Provider {
	case "github":
		if err := githubToURL(c); err != nil {
			return nil, err
		}
	}

	if err := c.loadChecksum(ctx, platforms, useCache); err != nil {
		return nil, fmt.Errorf("loading checksums: %w", err)
	}

	p, err := NewLock(ctx, c, platforms, useCache)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (c *Config) loadChecksum(ctx context.Context, platforms map[string][]string, useCache bool) error {
	c.Paths.fmt()
	if c.Paths.Checksums == nil {
		internal.Log().Debug().Str("program", c.Name).Msg("no remote checksum path provided")
		return nil
	}

	checksumSrcs, err := c.Paths.Checksums.generate(platforms, c.Vars)
	if err != nil {
		return fmt.Errorf("generating checksum source(s): %w", err)
	}

	rawChecksums := ""

	for _, checksumSrc := range checksumSrcs {
		bundle, err := checksumSrc.values(ctx, useCache)
		if err != nil {
			return fmt.Errorf("extracting checksum values from '%s': %w", checksumSrc.Artifact, err)
		}
		if bundle.Signed() {
			internal.Log().Debug().Str("program", c.Name).Msg("found signature specification")
			if err := bundle.VerifySignature(ctx); err != nil {
				return fmt.Errorf("verifying signed checksum: %w", err)
			}
			internal.Log().Info().Str("program", c.Name).Msg("cosign signature valid")
			c.Paths.Cosign = append(c.Paths.Cosign, bundle)
		}
		rawChecksums += "\n" + bundle.Artifact
	}

	checksums := map[string]string{}

	data, err := parseChecksumRef(rawChecksums)
	if err != nil {
		return fmt.Errorf("reading checksums: %w", err)
	}
	for f, cs := range data {
		internal.Log().Debug().Str("program", c.Name).Str(f, cs).Msg("found checksum")
		checksums[f] = cs
	}

	// Override the downloaded result with any explicitly specified checksum
	for f, cs := range c.Checksums {
		internal.Log().Warn().Str("program", c.Name).Str(f, cs).Msg("overwrite retrieved checksum")
		checksums[f] = cs
	}
	c.Checksums = checksums

	return nil
}

type RemotePath struct {
	Base   string `json:"base"`
	Target string `json:"target"`

	Checksums *ChecksumPaths `json:"checksums,omitempty"`

	Cosign []*CosignBundle `json:"-"`

	fmtOnce sync.Once
}

func (p *RemotePath) fmt() {
	p.fmtOnce.Do(func() {
		if l := len(p.Base); p.Base[l-1] != '/' {
			p.Base += "/"
		}
		if p.Checksums != nil {
			p.Checksums.rebase(p.Base)
		}
	})
}

func (p *RemotePath) target() string {
	p.fmt()
	return p.Base + p.Target
}

type ChecksumPaths struct {
	Artifact    string `json:"artifact"`
	Certificate string `json:"certificate,omitempty"`
	Signature   string `json:"signature,omitempty"`
}

func (c *ChecksumPaths) rebase(base string) {
	if !strings.HasPrefix(c.Artifact, "http") {
		c.Artifact = base + c.Artifact
	}
	if c.Certificate != "" && !strings.HasPrefix(c.Certificate, "http") {
		c.Certificate = base + c.Certificate
	}
	if c.Signature != "" && !strings.HasPrefix(c.Signature, "http") {
		c.Signature = base + c.Signature
	}
}

func (c *ChecksumPaths) generate(platforms map[string][]string, varsFn func(os, arch string) map[string]string) ([]*ChecksumPaths, error) {
	tArtifact, err := template.New("url").Parse(c.Artifact)
	if err != nil {
		return nil, err
	}
	tCert, err := template.New("url").Parse(c.Certificate)
	if err != nil {
		return nil, err
	}
	tSig, err := template.New("url").Parse(c.Signature)
	if err != nil {
		return nil, err
	}

	// Use map to ensure uniqueness
	bundles := map[string]*ChecksumPaths{}
	for os, archs := range platforms {
		for _, arch := range archs {
			var artifact, cert, sig bytes.Buffer
			if err := tArtifact.Execute(&artifact, varsFn(os, arch)); err != nil {
				return nil, fmt.Errorf("retrieving checksum artifact for %s/%s: %w", os, arch, err)
			}
			if err := tCert.Execute(&cert, varsFn(os, arch)); err != nil {
				return nil, fmt.Errorf("retrieving checksum certificate for %s/%s: %w", os, arch, err)
			}
			if err := tSig.Execute(&sig, varsFn(os, arch)); err != nil {
				return nil, fmt.Errorf("retrieving checksum certificate for %s/%s: %w", os, arch, err)
			}
			bundles[artifact.String()] = &ChecksumPaths{
				Artifact:    artifact.String(),
				Certificate: cert.String(),
				Signature:   sig.String(),
			}
		}
	}

	result := []*ChecksumPaths{}
	for _, bundle := range bundles {
		result = append(result, bundle)
	}
	return result, nil
}

func (c *ChecksumPaths) values(ctx context.Context, useCache bool) (*CosignBundle, error) {
	download := func(dst *string, src string) error {
		if src == "" {
			return nil
		}
		d := &download.HTTP{UseCache: useCache}
		r, err := d.Get(ctx, src)
		defer d.Close()
		if err != nil {
			return err
		}
		v, err := io.ReadAll(r)
		if err != nil {
			return err
		}
		*dst = string(v)
		return nil
	}
	b := CosignBundle{}
	if err := download(&b.Artifact, c.Artifact); err != nil {
		return nil, err
	}
	if err := download(&b.Certificate, c.Certificate); err != nil {
		return nil, err
	}
	if err := download(&b.Signature, c.Signature); err != nil {
		return nil, err
	}

	rawCert, err := base64.StdEncoding.DecodeString(b.Certificate)
	if err != nil {
		return nil, fmt.Errorf("base64 decoding certificate: %w", err)
	}
	b.Certificate = string(rawCert)

	return &b, nil
}
