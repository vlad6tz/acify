// Package bundler builds detect-theme JSON exports.
package bundler

import (
	"encoding/json"
	"fmt"

	"github.com/vlad6tz/acify.git/internal/config"
)

// MarshalExport serializes a detect-theme export as pretty-printed JSON.
func MarshalExport(exp config.Export) ([]byte, error) {
	data, err := json.MarshalIndent(exp, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal detect-theme export: %w", err)
	}
	return append(data, '\n'), nil
}
