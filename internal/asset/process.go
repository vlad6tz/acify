package asset

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	stddraw "image/draw"
	"image/jpeg"
	"image/png"
	"os"

	xdraw "golang.org/x/image/draw" // for CatmullRom scaling
)

// ProcessOptions controls how an image is re-encoded.
type ProcessOptions struct {
	// MaxW and MaxH bound the image within a box, preserving aspect ratio.
	// Zero values are ignored.
	MaxW, MaxH int
	// Format re-encodes to "png" or "jpeg" when non-empty. Empty keeps the
	// original format when it can be encoded, otherwise falls back to png.
	Format string
	// Quality is the JPEG encoder quality (0-100). Ignored for PNG.
	Quality int
}

// Optimize decodes the image at path, applies the processing options, and
// returns the resulting data URI plus a human-readable summary of what was
// done (e.g. "jpeg 640x360 q80"). When opts is empty the file is inlined
// unchanged.
func Optimize(path string, opts ProcessOptions) (uri string, report string, err error) {
	if opts == (ProcessOptions{}) {
		uri, err := ImageDataURI(path)
		return uri, "", err
	}

	data, dims, format, err := Process(path, opts)
	if err != nil {
		return "", "", err
	}
	mime := "image/" + format
	uri = fmt.Sprintf("data:%s;base64,%s", mime, base64.StdEncoding.EncodeToString(data))

	if format == "jpeg" {
		report = fmt.Sprintf("jpeg %dx%d q%d", dims.X, dims.Y, opts.Quality)
	} else {
		report = fmt.Sprintf("png %dx%d", dims.X, dims.Y)
	}
	return uri, report, nil
}

// Process reads the image at path, optionally resizes and re-encodes it, and
// returns the encoded bytes, the resulting dimensions, and the format the
// bytes were actually encoded as ("png" or "jpeg").
func Process(path string, opts ProcessOptions) ([]byte, image.Point, string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, image.Point{}, "", err
	}
	defer f.Close()

	src, format, err := image.Decode(f)
	if err != nil {
		return nil, image.Point{}, "", err
	}
	bounds := src.Bounds()
	w, h := bounds.Dx(), bounds.Dy()

	if opts.MaxW > 0 || opts.MaxH > 0 {
		w, h = containFit(w, h, opts.MaxW, opts.MaxH)
	}

	out := src
	if w != bounds.Dx() || h != bounds.Dy() {
		dst := image.NewRGBA(image.Rect(0, 0, w, h))
		xdraw.CatmullRom.Scale(dst, dst.Bounds(), src, bounds, xdraw.Over, nil)
		out = dst
	}

	if opts.Format != "" && opts.Format != format {
		format = opts.Format
	}
	// Go can only encode png/jpeg; anything else decodes to png.
	if format != "png" && format != "jpeg" {
		format = "png"
	}

	var buf bytes.Buffer
	switch format {
	case "jpeg":
		if opts.Quality <= 0 {
			opts.Quality = 85
		}
		if opts.Quality > 100 {
			opts.Quality = 100
		}
		// JPEG has no alpha; flatten onto black before encoding.
		opaque := flatten(out)
		err = jpeg.Encode(&buf, opaque, &jpeg.Options{Quality: opts.Quality})
	case "png":
		err = png.Encode(&buf, out)
	default:
		return nil, image.Point{}, "", fmt.Errorf("unsupported format %q", format)
	}
	if err != nil {
		return nil, image.Point{}, "", err
	}
	return buf.Bytes(), image.Pt(w, h), format, nil
}

// containFit scales a w x h box to fit within maxW x maxH, preserving aspect
// ratio and never upscaling.
func containFit(w, h, maxW, maxH int) (int, int) {
	if maxW <= 0 {
		maxW = w
	}
	if maxH <= 0 {
		maxH = h
	}
	sx, sy := 1.0, 1.0
	if w > maxW {
		sx = float64(maxW) / float64(w)
	}
	if h > maxH {
		sy = float64(maxH) / float64(h)
	}
	scale := sx
	if sy < sx {
		scale = sy
	}
	if scale >= 1.0 {
		return w, h
	}
	return max(1, int(float64(w)*scale)), max(1, int(float64(h)*scale))
}

// flatten draws img onto an opaque black background, returning an image that
// has no alpha channel (needed before JPEG encoding).
func flatten(img image.Image) image.Image {
	b := img.Bounds()
	opaque := image.NewRGBA(b)
	stddraw.Draw(opaque, b, image.NewUniform(image.Black), image.Point{}, stddraw.Src)
	stddraw.Draw(opaque, b, img, b.Min, stddraw.Over)
	return opaque
}
