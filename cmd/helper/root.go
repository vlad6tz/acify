package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/vlad6tz/acify.git/internal/asset"
	"github.com/vlad6tz/acify.git/internal/bundler"
	"github.com/vlad6tz/acify.git/internal/config"
	"github.com/vlad6tz/acify.git/internal/validator"
)

type options struct {
	background    string
	spinner       string
	noSpin        bool
	noSign        bool
	out           string
	templates     string
	name          string
	versionNumber int
	author        string
	themeID       int
	url           string
	themeType     string
	overlay       float64
	accent        string
	customCSS     string
	resize        string
	to            string
	quality       int
	fit           bool
	validateOnly  bool
	inline        bool
	verbose       bool
}

var opts options

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "acify",
		Short: "Build a detect.ac detect-theme JSON export from a background image",
		Long: "acify inlines a background image (and optional spinner image) into the\n" +
			"default HTML/CSS/JS templates and writes a single detect-theme JSON file\n" +
			"compatible with detect.ac theme import.\n\n" +
			"The generated file is not uploaded anywhere - import it into detect.ac manually.",
		Args:         cobra.NoArgs,
		RunE:         run,
		SilenceUsage: true,
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.background, "background", "b", "", "required path to the background image")
	flags.StringVarP(&opts.spinner, "spinner", "s", "", "optional path to the spinner image")
	flags.BoolVar(&opts.noSpin, "no-spin", false, "disable spinner rotation")
	flags.BoolVar(&opts.noSign, "no-sign", false, "drop the 'built with acify' signature comment")
	flags.StringVarP(&opts.out, "out", "o", ".", "directory in which to write the detect-theme JSON file")
	flags.StringVar(&opts.templates, "templates", "templates", "directory containing templates/default.{html,css,js}")
	flags.StringVar(&opts.name, "name", "acify", "theme name")
	flags.IntVar(&opts.versionNumber, "version-number", 1, "source version_number in the export")
	flags.StringVar(&opts.author, "author", "acify", "source author in the export")
	flags.IntVar(&opts.themeID, "theme-id", 0, "source theme_id in the export")
	flags.StringVar(&opts.url, "url", "", "source url in the export")
	flags.StringVar(&opts.themeType, "theme-type", "advanced", "theme type in the export")
	flags.Float64Var(&opts.overlay, "overlay", 0.4, "dark tint over the background image for readability (0.0-1.0)")
	flags.StringVar(&opts.accent, "accent", "", "override the --accent CSS color, e.g. #ff5c72")
	flags.StringVar(&opts.customCSS, "css", "", "append a custom CSS file to the bundle")
	flags.StringVar(&opts.resize, "resize", "", "downscale the background to fit within WxH, preserving aspect ratio (e.g. 1280x720)")
	flags.StringVar(&opts.to, "to", "", "re-encode the background image (png or jpeg)")
	flags.IntVar(&opts.quality, "quality", 85, "jpeg encoder quality (0-100), ignored for png")
	flags.BoolVar(&opts.fit, "fit", false, "auto-shrink the background until it fits a 256 KiB soft budget")
	flags.BoolVar(&opts.validateOnly, "validate-only", false, "validate templates, images, and output without writing a file")
	flags.BoolVarP(&opts.inline, "inline", "i", false, "print the generated CSS instead of writing a file")
	flags.BoolVarP(&opts.verbose, "verbose", "v", false, "verbose output")

	if err := cmd.MarkFlagRequired("background"); err != nil {
		panic(err)
	}
	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	if opts.overlay < 0 || opts.overlay > 1 {
		return fmt.Errorf("--overlay must be between 0.0 and 1.0")
	}
	if opts.quality < 0 || opts.quality > 100 {
		return fmt.Errorf("--quality must be between 0 and 100")
	}
	if opts.to != "" && opts.to != "png" && opts.to != "jpeg" {
		return fmt.Errorf("--to must be either png or jpeg")
	}
	var maxW, maxH int
	if opts.resize != "" {
		m := regexp.MustCompile(`^(\d+)x(\d+)$`).FindStringSubmatch(opts.resize)
		if m == nil {
			return fmt.Errorf("--resize must be in WxH form, e.g. 1280x720")
		}
		fmt.Sscanf(m[1], "%d", &maxW)
		fmt.Sscanf(m[2], "%d", &maxH)
	}

	name := strings.TrimSpace(opts.name)
	if name == "" {
		name = "acify"
	}
	versionNumber := opts.versionNumber
	if versionNumber <= 0 {
		versionNumber = 1
	}
	themeType := strings.TrimSpace(opts.themeType)
	if themeType == "" {
		themeType = "advanced"
	}

	html, css, js, err := loadTemplates(opts.templates)
	if err != nil {
		return err
	}

	procOpts := asset.ProcessOptions{MaxW: maxW, MaxH: maxH, Format: opts.to}
	if cmd.Flags().Changed("quality") {
		procOpts.Quality = opts.quality
	}
	bgURI := ""
	bgReport := ""
	if opts.fit {
		budget := validator.MaxFileBytes - len(css) - 4096
		var fits bool
		bgURI, bgReport, fits, err = asset.Fit(opts.background, budget, procOpts)
		if err != nil {
			return fmt.Errorf("background image: %w", err)
		}
		if !fits {
			fmt.Fprintf(os.Stderr, "warning: background could not be fitted under %d bytes (smallest %d B)\n",
				budget, len(bgURI))
		}
		fmt.Printf("background: fitted to %s (%d B)\n", bgReport, len(bgURI))
	} else {
		bgURI, bgReport, err = asset.Optimize(opts.background, procOpts)
		if err != nil {
			return fmt.Errorf("background image: %w", err)
		}
	}
	spURI, err := optionalImage("spinner", opts.spinner)
	if err != nil {
		return err
	}

	css = bundler.InjectCSS(css, bgURI, spURI, opts.accent, opts.overlay, !opts.noSpin)
	if opts.customCSS != "" {
		extra, err := os.ReadFile(opts.customCSS)
		if err != nil {
			return fmt.Errorf("read --css file: %w", err)
		}
		css += "\n" + string(extra)
	}

	if err := validator.CSS(css); err != nil {
		return fmt.Errorf("injected css: %w", err)
	}

	if !opts.noSign {
		html = sign("html", html)
		css = sign("css", css)
		js = sign("js", js)
	}

	exp := config.NewExport(name, versionNumber, opts.author, html, css, js)
	exp.ExportedAt = time.Now().UTC().Format(time.RFC3339)
	exp.Source.ThemeID = opts.themeID
	exp.Source.URL = opts.url
	exp.Theme.ThemeType = themeType
	if opts.accent != "" {
		exp.Theme.Tokens.Accent = opts.accent
	}
	exp.Theme.Tokens.Title = name

	if opts.validateOnly {
		data, err := bundler.MarshalExport(exp)
		if err != nil {
			return err
		}
		fmt.Printf("OK: html=%d B css=%d B js=%d B json=%d B\n", len(html), len(css), len(js), len(data))
		return nil
	}

	if opts.inline {
		fmt.Println("--- css ---")
		fmt.Println(css)
		return nil
	}

	data, err := bundler.MarshalExport(exp)
	if err != nil {
		return err
	}

	outFile := filepath.Join(opts.out, safeName(name)+"-v"+fmt.Sprint(versionNumber)+".detect-theme.json")
	if err := os.MkdirAll(filepath.Dir(outFile), 0o755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}
	if err := os.WriteFile(outFile, data, 0o644); err != nil {
		return fmt.Errorf("write %s: %w", outFile, err)
	}

	fmt.Printf("wrote %s (%d bytes: html=%d css=%d js=%d) as %s v%d\n",
		outFile, len(data), len(html), len(css), len(js), name, versionNumber)
	if opts.verbose {
		fmt.Printf("  background: %s\n", opts.background)
		if bgReport != "" {
			fmt.Printf("  background re-encoded: %s\n", bgReport)
		}
		if opts.spinner != "" {
			fmt.Printf("  spinner: %s\n", opts.spinner)
		}
		if opts.noSpin {
			fmt.Println("  no-spin: true")
		}
	}
	return nil
}

