package bundler

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	markerBackground = "/* acify:background */"
	markerSpinner    = "/* acify:spinner */"
)

// InjectCSS replaces the background and spinner markers in css with data-URI declarations.
// Pass empty strings to leave the template default in place,
// overlay controls the dark tint over a background image (0.0-1.0),
// accent overrides the --accent CSS variable when non-empty.
func InjectCSS(css, background, spinner, accent string, overlay float64, rotate bool) string {
	if background != "" {
		rule := fmt.Sprintf("background-image: linear-gradient(rgba(8,8,10,%.2f), rgba(8,8,10,%.2f)), url(\"%s\");",
			clamp(overlay), clamp(overlay), background)
		css = replaceOrAppend(css, markerBackground, rule, "body { "+rule+" }")
	} else {
		css = dropMarker(css, markerBackground)
	}

	switch {
	case spinner != "" && rotate:
		decl := "background: url(\"" + spinner + "\") center/contain no-repeat; border: none; filter: none;"
		css = replaceOrAppend(css, markerSpinner, decl, ".spinner { "+decl+" }")
	case spinner != "" && !rotate:
		decl := "background: url(\"" + spinner + "\") center/contain no-repeat; border: none; filter: none; animation: none;"
		css = replaceOrAppend(css, markerSpinner, decl, ".spinner { "+decl+" }")
	case spinner == "" && !rotate:
		css = replaceOrAppend(css, markerSpinner, "animation: none;", ".spinner { animation: none; }")
	default: // no spinner, rotating -> leave the template wireframe as-is
		css = dropMarker(css, markerSpinner)
	}

	if accent != "" {
		css += accentOverride(accent)
	}
	return css
}

// replaceOrAppend puts replacement where marker sits, or appends fallback rule
// when the marker is absent (custom template).
func replaceOrAppend(css, marker, replacement, fallback string) string {
	if strings.Contains(css, marker) {
		return strings.ReplaceAll(css, marker, replacement)
	}
	return css + "\n" + fallback + "\n"
}

func dropMarker(css, marker string) string {
	if strings.Contains(css, marker) {
		return strings.ReplaceAll(css, marker, "")
	}
	return css
}

func accentOverride(accent string) string {
	accent = strings.TrimSpace(accent)
	if r, g, b, ok := hexToRGB(accent); ok {
		return fmt.Sprintf("\n:root { --accent: %s; --accent-soft: rgba(%d,%d,%d,0.10); --accent-glow: rgba(%d,%d,%d,0.35); }\n",
			accent, r, g, b, r, g, b)
	}
	return fmt.Sprintf("\n:root { --accent: %s; }\n", accent)
}

func hexToRGB(s string) (uint32, uint32, uint32, bool) {
	h := strings.TrimPrefix(s, "#")
	if len(h) != 6 {
		return 0, 0, 0, false
	}
	v, err := strconv.ParseUint(h, 16, 24)
	if err != nil {
		return 0, 0, 0, false
	}
	return uint32(v >> 16 & 0xff), uint32(v >> 8 & 0xff), uint32(v & 0xff), true
}

func clamp(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}
