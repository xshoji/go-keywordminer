package english

import (
	"regexp"
	"sort"
	"strings"
)

// ExtractEnglishKeywords 英語テキストからキーワードを抽出（頻度順、正規化、代表単語選択）
func ExtractEnglishKeywords(text string, stopWords map[string]int, normalizeKeyword func(string) string) []string {
	clean := strings.ToLower(text)
	clean = regexp.MustCompile(`[^\w\s-]`).ReplaceAllString(clean, " ")
	clean = regexp.MustCompile(`-{2,}`).ReplaceAllString(clean, "-")
	words := strings.Fields(clean)

	wordFreq := make(map[string]int)
	normalizedWords := make(map[string][]string) // 正規化→元の単語のマッピング

	for _, w := range words {
		if _, skip := stopWords[w]; skip || len(w) <= 1 || w == "-" {
			continue
		}
		norm := normalizeKeyword(w)
		if _, skip := stopWords[norm]; !skip && len(norm) > 1 && norm != "-" {
			wordFreq[w]++
			if norm != w {
				normalizedWords[norm] = append(normalizedWords[norm], w)
			}
		}
	}

	normalizedScores := make(map[string]int)
	for word, freq := range wordFreq {
		norm := normalizeKeyword(word)
		normalizedScores[norm] += freq
	}

	type keywordWithScore struct {
		Keyword string
		Score   int
	}
	var resultList []keywordWithScore
	for norm, score := range normalizedScores {
		bestWord := norm
		bestScore := 0
		if originals, exists := normalizedWords[norm]; exists && len(originals) > 0 {
			for _, original := range originals {
				if wordFreq[original] > bestScore {
					bestWord = original
					bestScore = wordFreq[original]
				}
			}
		}
		resultList = append(resultList, keywordWithScore{
			Keyword: bestWord,
			Score:   score,
		})
	}

	sort.Slice(resultList, func(i, j int) bool {
		return resultList[i].Score > resultList[j].Score
	})

	var result []string
	for _, kw := range resultList {
		result = append(result, kw.Keyword)
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
