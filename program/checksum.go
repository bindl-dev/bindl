package program

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
)

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
