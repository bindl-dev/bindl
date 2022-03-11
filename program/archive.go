package program

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"
)

type Archive struct {
	Name string
	Data []byte

	checksums *ArchiveChecksum
}

const (
	archiveTarSuffix   = ".tar"
	archiveTarGzSuffix = ".tar.gz"
	// TODO: archiveZipSuffix   = ".zip"
)

func (a *Archive) Extract(filename string) ([]byte, error) {
	binary, err := a.extractBinaryNoChecksum(filename)
	if err != nil {
		return nil, fmt.Errorf("extracting binary source from archive: %w", err)
	}

	expect := a.checksums.Binaries[filename]

	if err := checksumSHA256(binary, []byte(expect)); err != nil {
		return nil, fmt.Errorf("checksum validation for '%s': %w", filename, err)
	}

	return binary, nil
}

func (a *Archive) extractBinaryNoChecksum(filename string) ([]byte, error) {
	var buf bytes.Buffer
	var data []byte
	var err error

	r := bytes.NewBuffer(a.Data)

	switch {
	case strings.HasSuffix(a.Name, archiveTarSuffix):
		err = untar(&buf, r, filename)
	case strings.HasSuffix(a.Name, archiveTarGzSuffix):
		err = untargz(&buf, r, filename)
	default:
		// assume currently downloaded file is the binary
		data = a.Data
	}

	if err != nil {
		return nil, err
	}
	if data == nil {
		data = buf.Bytes()
	}
	return data, err
}

func untargz(w io.Writer, rawTarGz io.Reader, filename string) error {
	gzReader, err := gzip.NewReader(rawTarGz)
	if err != nil {
		return err
	}
	defer gzReader.Close()

	return untar(w, gzReader, filename)
}

func untar(w io.Writer, rawTar io.Reader, filename string) error {
	tarReader := tar.NewReader(rawTar)

	header, err := tarReader.Next()
	for {
		if err != nil {
			break
		}
		if header.Typeflag != tar.TypeReg || filepath.Base(header.Name) != filename {
			header, err = tarReader.Next()
			continue
		}

		_, err = io.Copy(w, tarReader)
		break
	}

	if errors.Is(err, io.EOF) {
		err = fmt.Errorf("unable to find '%s' in archive: %w", filename, err)
	}
	return err
}

const archiveChecksumKey = "_archive"

type ArchiveChecksum struct {
	Archive  string
	Binaries map[string]string
}

func (c *ArchiveChecksum) MarshalJSON() ([]byte, error) {
	raw := map[string]string{}

	raw[archiveChecksumKey] = c.Archive
	for b, cs := range c.Binaries {
		raw[b] = cs
	}

	return json.Marshal(raw)
}

func (c *ArchiveChecksum) UnmarshalJSON(b []byte) error {
	if c.Binaries == nil {
		c.Binaries = map[string]string{}
	}

	raw := map[string]string{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return fmt.Errorf("reading partial json: %w", err)
	}

	for name, hash := range raw {
		if name == archiveChecksumKey {
			c.Archive = hash
		} else {
			c.Binaries[name] = hash
		}
	}

	return nil
}
