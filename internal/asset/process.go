package asset

import (
	"bytes"
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
	// MaxW and MaxH bound the image to fit within a box, preserving aspect
	// ratio. Zero values are ignored.
	MaxW, MaxH int
	// Format is "png" or "jpeg". Empty keeps the original format.
	Format string
	// Quality is the JPEG encoder quality (0-100). Ignored for PNG.
	Quality int
}

// Process reads the image at path, optionally resizes and re-encodes it, and
// returns the encoded bytes plus the resulting dimensions.
func Process(path string, opts ProcessOptions) ([]byte, image.Point, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, image.Point{}, err
	}
	defer f.Close()

	src, format, err := image.Decode(f)
	if err != nil {
		return nil, image.Point{}, err
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
		return nil, image.Point{}, fmt.Errorf("unsupported format %q", format)
	}
	if err != nil {
		return nil, image.Point{}, err
	}
	return buf.Bytes(), image.Pt(w, h), nil
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
