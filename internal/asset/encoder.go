// Package asset reads image files and generates base64 data URIs for inlining.
package asset

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var extMime = map[string]string{
	".png":  "image/png",
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".gif":  "image/gif",
	".webp": "image/webp",
	".bmp":  "image/bmp",
	".svg":  "image/svg+xml",
}

// EncodeBase64 reads all of r and returns its standard base64 encoding.
func EncodeBase64(r io.Reader) (string, error) {
	buf := &bytes.Buffer{}
	if _, err := io.Copy(buf, r); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

// DataURI reads all of r and returns it as a base64 data URI. When mime is
// empty, the MIME type is sniffed from the content.
func DataURI(r io.Reader, mime string) (string, error) {
	buf := &bytes.Buffer{}
	if _, err := io.Copy(buf, r); err != nil {
		return "", err
	}
	if mime == "" {
		mime = http.DetectContentType(buf.Bytes())
	}
	return fmt.Sprintf("data:%s;base64,%s", mime, base64.StdEncoding.EncodeToString(buf.Bytes())), nil
}

// Image validates path and returns its data URI, picking the MIME from the
// file extension when known (more reliable than sniffing for webp/svg).
func Image(path string) (string, error) {
	if err := ValidateImage(path); err != nil {
		return "", err
	}
	return ImageDataURI(path)
}

// ImageDataURI reads the image file at path and returns it as a data URI.
func ImageDataURI(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	return DataURI(f, mimeFor(path))
}

// ValidateImage checks that path exists, is a file, and (for raster formats)
// decodes as an image. SVG is accepted without decoding.
func ValidateImage(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return fmt.Errorf("%s: is a directory, expected an image", path)
	}
	if strings.EqualFold(filepath.Ext(path), ".svg") {
		return nil
	}
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, _, err := image.DecodeConfig(f); err != nil {
		return fmt.Errorf("%s: not a valid image: %w", path, err)
	}
	return nil
}

func mimeFor(path string) string {
	if m, ok := extMime[strings.ToLower(filepath.Ext(path))]; ok {
		return m
	}
	return ""
}
