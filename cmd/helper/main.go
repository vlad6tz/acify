// Command helper is the acify entry point: it parses CLI flags, validates the
// templates, inlines any image assets, and bundles everything into a zip.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/vlad6tz/acify.git/internal/asset"
	"github.com/vlad6tz/acify.git/internal/bundler"
	"github.com/vlad6tz/acify.git/internal/config"
	"github.com/vlad6tz/acify.git/internal/validator"
)

func main() {
	var (
		out      = flag.String("out", "style.ac.bundle", "output bundle path")
		dir      = flag.String("templates", "templates", "directory containing templates/default.{html,css,js}")
		name     = flag.String("name", "", "theme name (default: env or 'detect')")
		version  = flag.String("version", "", "theme version (default: env or '1.0.0')")
		assetArg = flag.String("assets", "", "comma-separated image paths to inline as data URIs")
	)
	flag.Parse()

	if err := run(*dir, *name, *version, *out, *assetArg); err != nil {
		log.Fatalf("acify: %v", err)
	}
}

func run(templatesDir, name, version, out, assetList string) error {
	theme := config.Default()
	if name != "" {
		theme.Name = name
	}
	if version != "" {
		theme.Version = version
	}

	exts := []string{"html", "css", "js"}
	read := make(map[string][]byte, len(exts))
	for _, ext := range exts {
		p := filepath.Join(templatesDir, "default."+ext)
		if err := validator.File(p); err != nil {
			return err
		}
		b, err := os.ReadFile(p)
		if err != nil {
			return fmt.Errorf("read %s: %w", p, err)
		}
		read[ext] = b
	}
	if err := validator.HTML(string(read["html"])); err != nil {
		return err
	}
	if err := validator.CSS(string(read["css"])); err != nil {
		return err
	}
	if err := validator.JS(string(read["js"])); err != nil {
		return err
	}

	files := []bundler.File{
		{Name: "theme.html", Body: read["html"]},
		{Name: "theme.css", Body: read["css"]},
		{Name: "theme.js", Body: read["js"]},
	}

	if assetList != "" {
		for _, p := range strings.Split(assetList, ",") {
			p = strings.TrimSpace(p)
			if p == "" {
				continue
			}
			uri, err := asset.ImageDataURI(p)
			if err != nil {
				return fmt.Errorf("inline asset %s: %w", p, err)
			}
			theme.Assets = append(theme.Assets, config.Asset{Name: filepath.Base(p), Path: p})
			files = append(files, bundler.File{Name: filepath.Base(p), Body: []byte(uri)})
		}
	}

	data, err := bundler.Bundle(files)
	if err != nil {
		return err
	}
	if err := os.WriteFile(out, data, 0o644); err != nil {
		return fmt.Errorf("write %s: %w", out, err)
	}
	fmt.Printf("wrote %s (%d bytes, %d files) as %s v%s\n", out, len(data), len(files), theme.Name, theme.Version)
	return nil
}
