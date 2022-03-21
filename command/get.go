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
	"path"

	"go.xargs.dev/bindl/config"
	"go.xargs.dev/bindl/download"
	"go.xargs.dev/bindl/internal"
	"go.xargs.dev/bindl/program"
)

// Get implements lockfileProgramCommandFunc, therefore needs to be concurrent-safe
func Get(ctx context.Context, conf *config.Runtime, p *program.URLProgram) error {
	if err := Verify(ctx, conf, p); err == nil {
		internal.Log().Debug().Str("program", p.Name()).Msg("found existing, skipping")
		return nil
	}
	a, err := p.DownloadArchive(ctx, &download.HTTP{}, conf.OS, conf.Arch)
	if err != nil {
		return err
	}
	internal.Log().Debug().Str("program", p.Name()).Msg("extracting archive")
	bin, err := a.Extract(p.Name())
	if err != nil {
		return err
	}
	internal.Log().Debug().Str("program", p.Name()).Msg("found binary")

	if err = os.MkdirAll(conf.OutputDir, 0755); err != nil {
		return err
	}
	loc := path.Join(conf.OutputDir, p.Name())
	err = os.WriteFile(loc, bin, 0755)
	if err != nil {
		return err
	}
	internal.Log().Info().Str("output", loc).Str("program", p.Name()).Msg("downloaded")
	return nil
}
