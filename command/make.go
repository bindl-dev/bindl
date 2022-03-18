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
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"go.xargs.dev/bindl/config"
)

var rawMakefileTmpl = "\n{{ .OutputDir }}/{{ .Name }}: {{ .OutputDir }}/bindl\n\t{{ .OutputDir }}/bindl get {{ .Name }}\n"

var makefileTmpl = template.Must(template.New("makefile").Parse(rawMakefileTmpl))

var makefileHeader = `# DO NOT MODIFY - THIS FILE WAS GENERATED BY BINDL
# ANY MODIFICATIONS WILL BE OVERWRITTEN

# If file content is outdated, run 'bindl make' to regenerate.
`

func GenerateMakefile(conf *config.Runtime, path string) error {
	l, err := config.ParseLock(conf.LockfilePath)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(path, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("opening file '%s': %w", path, err)
	}

	if _, err := f.Write([]byte(makefileHeader)); err != nil {
		return fmt.Errorf("writing makefile header: %w", err)
	}

	m := map[string]string{
		"OutputDir": filepath.Base(conf.OutputDir),
	}

	for _, prog := range l.Programs {
		m["Name"] = prog.PName
		if err := makefileTmpl.Execute(f, m); err != nil {
			return fmt.Errorf("writing template: %w", err)
		}
	}

	return nil
}
