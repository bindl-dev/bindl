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
	"os"
	"path/filepath"
	"time"

	"github.com/bindl-dev/bindl/config"
	"github.com/bindl-dev/bindl/download"
	"github.com/bindl-dev/bindl/internal"
	"github.com/bindl-dev/bindl/program"
	"github.com/rs/zerolog"
)

func symlink(binDir, progDir string, p *program.Lock) error {
	relProgDir := filepath.Join(progDir, p.Name)

	// Update program atime and mtime to prevent Makefile from rebuilding
	// ref: https://stackoverflow.com/a/35276091
	now := time.Now().Local()
	if err := os.Chtimes(filepath.Join(binDir, relProgDir), now, now); err != nil {
		return err
	}

	symlinkPath := filepath.Join(binDir, p.Name)
	internal.Log().Debug().
		Str("program", p.Name).
		Dict("symlink", zerolog.Dict().
			Str("ref", relProgDir).
			Str("target", symlinkPath)).
		Msg("symlink program")
	_ = os.Remove(symlinkPath)
	return os.Symlink(filepath.Join(progDir, p.Name), symlinkPath)
}

// Get implements ProgramCommandFunc, therefore needs to be concurrent-safe.
// Before downloading, Get attempts to:
// - validate the existing installation
// - if it failed, redo symlink, then validate again
// - if it still fails, then attempt to download
// This is useful when a project is working on branches with different versions of
// a given program, ensuring that we only download when absolutely necessary.
func Get(ctx context.Context, conf *config.Runtime, p *program.Lock) error {
	archiveName, err := p.ArchiveName(conf.OS, conf.Arch)
	if err != nil {
		return err
	}

	progDir := filepath.Join(conf.ProgDir, p.Checksums[archiveName].Binary+"-"+p.Name)
	if err := Verify(ctx, conf, p); err == nil {
		// Re-run symlink to renew atime and mtime, so that GNU Make will not rebuild in the future
		internal.Log().Debug().Str("program", p.Name).Msg("found valid existing, re-linking")
		return symlink(conf.BinDir, progDir, p)
	}

	internal.Log().Debug().Err(err).Msg("verification failed")

	// Looks like verify failed, let's assume that the right version exists,
	// but was symlinked to the wrong one, therefore fix symlink and re-verify
	if err := symlink(conf.BinDir, progDir, p); err != nil {
		internal.Log().Debug().Err(err).Msg("failed symlink, downloading program")
	} else {
		if err := Verify(ctx, conf, p); err == nil {
			internal.Log().Debug().Str("program", p.Name).Msg("re-linked to appropriate version")
			// No need to return symlink() here, because we just ran symlink()
			return nil
		}
	}

	internal.Log().Debug().Err(err).Msg("verification failed after fixing symlink, redownloading")

	a, err := p.DownloadArchive(ctx, &download.HTTP{UseCache: conf.UseCache}, conf.OS, conf.Arch)
	if err != nil {
		return err
	}
	internal.Log().Debug().Str("program", p.Name).Msg("extracting archive")
	bin, err := a.Extract(p.Name)
	if err != nil {
		return err
	}
	internal.Log().Debug().Str("program", p.Name).Msg("found binary")

	fullProgDir := filepath.Join(conf.BinDir, progDir)
	if err = os.MkdirAll(fullProgDir, 0755); err != nil {
		return err
	}
	binPath := filepath.Join(fullProgDir, p.Name)
	err = os.WriteFile(binPath, bin, 0755)
	if err != nil {
		return err
	}
	internal.Log().Debug().Str("output", binPath).Str("program", p.Name).Msg("downloaded")

	return symlink(conf.BinDir, progDir, p)
}
