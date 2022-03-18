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
	"encoding/json"
	"fmt"
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
	archiveZipSuffix   = ".zip"
)

func (a *Archive) BinaryChecksum(binaryName string) ([]byte, error) {
	binary, err := a.extractBinaryNoChecksum(binaryName)
	if err != nil {
		return nil, err
	}
	return checksumSHA256(binary), nil
}

func (a *Archive) Extract(binaryName string) ([]byte, error) {
	binary, err := a.extractBinaryNoChecksum(binaryName)
	if err != nil {
		return nil, fmt.Errorf("extracting binary source from archive: %w", err)
	}

	expect := a.checksums.Binaries[binaryName]

	if err := assertChecksumSHA256(binary, []byte(expect)); err != nil {
		return nil, fmt.Errorf("checksum validation for '%s': %w", binaryName, err)
	}

	return binary, nil
}

func (a *Archive) extractBinaryNoChecksum(binaryName string) ([]byte, error) {
	var buf bytes.Buffer
	var data []byte
	var err error

	r := bytes.NewReader(a.Data)

	switch {
	case strings.HasSuffix(a.Name, archiveZipSuffix):
		err = unzip(&buf, r, int64(len(a.Data)), binaryName)
	case strings.HasSuffix(a.Name, archiveTarSuffix):
		err = untar(&buf, r, binaryName)
	case strings.HasSuffix(a.Name, archiveTarGzSuffix):
		err = untargz(&buf, r, binaryName)
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
