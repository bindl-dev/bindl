package program

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"path/filepath"
	"text/template"

	"go.xargs.dev/bindl/download"
	"go.xargs.dev/bindl/internal"
)

type URLProgram struct {
	Base

	URLTemplate string                      `json:"url"`
	Checksums   map[string]*ArchiveChecksum `json:"checksums,omitempty"`
}

func NewURLProgram(c *Config) (*URLProgram, error) {
	checksums := map[string]*ArchiveChecksum{}
	for f, cs := range c.Checksums {
		checksums[f] = &ArchiveChecksum{Archive: cs}
	}
	p := &URLProgram{
		Base: Base{
			PName:   c.PName,
			Version: c.Version,
			Overlay: c.Overlay,
		},
		URLTemplate: c.Path,
		Checksums:   checksums,
	}
	return p, nil
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