func loadTemplates(dir string) (html, css, js string, err error) {
	exts := []string{"html", "css", "js"}
	read := make(map[string][]byte, len(exts))
	for _, ext := range exts {
		p := filepath.Join(dir, "default."+ext)
		if err := validator.File(p); err != nil {
			return "", "", "", err
		}
		b, err := os.ReadFile(p)
		if err != nil {
			return "", "", "", fmt.Errorf("read %s: %w", p, err)
		}
		read[ext] = b
	}
	if err := validator.HTML(string(read["html"])); err != nil {
		return "", "", "", err
	}
	if err := validator.CSS(string(read["css"])); err != nil {
		return "", "", "", err
	}
	if err := validator.JS(string(read["js"])); err != nil {
		return "", "", "", err
	}
	return string(read["html"]), string(read["css"]), string(read["js"]), nil
}

func optionalImage(kind, path string) (string, error) {
	if path == "" {
		return "", nil
	}
	uri, err := asset.Image(path)
	if err != nil {
		return "", fmt.Errorf("%s image: %w", kind, err)
	}
	return uri, nil
}

func safeName(s string) string {
	return strings.NewReplacer(" ", "_", "/", "_", "\\", "_", ":", "_").Replace(s)
}

// repoURL identifies the project the bundle was built with.
const repoURL = "https://github.com/vlad6tz/acify"

// sign prepends a "built with acify" comment appropriate for the file kind.
func sign(kind, content string) string {
	var open, close string
	switch kind {
	case "html":
		open, close = "<!-- ", " -->"
	default:
		open, close = "/* ", " */"
	}
	return open + "built with acify - " + repoURL + close + "\n" + content
}
