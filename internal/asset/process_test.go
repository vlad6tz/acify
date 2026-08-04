package asset

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

func testImage(t *testing.T, w, h int) image.Image {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.Draw(img, img.Bounds(), image.NewUniform(color.RGBA{90, 120, 200, 255}), image.Point{}, draw.Src)
	for y := 0; y < h; y += 2 {
		for x := 0; x < w; x += 2 {
			img.Set(x, y, color.RGBA{255, 60, 40, 255})
		}
	}
	return img
}

func writeTestImage(t *testing.T, dir, name string, img image.Image) string {
	t.Helper()
	p := filepath.Join(dir, name)
	f, err := os.Create(p)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestContainFit(t *testing.T) {
	cases := []struct {
		w, h, maxW, maxH, wantW, wantH int
	}{
		{200, 100, 100, 100, 100, 50},   // scale down to fit box
		{50, 30, 200, 200, 50, 30},      // never upscale
		{800, 600, 400, 0, 400, 300},    // maxW only
		{800, 600, 0, 300, 400, 300},    // maxH only
		{100, 100, 0, 0, 100, 100},      // no limits -> unchanged
		{10, 10, 3, 2, 2, 2},            // very small target, floor at 1
		{1200, 300, 800, 800, 800, 200}, // wide image in square box
	}
	for _, c := range cases {
		gotW, gotH := containFit(c.w, c.h, c.maxW, c.maxH)
		if gotW != c.wantW || gotH != c.wantH {
			t.Errorf("containFit(%d,%d,%d,%d) = %dx%d, want %dx%d",
				c.w, c.h, c.maxW, c.maxH, gotW, gotH, c.wantW, c.wantH)
		}
	}
}

func TestProcessResizePng(t *testing.T) {
	dir := t.TempDir()
	p := writeTestImage(t, dir, "big.png", testImage(t, 200, 100))

	data, dims, _, err := Process(p, ProcessOptions{MaxW: 100, MaxH: 100})
	if err != nil {
		t.Fatal(err)
	}
	if dims.X != 100 || dims.Y != 50 {
		t.Fatalf("resized dims = %dx%d, want 100x50", dims.X, dims.Y)
	}
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("result is not a valid png: %v", err)
	}
	if img.Bounds().Dx() != 100 || img.Bounds().Dy() != 50 {
		t.Fatalf("decoded png dims = %dx%d, want 100x50", img.Bounds().Dx(), img.Bounds().Dy())
	}
}

func TestProcessNoUpscale(t *testing.T) {
	dir := t.TempDir()
	p := writeTestImage(t, dir, "small.png", testImage(t, 50, 30))

	_, dims, _, err := Process(p, ProcessOptions{MaxW: 400, MaxH: 400})
	if err != nil {
		t.Fatal(err)
	}
	if dims.X != 50 || dims.Y != 30 {
		t.Fatalf("upscaled dims = %dx%d, want original 50x30", dims.X, dims.Y)
	}
}

func TestProcessJpeg(t *testing.T) {
	dir := t.TempDir()
	p := writeTestImage(t, dir, "img.png", testImage(t, 64, 64))

	data, dims, _, err := Process(p, ProcessOptions{Format: "jpeg", Quality: 80})
	if err != nil {
		t.Fatal(err)
	}
	if dims.X != 64 || dims.Y != 64 {
		t.Fatalf("dims = %dx%d, want 64x64", dims.X, dims.Y)
	}
	if _, err := jpeg.Decode(bytes.NewReader(data)); err != nil {
		t.Fatalf("result is not a valid jpeg: %v", err)
	}
}

func TestProcessQualityBounds(t *testing.T) {
	dir := t.TempDir()
	p := writeTestImage(t, dir, "img.png", testImage(t, 32, 32))

	for _, q := range []int{0, 150, -5} {
		if _, _, _, err := Process(p, ProcessOptions{Format: "jpeg", Quality: q}); err != nil {
			t.Fatalf("quality %d: %v", q, err)
		}
	}
}

func TestProcessAnimatedGifRejected(t *testing.T) {
	dir := t.TempDir()
	pal := color.Palette{color.Black, color.White}
	imgs := []*image.Paletted{
		image.NewPaletted(image.Rect(0, 0, 4, 4), pal),
		image.NewPaletted(image.Rect(0, 0, 4, 4), pal),
	}
	g := &gif.GIF{Image: imgs, Delay: []int{10, 10}}
	p := filepath.Join(dir, "anim.gif")
	f, err := os.Create(p)
	if err != nil {
		t.Fatal(err)
	}
	if err := gif.EncodeAll(f, g); err != nil {
		t.Fatal(err)
	}
	f.Close()

	if _, _, _, err := Process(p, ProcessOptions{Format: "jpeg"}); err == nil {
		t.Fatal("expected animated gif to be rejected, got nil error")
	}
	if _, err := ImageDataURI(p); err != nil {
		t.Fatalf("pass-through inlining of an animated gif should still work: %v", err)
	}
}

func TestProcessSvgRejected(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "img.svg")
	if err := os.WriteFile(p, []byte("<svg xmlns='http://www.w3.org/2000/svg' width='4' height='4'/>"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, _, _, err := Process(p, ProcessOptions{Format: "jpeg"}); err == nil {
		t.Fatal("expected svg to be rejected, got nil error")
	}
}
