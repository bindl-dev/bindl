package download

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

type HTTP struct {
	response *http.Response
}

func (d *HTTP) Get(ctx context.Context, url string) (io.Reader, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("composing request: %w", err)
	}

	c := http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("downloading program: %w", err)
	}
	d.response = resp

	return d.response.Body, nil
}

func (d *HTTP) Close() {
	if d.response == nil {
		return
	}
	if d.response.Body == nil {
		return
	}
	d.response.Body.Close()
}
