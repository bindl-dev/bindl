package program

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"path/filepath"
	"text/template"

	"go.xargs.dev/bindl/download"
)

type URLProgram struct {
	PName          string                       `json:"name"`
	Version        string                       `json:"version"`
	URLTemplate    string                       `json:"url"`
	Overlay        map[string]map[string]string `json:"overlay,omitempty"`
	ChecksumSource string                       `json:"checksum,omitempty"`
	Checksums      map[string]*ArchiveChecksum  `json:"checksums,omitempty"`
}

func (p *URLProgram) Name() string {
	return p.PName
}

func (p *URLProgram) URL(goOS, goArch string) (string, error) {
	val := map[string]string{
		"Name":    p.PName,
		"Version": p.Version,

		"OS":   goOS,
		"Arch": goArch,
	}

	for label, originals := range p.Overlay {
		for original, replacement := range originals {
			if val[label] == original {
				val[label] = replacement
			}
		}
	}

	t, err := template.New("url").Parse(p.URLTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, val); err != nil {
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

	_, err = io.Copy(w, body)
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
