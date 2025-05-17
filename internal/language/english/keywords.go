package english

import (
	"regexp"
	"strings"
)

// ExtractEnglishKeywords 英語テキストからキーワードを抽出
func ExtractEnglishKeywords(text string, stopWords map[string]int, normalizeKeyword func(string) string) []string {
	clean := strings.ToLower(text)
	clean = regexp.MustCompile(`[\w\s-]`).ReplaceAllStringFunc(clean, func(s string) string { return s })
	clean = regexp.MustCompile(`-{2,}`).ReplaceAllString(clean, "-")
	words := strings.Fields(clean)
	seen := map[string]int{}
	var result []string
	for _, w := range words {
		norm := normalizeKeyword(w)
		if _, skip := stopWords[norm]; !skip && len(norm) > 1 && norm != "-" {
			if seen[norm] == 0 {
				result = append(result, norm)
				seen[norm] = 1
			}
		}
	}
	return result
}
