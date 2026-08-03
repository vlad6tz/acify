package validator

import "testing"

// base64 of a small PNG chosen to contain "//" so that naive scanners would
// mistake it for a comment.
const dataURIWithSlash = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJ//AAAADUlEQVR42mP8z8BQDwAEhQGAhKmMIQAAAABJRU5ErkJggg=="

func TestCSSQuotedDataURI(t *testing.T) {
	css := "body { background-image: url(\"" + dataURIWithSlash + "\"); }"
	if err := CSS(css); err != nil {
		t.Fatalf("CSS() = %v, want nil", err)
	}
}

func TestCSSUnquotedDataURI(t *testing.T) {
	css := "body { background-image: url(" + dataURIWithSlash + "); }"
	if err := CSS(css); err != nil {
		t.Fatalf("CSS() = %v, want nil", err)
	}
}

func TestCSSBlockAfterDataURI(t *testing.T) {
	css := "body { background-image: url(" + dataURIWithSlash + "); }\n.spinner { border: none; }"
	if err := CSS(css); err != nil {
		t.Fatalf("CSS() = %v, want nil", err)
	}
}

func TestCSSMismatchStillDetected(t *testing.T) {
	css := "body { background-image: url(no); } }"
	if err := CSS(css); err == nil {
		t.Fatal("CSS() = nil, want mismatch error")
	}
}

func TestSize(t *testing.T) {
	if err := Size("theme.css", 255*1024); err != nil {
		t.Fatalf("Size(255 KiB) = %v, want nil", err)
	}
	if err := Size("theme.css", 256*1024); err == nil {
		t.Fatal("Size(256 KiB) = nil, want error")
	}
}
