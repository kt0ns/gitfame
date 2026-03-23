package filter

import (
	"path/filepath"
)

type Condition func(filePath string) bool

func ChainFilters(files []string, conditions ...Condition) []string {
	var result []string
	for _, file := range files {
		keep := true
		for _, cond := range conditions {
			if !cond(file) {
				keep = false
				break
			}
		}
		if keep {
			result = append(result, file)
		}
	}
	return result
}

func ByExtensions(allowedExts []string) Condition {
	if len(allowedExts) == 0 {
		return func(filePath string) bool { return true }
	}

	extMap := make(map[string]bool, len(allowedExts))
	for _, ext := range allowedExts {
		extMap[ext] = true
	}

	return func(filePath string) bool {
		return extMap[filepath.Ext(filePath)]
	}
}

func ExcludePatterns(patterns []string) Condition {
	if len(patterns) == 0 {
		return func(filePath string) bool { return true }
	}

	return func(filePath string) bool {
		for _, pattern := range patterns {
			matched, err := filepath.Match(pattern, filePath)
			if err == nil && matched {
				return false // совпало с Exclude -> отбрасываем
			}
		}
		return true
	}
}

func RestrictToPatterns(patterns []string) Condition {
	if len(patterns) == 0 {
		return func(filePath string) bool { return true }
	}

	return func(filePath string) bool {
		for _, pattern := range patterns {
			matched, err := filepath.Match(pattern, filePath)
			if err == nil && matched {
				return true // совпало -> оставляем
			}
		}
		return false
	}
}
