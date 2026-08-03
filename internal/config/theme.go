// Package config defines the theme manifest structures.
package config

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

// Default returns a Theme with the standard defaults, suitable as a starting
// manifest for a new style.
func Default() Theme {
	return Theme{
		Name:    "acify",
		Version: "1.0.0",
	}
}
