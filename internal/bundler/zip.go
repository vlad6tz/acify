// Package bundler packs the HTML, CSS, JS, and asset files into an
// exportable zip archive.
package bundler

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"path/filepath"
)

// File is a single entry to write into the bundle.
type File struct {
	Name string
	Body []byte
}

// Bundle packs files into an in-memory zip archive and returns its bytes.
func Bundle(files []File) ([]byte, error) {
	buf := &bytes.Buffer{}
	zw := zip.NewWriter(buf)
	for _, f := range files {
		if f.Name == "" {
			continue
		}
		w, err := zw.Create(filepath.ToSlash(f.Name))
		if err != nil {
			return nil, fmt.Errorf("create zip entry %q: %w", f.Name, err)
		}
		if _, err := w.Write(f.Body); err != nil {
			return nil, fmt.Errorf("write zip entry %q: %w", f.Name, err)
		}
	}
	if err := zw.Close(); err != nil {
		return nil, fmt.Errorf("finalize zip: %w", err)
	}
	return buf.Bytes(), nil
}

// WriteBundle packs files into a zip and writes it to w.
func WriteBundle(w io.Writer, files []File) error {
	data, err := Bundle(files)
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}
