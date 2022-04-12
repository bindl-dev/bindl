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

package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/bindl-dev/bindl/command"
	"github.com/bindl-dev/bindl/config"
	"github.com/bindl-dev/bindl/internal"
	"github.com/bindl-dev/bindl/program/bootstrap"
)

var bindlGetBootstrap = false

func init() {
	BindlGet.Flags().BoolVar(&bindlGetBootstrap, "bootstrap", bindlGetBootstrap, "get bootstrapped program")
}

var BindlGet = &cobra.Command{
	Use:   "get [name, ...]",
	Short: "Get local copy of program",
	Long: `Get downloads the names program, which must already exist in bindl.yaml,
and ensures the program is ready to be used by setting executable flag. If no 
program name is specified through args, then all programs in lockfile will be selected.

While it is unlikely for end-user to need it, the flag --bootstrap is provided to download
internally trusted program. Bootstrap mode uses pre-defined values of program validations
at compile time. In bootstrap mode, program name must be specified in args.`,
	RunE: func(cmd *cobra.Command, names []string) error {
		if bindlGetBootstrap {
			return getBootstrap(cmd.Context(), names)
		}
		err := command.IterateLockfilePrograms(
			cmd.Context(),
			conf,
			names,
			command.Get)
		return err
	},
}

func getBootstrap(ctx context.Context, names []string) error {
	if len(names) < 1 {
		return fmt.Errorf("bootstrap mode requires program name to be specified")
	}

	for _, name := range names {
		manifest, err := bootstrap.Lock(name)
		if err != nil {
			return err
		}
		lock, err := config.ParseLockBytes(manifest)
		if err != nil {
			return err
		}
		if len(lock.Programs) == 0 {
			return fmt.Errorf("no programs were found in the bootstrap manifest for %v, please report this bug", name)
		}
		if len(lock.Programs) > 1 {
			return fmt.Errorf("multiple programs were found in the bootstrap manifest for %v, please report this bug", name)
		}
		if err := command.Get(ctx, conf, lock.Programs[0]); err != nil {
			return err
		}
		internal.Log().Info().Str("program", name).Msg("bootstrap successful")
	}
	return nil
}
