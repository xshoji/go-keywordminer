package scoring

import (
	"sort"
)

type KeywordWithScore struct {
	Keyword string
	Score   int
}

// RankKeywordsByScore はキーワードをスコア順にランク付けします
func RankKeywordsByScore(scoreMap map[string]int, originalMap map[string]string, limit int) []KeywordWithScore {
	type kv struct {
		Key   string
		Value int
	}
	var sorted []kv
	for k, v := range scoreMap {
		sorted = append(sorted, kv{k, v})
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Value > sorted[j].Value
	})

	var result []KeywordWithScore
	for _, kv := range sorted {
		originalKey := kv.Key
		if original, ok := originalMap[kv.Key]; ok {
			originalKey = original
		}
		result = append(result, KeywordWithScore{Keyword: originalKey, Score: kv.Value})
		if limit > 0 && len(result) >= limit {
			break
		}
	}
	return result
}
