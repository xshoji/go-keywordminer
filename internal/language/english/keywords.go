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

// NormalizeEnglishKeyword 英語の単語を正規化（単複変換・小文字化・invariant対応）
func NormalizeEnglishKeyword(word string, pluralSingularMap map[string]string, invariantWords map[string]bool) string {
	w := strings.ToLower(word)
	if invariantWords[w] {
		return w
	}
	if singular, ok := pluralSingularMap[w]; ok {
		return singular
	}
	// 単純なs, es, iesの変換
	if strings.HasSuffix(w, "ies") && len(w) > 3 {
		return w[:len(w)-3] + "y"
	}
	if strings.HasSuffix(w, "es") && len(w) > 2 {
		return w[:len(w)-2]
	}
	if strings.HasSuffix(w, "s") && len(w) > 1 {
		return w[:len(w)-1]
	}
	return w
}
