// Command acify builds detect.ac custom style bundles from a background image.
package main

import (
	"os"
)

func main() {
	if err := newRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
