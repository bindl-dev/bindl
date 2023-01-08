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
	"fmt"
	"strings"
)

type Archive struct {
	checksums *ArchiveChecksum

	Name string
	Data []byte
}

const (
	archiveTarSuffix   = ".tar"
	archiveTarGzSuffix = ".tar.gz"
	archiveGzSuffix    = ".gz"
	archiveZipSuffix   = ".zip"
)

// BinaryChecksum retrieves checksum value of binary in archive.
// Returned value is base64-encoded []byte.
// TODO: maybe we should just return string(checksumSHA256(binary))?
func (a *Archive) BinaryChecksum(binaryName string) ([]byte, error) {
	binary, err := a.extractBinaryNoChecksum(binaryName)
	if err != nil {
		return nil, err
	}
	return checksumSHA256(binary), nil
}

// Extract returns the binary data from archive with checksum guarantee.
// That is, if checksum fails, then the binary will not be returned.
func (a *Archive) Extract(binaryName string) ([]byte, error) {
	binary, err := a.extractBinaryNoChecksum(binaryName)
	if err != nil {
		return nil, fmt.Errorf("extracting binary source from archive: %w", err)
	}

	expect := a.checksums.Binary

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
	case strings.HasSuffix(a.Name, archiveGzSuffix):
		err = ungz(&buf, r, binaryName)
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

type ArchiveChecksum struct {
	Archive string `json:"archive"`
	Binary  string `json:"binary"`
}
