# AGENTS.md

Use this file to turn a user's theme preferences into the correct `acify` CLI
command, with no assistance from the user. If a preference is not stated, keep
the flag unset and rely on the default.

`acify` is a Go (1.25) / Cobra CLI. It inlines a required background image
(`-b`) and an optional spinner image (`-s`) into `templates/default.{html,css,js}`
as base64 data URIs, then writes a single detect-theme JSON export
(`<name>-v<version-number>.detect-theme.json`) compatible with detect.ac theme
import.

## Flag reference

| Flag | Shorthand | Default | Meaning |
| --- | --- | --- | --- |
| `--background` | `-b` | *(required)* | Path to the background image (`.png` `.jpg` `.jpeg` `.gif` `.webp` `.bmp` `.svg`) |
| `--spinner` | `-s` | *(none)* | Optional spinner image |
| `--no-spin` | | `false` | Disable spinner rotation (spins by default) |
| `--no-sign` | | `false` | Drop the "built with acify" signature comment |
| `--out` | `-o` | `.` | Directory for the detect-theme JSON file |
| `--templates` | | `templates` | Dir containing `default.{html,css,js}` |
| `--name` | | `acify` | Theme name |
| `--version-number` | | `1` | Source `version_number` in the export |
| `--author` | | `acify` | Source `author` in the export |
| `--theme-id` | | `0` | Source `theme_id` in the export |
| `--url` | | *(none)* | Source `url` in the export |
| `--theme-type` | | `advanced` | Theme `theme_type` in the export |
| `--overlay` | | `0.4` | Dark tint over the background (0.0-1.0) |
| `--accent` | | *(none)* | Override the `--accent` color, e.g. `#ff5c72` |
| `--css` | | *(none)* | Append a custom CSS file to the bundle |
| `--resize` | | *(none)* | Downscale the background to fit within `WxH`, e.g. `1280x720` (preserves aspect, never upscales) |
| `--to` | | *(none)* | Re-encode the background (`png` or `jpeg`) |
| `--quality` | | `85` | JPEG encoder quality (0-100), ignored for PNG |
| `--fit` | | `false` | Auto-shrink the background to fit a 256 KiB soft budget |
| `--validate-only` | | `false` | Validate templates/images/output without writing a file |
| `--inline` | `-i` | `false` | Print generated CSS instead of writing a file |
| `--verbose` | `-v` | `false` | Verbose output |
| `--help` | `-h` | *(auto)* | Show usage and all flags (added by Cobra) |

## Command-building rules

- `-b` is always required; derive it from the user's chosen image and verify the
  format is in the supported list above.
- Prefer short flags (`-b`, `-s`). Leave anything the user did not mention unset
  - the defaults are safe.
- Validate inputs before emitting: `--overlay` must be between 0.0 and 1.0;
  `--quality` between 0 and 100; `--to` is `png` or `jpeg`; `--resize` is
  `WxH`; `--version-number` and `--theme-id` are integers.
- Suggest `--resize`/`--to`/`--quality` when a large image may bloat the export,
  or `--fit` to let the tool auto-shrink until it fits a 256 KiB soft budget.
- Spinner rotates by default. Only add `--no-spin` when the user wants a
  non-rotating spinner.
- Add `--no-sign` when the user wants clean files without the "built with
  acify" signature comment.
- Contract: output is a single detect-theme JSON file with
  `format="detect-theme"`, `format_version=1`, `source`, and
  `theme{name,theme_type,tokens,texts,html,css,js}`.

## Build & verify

```sh
go build -o acify.exe ./cmd/helper
.\acify.exe -b bg.png -o dist        # PowerShell: note the .\
gofmt -l .
go vet ./...
go build ./...
go test ./...
```

Self-check: after building, run
`.\acify.exe -b <image> --validate-only -o dist` to prove the generated command
is valid before telling the user it works.

## Editing invariants

- `internal/bundler/inject.go` - `InjectCSS(...)` replaces the markers
  `/* acify:background */` and `/* acify:spinner */` in `templates/default.css`.
- Injected data URIs must stay double-quoted (`url("data:...")`) - a bare `//`
  in unquoted base64 breaks brace-balancing in `internal/validator`.
- `cmd/helper/root.go` holds the `options` struct and all flag definitions.
- `internal/config/theme.go` defines the detect-theme export schema; keep it
  aligned with the target `*.detect-theme.json` format.