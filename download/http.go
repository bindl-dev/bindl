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
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/bindl-dev/bindl/internal"
	"github.com/bindl-dev/httpcache"
	"github.com/bindl-dev/httpcache/diskcache"
)

// HTTP implements Downloader which downloads programs through net/http
//nolint:govet  // bytes saved isn't worth the reduced visibility
type HTTP struct {
	UseCache bool

	closeBodyOnce sync.Once
	body          io.ReadCloser
}

var client *http.Client

func init() {
	var cacheDir = filepath.Join(".bindlcache", "http")
	userCacheDir, err := os.UserCacheDir()
	if err != nil {
		internal.Log().Debug().Err(err).Msg("finding user cache directory")
	} else {
		cacheDir = filepath.Join(userCacheDir, "bindl", "http")
	}
	client = httpcache.NewTransport(diskcache.New(cacheDir)).Client()
}

func (d *HTTP) Get(ctx context.Context, url string) (io.Reader, error) {
	var c *http.Client
	if d.UseCache {
		c = client
	} else {
		c = http.DefaultClient
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("composing request: %w", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("downloading program: %w", err)
	}
	d.body = resp.Body

	if resp.Header.Get(httpcache.XFromCache) != "" {
		internal.Log().Debug().Str("url", url).Msg("cached response")
	}

	return d.body, nil
}

func (d *HTTP) Close() {
	d.closeBodyOnce.Do(func() {
		if d.body == nil {
			return
		}
		d.body.Close()
	})
}
