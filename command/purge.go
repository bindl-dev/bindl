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
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/bindl-dev/bindl/config"
	"github.com/bindl-dev/bindl/internal"
)

// Purge deletes downloaded binaries. By default, it only deletes
// downloaded binaries which no longer mentioned by the lockfile.
// Passing `all` will ignore lockfile check and deletes all programs.
func Purge(ctx context.Context, conf *config.Runtime, all, dryRun bool) error {
	progDir := filepath.Join(conf.BinDir, conf.ProgDir)
	if all {
		return removeAll(progDir, dryRun)
	}

	l, err := config.ParseLock(conf.LockfilePath)
	if err != nil {
		return fmt.Errorf("parsing lockfile: %w", err)
	}

	keep := map[string]bool{}
	for _, p := range l.Programs {
		archiveName, err := p.ArchiveName(conf.OS, conf.Arch)
		if err != nil {
			return fmt.Errorf("generating archive name for '%s': %w", p.Name, err)
		}
		checksum := p.Checksums[archiveName].Binaries[p.Name]
		keepPath := checksum + "-" + p.Name
		internal.Log().Debug().Str("program", keepPath).Msg("to keep")
		keep[keepPath] = true
	}

	return filepath.WalkDir(progDir, func(path string, _ fs.DirEntry, err error) error {
		basePath := filepath.Base(path)
		if len(basePath) < hex.EncodedLen(sha256.Size)+1 {
			internal.Log().Debug().Str("path", path).Msg("not program directory, skip")
			return nil
		}
		if keep[basePath] {
			internal.Log().Debug().Str("path", path).Str("base", basePath).Msg("keep, skip")
			return fs.SkipDir
		}
		if err := removeAll(path, dryRun); err != nil {
			return err
		}
		return fs.SkipDir
	})
}

func removeAll(path string, dryRun bool) error {
	internal.Log().Info().Str("path", path).Msg("purge")
	if !dryRun {
		return os.RemoveAll(path)
	}
	return nil
}
