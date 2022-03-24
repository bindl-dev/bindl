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

package command

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"go.xargs.dev/bindl/config"
	"go.xargs.dev/bindl/internal"
)

func isNewline(c byte) bool {
	return c == '\n' || c == '\r'
}

func getNumNewlinesRequired(f *os.File) (int, error) {
	// Empty file?
	off, err := f.Seek(0, os.SEEK_END)
	if off == 0 || err == io.EOF {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}

	// Check last character.
	buf := make([]byte, 1)
	if _, err = f.Seek(-1, os.SEEK_END); err != nil {
		return 0, err
	}
	if _, err = f.Read(buf); err != nil {
		return 0, err
	}
	if !isNewline(buf[0]) {
		return 2, nil
	}

	// Check 2nd-to-last character.
	if _, err = f.Seek(-2, os.SEEK_END); err != nil {
		return 0, err
	}
	if _, err = f.Read(buf); err != nil {
		return 0, err
	}
	if !isNewline(buf[0]) {
		return 1, nil
	}

	// Already have 2 newlines. No need to format.
	return 0, nil
}

func getValidTargets(binDir string) map[string]bool {
	noTrailingSlash := strings.TrimSuffix(binDir, "/")
	return map[string]bool{
		filepath.Join(binDir, "*"): true,
		noTrailingSlash + "/":      true,
		noTrailingSlash:            true,
	}
}

func getIgnoreEntry(binDir string, numPrefixNewlines int) string {
	prefix := ""
	for i := 0; i < numPrefixNewlines; i++ {
		prefix += "\n"
	}

	internal.Log().Debug().Int("newlines", numPrefixNewlines).Msg("creating entry to add to file")
	return prefix + "# Development and tool binaries\n" + filepath.Join(binDir, "*") + "\n"
}

func UpdateIgnoreFile(conf *config.Runtime, path string) error {
	internal.Log().Debug().Str("ignore", path).Msg("attempting to update ignore file")

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("opening file '%s': %w", path, err)
	}
	defer f.Close()

	targets := getValidTargets(conf.BinDir)

	n := 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		n++
		l := strings.TrimSpace(scanner.Text())
		internal.Log().Debug().Int("line", n).Str("ignore", l).Msg("parsing ignore file content")
		if targets[l] {
			internal.Log().Info().Int("line", n).Msg("found a qualifying ignore file entry, not modifying")
			return nil
		}
	}

	internal.Log().Debug().Msg("no qualifying ignore file entry found, adding")
	numNewlinesRequired, err := getNumNewlinesRequired(f)
	if err != nil {
		return err
	}

	_, err = f.WriteString(getIgnoreEntry(conf.BinDir, numNewlinesRequired))
	return err
}
