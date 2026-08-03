// Package validator performs syntax and sanity checks on the HTML, CSS, and
// JS that will be injected into a style.
package validator

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/net/html"
)

// MaxStyleBytes is the detect.ac per-file limit (strictly under 256 KiB).
const MaxStyleBytes = 256 * 1024

// Size reports an error when n is not strictly below the per-file limit.
func Size(kind string, n int) error {
	if n >= MaxStyleBytes {
		return fmt.Errorf("%s is %d bytes, which is not under the %d-byte (256 KiB) limit",
			kind, n, MaxStyleBytes)
	}
	return nil
}

// File checks that the file exists and is not empty.
func File(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}
	if info.IsDir() {
		return fmt.Errorf("%s: is a directory, expected a file", path)
	}
	if info.Size() == 0 {
		return fmt.Errorf("%s: file is empty", path)
	}
	return nil
}

// HTML parses the document and returns any well-formedness error.
func HTML(doc string) error {
	if err := nonEmpty("html", doc); err != nil {
		return err
	}
	_, err := html.Parse(strings.NewReader(doc))
	if err != nil {
		return fmt.Errorf("html: %w", err)
	}
	return nil
}

// CSS performs a structural check: matching braces and parentheses.
func CSS(s string) error {
	if err := nonEmpty("css", s); err != nil {
		return err
	}
	return checkBalanced("css", s)
}

// JS performs a structural check: matching (), [], and {} pairs, ignoring
// string and comment contents.
func JS(s string) error {
	if err := nonEmpty("js", s); err != nil {
		return err
	}
	return checkBalanced("js", s)
}

func nonEmpty(kind, s string) error {
	if strings.TrimSpace(s) == "" {
		return fmt.Errorf("%s: content is empty", kind)
	}
	return nil
}

// checkBalanced verifies that braces, brackets, and parentheses are balanced
// and correctly nested, skipping quoted strings and line/block comments.
func checkBalanced(kind, s string) error {
	pairs := map[rune]rune{'{': '}', '[': ']', '(': ')'}
	stack := []rune{}
	var quote rune
	comment := 0
	runes := []rune(s)
	for i := 0; i < len(runes); i++ {
		c := runes[i]
		if comment > 0 {
			if c == '/' && i+1 < len(runes) && runes[i+1] == '*' {
				comment++
				i++
			} else if c == '*' && i+1 < len(runes) && runes[i+1] == '/' {
				comment--
				i++
			}
			continue
		}
		if quote != 0 {
			if c == quote {
				quote = 0
			} else if c == '\\' && i+1 < len(runes) {
				i++
			}
			continue
		}
		switch {
		case c == '/' && i+1 < len(runes) && runes[i+1] == '/':
			for i < len(runes) && runes[i] != '\n' {
				i++
			}
		case c == '/' && i+1 < len(runes) && runes[i+1] == '*':
			comment = 1
			i++
		case c == '\'' || c == '"' || c == '`':
			quote = c
		default:
			if close, ok := pairs[c]; ok {
				stack = append(stack, close)
				continue
			}
			for j := len(stack) - 1; j >= 0; j-- {
				if stack[j] == c {
					if j != len(stack)-1 {
						return fmt.Errorf("%s: mismatched %q", kind, string(c))
					}
					stack = stack[:j]
					break
				}
			}
		}
	}
	if quote != 0 {
		return fmt.Errorf("%s: unterminated string", kind)
	}
	if comment > 0 {
		return fmt.Errorf("%s: unterminated block comment", kind)
	}
	if len(stack) != 0 {
		return fmt.Errorf("%s: unbalanced delimiter %q", kind, string(stack[len(stack)-1]))
	}
	return nil
}
