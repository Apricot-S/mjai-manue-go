package shared

import (
	"fmt"
	"path/filepath"
	"slices"
)

func GlobAll(patterns []string) ([]string, error) {
	var result []string
	for _, pattern := range patterns {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return nil, fmt.Errorf("invalid glob pattern %q: %w", pattern, err)
		}
		result = slices.Concat(result, matches)
	}
	return result, nil
}
