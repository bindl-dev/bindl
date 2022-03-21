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
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"strings"
)

func readChecksumRef(r io.Reader) (map[string]string, error) {
	raw, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	result := map[string]string{}
	str := strings.ReplaceAll(string(raw), "\r\n", "\n")
	for _, line := range strings.Split(str, "\n") {
		if line == "" {
			continue
		}
		cs, f, ok := strings.Cut(line, "  ")
		if !ok {
			return nil, fmt.Errorf("failed to cut string: unexpected formatting: '%s'", line)
		}
		result[f] = cs
	}
	return result, nil
}

func checksumSHA256(binary []byte) []byte {
	v := make([]byte, hex.EncodedLen(sha256.Size))
	byteHash := sha256.Sum256(binary)
	_ = hex.Encode(v, byteHash[:])
	return v
}

func assertChecksumSHA256(binary []byte, expect []byte) error {
	got := checksumSHA256(binary)
	if !bytes.Equal(expect, got) {
		return fmt.Errorf("checksum mismatch:\n\texpected: %s\n\treceived: %s", expect, got)
	}
	return nil
}

type ChecksumCalculator struct {
	h hash.Hash
}

func (c *ChecksumCalculator) SHA256(dst io.Writer) io.Writer {
	h := sha256.New()
	c.h = h

	w := io.MultiWriter(dst, h)
	return w
}

// Error compares the Base64 encoded checksum
func (c *ChecksumCalculator) Error(expect []byte) error {
	got := make([]byte, hex.EncodedLen(sha256.Size))
	byteHash := c.h.Sum(nil)
	_ = hex.Encode(got, byteHash)
	if !bytes.Equal(expect, got) {
		return fmt.Errorf("checksum mismatch:\n\texpected: %s\n\treceived: %s", expect, got)
	}
	return nil
}
