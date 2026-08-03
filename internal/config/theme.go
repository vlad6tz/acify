// Package config defines the theme manifest structures.
package config

import "os"

// Theme is the manifest for a custom style bundle.
type Theme struct {
	Name    string
	Version string
	Assets  []Asset
	Options []Option
}

// Asset is an image file inlined into the style as a data URI.
type Asset struct {
	Name string
	Path string
}

// Option is a user-configurable theme option.
type Option struct {
	ID      string
	Label   string
	Type    string
	Default string
	Choices []string
}

// Default returns a Theme populated from env-provided values, suitable as a
// starting manifest for a new style.
func Default() Theme {
	return Theme{
		Name:    valueOr("ACIFY_NAME", "detect"),
		Version: valueOr("ACIFY_VERSION", "1.0.0"),
	}
}

func valueOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
