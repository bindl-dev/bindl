package program

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"path/filepath"

	"go.xargs.dev/bindl/internal"
)

func unzip(w io.Writer, rawZip io.ReaderAt, size int64, filename string) error {
	internal.Log().Debug().Int64("bytes", size).Str("file", filename).Msg("unzipping")
	r, err := zip.NewReader(rawZip, size)
	if err != nil {
		return fmt.Errorf("initializing zip reader: %w", err)
	}
	for _, f := range r.File {
		if f.FileInfo().IsDir() || filepath.Base(f.Name) != filename {
			continue
		}
		fd, err := f.Open()
		if err != nil {
			return fmt.Errorf("opening '%s' from zip archive: %w", f.Name, err)
		}
		_, err = io.Copy(w, fd)
		fd.Close()
		if err != nil {
			return fmt.Errorf("copying '%s' from zip archive: %w", f.Name, err)
		}
		return nil
	}
	return fmt.Errorf("unable to find '%s' from zip archive", filename)
}

func untar(w io.Writer, rawTar io.Reader, filename string) error {
	tarReader := tar.NewReader(rawTar)

	header, err := tarReader.Next()
	for {
		if err != nil {
			break
		}
		if header.Typeflag != tar.TypeReg || filepath.Base(header.Name) != filename {
			header, err = tarReader.Next()
			continue
		}

		_, err = io.Copy(w, tarReader)
		break
	}

	if errors.Is(err, io.EOF) {
		err = fmt.Errorf("unable to find '%s' in archive: %w", filename, err)
	}
	return err
}

func untargz(w io.Writer, rawTarGz io.Reader, filename string) error {
	gzReader, err := gzip.NewReader(rawTarGz)
	if err != nil {
		return err
	}
	defer gzReader.Close()

	return untar(w, gzReader, filename)
}
