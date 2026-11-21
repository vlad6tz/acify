// Package config defines the detect-theme export structures.
package config

// Export is the top-level detect-theme JSON document.
type Export struct {
	Format        string       `json:"format"`
	FormatVersion int          `json:"format_version"`
	ExportedAt    string       `json:"exported_at"`
	Source        Source       `json:"source"`
	Theme         ThemePayload `json:"theme"`
}

// Source describes where the theme came from.
type Source struct {
	ThemeID       int    `json:"theme_id"`
	Name          string `json:"name"`
	Author        string `json:"author"`
	VersionNumber int    `json:"version_number"`
	IncludesDraft bool   `json:"includes_draft"`
	URL           string `json:"url"`
}

// Tokens holds theme token metadata.
type Tokens struct {
	Accent       string `json:"accent"`
	Bg           string `json:"bg"`
	BgSoft       string `json:"bg_soft"`
	BgCard       string `json:"bg_card"`
	Line         string `json:"line"`
	LineStrong   string `json:"line_strong"`
	Fg           string `json:"fg"`
	FgDim        string `json:"fg_dim"`
	FgMute       string `json:"fg_mute"`
	Danger       string `json:"danger"`
	Title        string `json:"title"`
	CornerRadius int    `json:"corner_radius"`
}

// ThemePayload holds the theme content.
type ThemePayload struct {
	Name      string            `json:"name"`
	ThemeType string            `json:"theme_type"`
	Tokens    Tokens            `json:"tokens"`
	Texts     map[string]string `json:"texts"`
	HTML      string            `json:"html"`
	CSS       string            `json:"css"`
	JS        string            `json:"js"`
}

// DefaultTokens returns the default token set matching the acify template.
func DefaultTokens(name string) Tokens {
	return Tokens{
		Accent:       "#23EE75",
		Bg:           "#0A0A0A",
		BgSoft:       "#111111",
		BgCard:       "#131313",
		Line:         "#1F1F1F",
		LineStrong:   "#2A2A2A",
		Fg:           "#F0F0F0",
		FgDim:        "#9A9A9A",
		FgMute:       "#5E5E5E",
		Danger:       "#FF5577",
		Title:        name,
		CornerRadius: 8,
	}
}

// NewExport builds a detect-theme export with sensible defaults.
func NewExport(name string, versionNumber int, author string, html, css, js string) Export {
	if name == "" {
		name = "acify"
	}
	if versionNumber <= 0 {
		versionNumber = 1
	}
	if author == "" {
		author = "acify"
	}
	return Export{
		Format:        "detect-theme",
		FormatVersion: 1,
		Source: Source{
			ThemeID:       0,
			Name:          name,
			Author:        author,
			VersionNumber: versionNumber,
			IncludesDraft: false,
			URL:           "",
		},
		Theme: ThemePayload{
			Name:      name,
			ThemeType: "advanced",
			Tokens:    DefaultTokens(name),
			Texts:     map[string]string{},
			HTML:      html,
			CSS:       css,
			JS:        js,
		},
	}
}
