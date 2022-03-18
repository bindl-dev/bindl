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

package download

import (
	"context"
	"io"
)

type Downloader interface {
	Get(ctx context.Context, url string) (io.Reader, error)
	Close()
}

type LocalFile struct {
	r io.ReadCloser
}

func NewLocalFile(r io.ReadCloser) *LocalFile {
	return &LocalFile{r: r}
}

func (l *LocalFile) Get(ctx context.Context, url string) (io.Reader, error) {
	return l.r, nil
}

func (l *LocalFile) Close() {
	l.r.Close()
}
