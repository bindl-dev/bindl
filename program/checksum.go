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

func readChecksum(r io.Reader) (map[string]string, error) {
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

	return map[string]string{}, nil
}

func checksumSHA256(binary []byte, expect []byte) error {
	got := make([]byte, hex.EncodedLen(sha256.Size))
	byteHash := sha256.Sum256(binary)
	_ = hex.Encode(got, byteHash[:])
	if !bytes.Equal(expect, got) {
		return fmt.Errorf("checksum mismatch:\n\texpected: %s\n\treceived: %s", expect, got)
	}
	return nil
}

type checksumCalculator struct {
	h hash.Hash
}

func (c *checksumCalculator) SHA256(dst io.Writer) io.Writer {
	h := sha256.New()
	c.h = h

	w := io.MultiWriter(dst, h)
	return w
}

func (c *checksumCalculator) Error(expect []byte) error {
	got := make([]byte, hex.EncodedLen(sha256.Size))
	byteHash := c.h.Sum(nil)
	_ = hex.Encode(got, byteHash)
	if !bytes.Equal(expect, got) {
		return fmt.Errorf("checksum mismatch:\n\texpected: %s\n\treceived: %s", expect, got)
	}
	return nil
}