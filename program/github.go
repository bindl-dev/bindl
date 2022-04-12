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

package program

import (
	"bytes"
	"text/template"

	"github.com/bindl-dev/bindl/internal"
)

var (
	tmplGitHubBase = template.Must(template.New("url").Parse("https://github.com/{{ .Base }}/releases/download/"))
)

func githubToURL(c *Config) error {
	var buf bytes.Buffer
	if err := tmplGitHubBase.Execute(&buf, c.Paths); err != nil {
		return err
	}
	c.Paths.Base = buf.String() + "v{{ .Version }}"
	internal.Log().Debug().Str("url", c.Paths.Base).Msg("github to url")
	return nil
}
