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
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/bindl-dev/bindl/config"
	"github.com/bindl-dev/bindl/internal"
	"github.com/bindl-dev/bindl/program"
)

// Verify implements ProgramCommandFunc, therefore needs to be concurrent-safe
// It verifies existing if the exiting program is consistent with what is declared
// by the lockfile.
func Verify(ctx context.Context, conf *config.Runtime, prog *program.Lock) error {
	binPath := filepath.Join(conf.BinDir, prog.Name)
	f, err := os.Open(binPath)
	if err != nil {
		return fmt.Errorf("opening '%v': %w", binPath, err)
	}
	defer f.Close()

	archiveName, err := prog.ArchiveName(conf.OS, conf.Arch)
	if err != nil {
		return fmt.Errorf("generating filename for '%v': %w", prog.Name, err)
	}
	expected := prog.Checksums[archiveName].Binaries[prog.Name]

	c := &program.ChecksumCalculator{}
	w := c.SHA256(io.Discard)
	if _, err = io.Copy(w, f); err != nil {
		return fmt.Errorf("reading checksum for '%v': %w", prog.Name, err)
	}

	if err := c.Error([]byte(expected)); err != nil {
		return err
	}

	internal.Log().Debug().Str("program", prog.Name).Msg("validated")

	return nil
}
