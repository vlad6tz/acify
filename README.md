# acify

Build custom style bundles for [detect.ac](https://detect.ac) from a background image.

`acify` takes your background image (and optional spinner image), inlines it into the
default HTML/CSS/JS templates as base64 data URIs, and packs the three files into a zip.
Each file is signed with the repo URL and stays strictly under the 256 KiB per-file
limit required by detect.ac.

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
acify -b background.png                      # required - build from a background image
acify -b bg.png -s spinner.gif -o dist       # optional spinner + output directory
acify -b bg.png -s spinner.gif --no-spin     # spinner image without rotation
acify -b bg.png --name mine --version 1.2.0  # theme name / version in the zip filename
```

The generated zip contains `theme.html`, `theme.css`, and `theme.js`. It is not
uploaded anywhere - extract it and paste the three files into the detect.ac theme
editor manually.

## Flags

| Flag | Shorthand | Default | Description |
| --- | --- | --- | --- |
| `--background` | `-b` | *(required)* | Path to the background image |
| `--spinner` | `-s` | | Optional spinner image |
| `--no-spin` | | | Disable spinner rotation (spins by default) |
| `--out` | `-o` | `.` | Directory for the bundle zip |
| `--templates` | | `templates` | Dir containing `default.{html,css,js}` |
| `--name` | | `acify` | Theme name (used in the zip filename) |
| `--version` | | `1.0.0` | Theme version |
| `--overlay` | | `0.4` | Dark tint over the background (0.0-1.0) |
| `--accent` | | | Override the `--accent` color, e.g. `#ff5c72` |
| `--css` | | | Append a custom CSS file to the bundle |
| `--validate-only` | | | Check templates/images/sizes without writing a zip |
| `--inline` | `-i` | | Print data URIs + CSS instead of building |
| `--verbose` | `-v` | | Verbose output |

## Project layout

```
cmd/helper/            CLI entry point (Cobra)
internal/asset/        image validation & base64 data URIs
internal/bundler/      marker injection + zip packing
internal/config/       theme manifest structs
internal/validator/    HTML/CSS/JS checks + 256 KiB limit
templates/             default.{html,css,js} boilerplate
AGENTS.md              guidance for AI agents building acify commands
```

## Agent guidance

`AGENTS.md` describes the CLI's flags, defaults, and invariants so an AI agent
can build the correct `acify` command from a user's theme preferences without
extra assistance.

## License

MIT
