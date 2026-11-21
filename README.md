# acify

Build [detect.ac](https://detect.ac) detect-theme JSON exports from a background image.

`acify` takes your background image (and optional spinner image), inlines it into the
default HTML/CSS/JS templates as base64 data URIs, and writes a single
`*.detect-theme.json` file following the detect-theme export format. HTML/CSS/JS
are signed with the repo URL by default; pass `--no-sign` to drop the signature
comment.

## Requirements

- **Go 1.25+** to build the CLI.
- A **background image** (`.png`, `.jpg/.jpeg`, `.gif`, `.webp`, `.bmp`, or `.svg`).
  The background is mandatory.

## Install

Requires [Go 1.25+](https://go.dev/dl).

```sh
# clone the repo
git clone https://github.com/vlad6tz/acify.git
cd acify

# build the binary
go build -o acify.exe ./cmd/helper
```

## Usage

```sh
acify -b background.png                         # build a detect-theme JSON export
acify -b bg.png -s spinner.gif -o dist          # optional spinner + output directory
acify -b bg.png -s spinner.gif --no-spin        # spinner image without rotation
acify -b big.png --resize 1280x720 --to jpeg    # shrink/re-encode a large background
acify -b huge.gif --fit                         # auto-shrink to a 256 KiB soft budget
acify -b bg.png --name mine --version-number 2  # theme name / source version number
```

The generated file is a single detect-theme JSON export (for example
`acify-v1.detect-theme.json`). It is not uploaded anywhere - import it into
detect.ac manually.

## Flags

| Flag | Shorthand | Default | Description |
| --- | --- | --- | --- |
| `--background` | `-b` | *(required)* | Path to the background image |
| `--spinner` | `-s` | | Optional spinner image |
| `--no-spin` | | | Disable spinner rotation (spins by default) |
| `--no-sign` | | | Drop the "built with acify" signature comment |
| `--out` | `-o` | `.` | Directory for the detect-theme JSON file |
| `--templates` | | `templates` | Dir containing `default.{html,css,js}` |
| `--name` | | `acify` | Theme name |
| `--version-number` | | `1` | Source `version_number` in the export |
| `--author` | | `acify` | Source `author` in the export |
| `--theme-id` | | `0` | Source `theme_id` in the export |
| `--url` | | | Source `url` in the export |
| `--theme-type` | | `advanced` | Theme `theme_type` in the export |
| `--overlay` | | `0.4` | Dark tint over the background (0.0-1.0) |
| `--accent` | | | Override the `--accent` color, e.g. `#ff5c72` |
| `--css` | | | Append a custom CSS file to the bundle |
| `--resize` | | | Downscale the background to fit within `WxH`, preserving aspect ratio (e.g. `1280x720`) |
| `--to` | | | Re-encode the background image (`png` or `jpeg`) |
| `--quality` | | `85` | JPEG encoder quality (0-100), ignored for PNG |
| `--fit` | | | Auto-shrink the background to fit a 256 KiB soft budget |
| `--validate-only` | | | Validate templates/images/output without writing a file |
| `--inline` | `-i` | | Print generated CSS instead of building a file |
| `--verbose` | `-v` | | Verbose output |
| `--help` | `-h` | | Show usage and all flags (added by Cobra) |

## Output format

The generated JSON follows the detect-theme export pattern:

```json
{
  "format": "detect-theme",
  "format_version": 1,
  "exported_at": "2026-11-21T14:32:28Z",
  "source": { ... },
  "theme": {
    "name": "...",
    "theme_type": "advanced",
    "tokens": { ... },
    "texts": {},
    "html": "...",
    "css": "...",
    "js": "..."
  }
}
```

## Project layout

```
cmd/helper/            CLI entry point (Cobra)
internal/asset/        image validation & base64 data URIs
internal/bundler/      CSS marker injection + detect-theme JSON export
internal/config/       detect-theme export structs
internal/validator/    HTML/CSS/JS structural checks
templates/             default.{html,css,js} boilerplate
AGENTS.md              guidance for AI agents building acify commands
```

## Agent guidance

`AGENTS.md` describes the CLI's flags, defaults, and invariants so an AI agent
can build the correct `acify` command from a user's theme preferences without
extra assistance.

## Contributing

Contributions are welcome. To work on the tool locally:

```sh
# fork the repo on GitHub, then clone your fork
git clone https://github.com/<your-username>/acify.git
cd acify
git remote add upstream https://github.com/vlad6tz/acify.git
git checkout -b my-feature

# build
go build -o acify.exe ./cmd/helper

# verify before submitting
gofmt -l .
go vet ./...
go build ./...
go test ./...
```

Run the checks above and add tests for any new behavior. Push your branch to
your fork and open a pull request against `upstream`. When changing the CLI
surface, keep `README.md` and `AGENTS.md` flag tables in sync.

## License

MIT
