// Package asset reads image files and generates base64 data URIs for inlining.
package asset

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
)

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

// ImageDataURI reads the image file at path and returns it as a data URI.
func ImageDataURI(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	return DataURI(f, "")
}
