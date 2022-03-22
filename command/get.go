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

	"go.xargs.dev/bindl/config"
	"go.xargs.dev/bindl/download"
	"go.xargs.dev/bindl/internal"
	"go.xargs.dev/bindl/program"
)

func symlink(progDir, symlinkPath string, p *program.URLProgram) error {
	internal.Log().Debug().Str("program", p.PName).Str("progDir", progDir).Msg("symlink program")
	_ = os.Remove(symlinkPath)
	return os.Symlink(filepath.Join(progDir, p.PName), symlinkPath)
}

// Get implements ProgramCommandFunc, therefore needs to be concurrent-safe
func Get(ctx context.Context, conf *config.Runtime, p *program.URLProgram) error {
	err := Verify(ctx, conf, p)
	if err == nil {
		internal.Log().Debug().Str("program", p.PName).Msg("found existing, skipping")
		return nil
	}

	internal.Log().Debug().Err(err).Msg("verification failed")

	// Looks like verify failed, let's assume that the right version exists,
	// but was symlinked to the wrong one, therefore fix symlink and re-verify
	archiveName, err := p.ArchiveName(conf.OS, conf.Arch)
	if err != nil {
		return err
	}
	progDir := filepath.Join(conf.ProgDir, p.Checksums[archiveName].Binaries[p.PName]+"-"+p.PName)
	fullProgDir := filepath.Join(conf.BinPathDir, progDir)
	symlinkPath := filepath.Join(conf.BinPathDir, p.PName)
	if err := symlink(progDir, symlinkPath, p); err != nil {
		internal.Log().Debug().Err(err).Msg("failed symlink, donwloading program")
	} else {
		if err = Verify(ctx, conf, p); err == nil {
			internal.Log().Debug().Str("program", p.PName).Msg("re-linked to appropriate version")
			return nil
		}
	}

	internal.Log().Debug().Err(err).Msg("verification failed after fixing symlink, redownloading")

	a, err := p.DownloadArchive(ctx, &download.HTTP{}, conf.OS, conf.Arch)
	if err != nil {
		return err
	}
	internal.Log().Debug().Str("program", p.PName).Msg("extracting archive")
	bin, err := a.Extract(p.PName)
	if err != nil {
		return err
	}
	internal.Log().Debug().Str("program", p.PName).Msg("found binary")

	if err = os.MkdirAll(fullProgDir, 0755); err != nil {
		return err
	}
	binPath := filepath.Join(fullProgDir, p.PName)
	err = os.WriteFile(binPath, bin, 0755)
	if err != nil {
		return err
	}
	internal.Log().Debug().Str("output", binPath).Str("program", p.PName).Msg("downloaded")

	return symlink(progDir, symlinkPath, p)
}
