package storage

import "strings"

// EscapeLikePattern экранирует специальные символы для безопасного использования в LIKE запросах
func EscapeLikePattern(pattern string) string {
	replacer := strings.NewReplacer(
		"\\", "\\\\",
		"%", "\\%",
		"_", "\\_",
	)
	return replacer.Replace(pattern)
}
