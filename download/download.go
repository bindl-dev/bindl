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
