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

func TestMaxFileBytes(t *testing.T) {
	if MaxFileBytes != 256*1024 {
		t.Fatalf("MaxFileBytes = %d, want %d", MaxFileBytes, 256*1024)
	}
}
