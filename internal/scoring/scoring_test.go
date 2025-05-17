package scoring

import "testing"

func TestRankKeywordsByScore(t *testing.T) {
	scoreMap := map[string]int{
		"go":    10,
		"python": 5,
		"java":  7,
	}
	originalMap := map[string]string{
		"go":    "Go",
		"python": "Python",
		"java":  "Java",
	}

	result := RankKeywordsByScore(scoreMap, originalMap, 2)
	if len(result) != 2 {
		t.Errorf("expected 2 results, got %d", len(result))
	}
	if result[0].Keyword != "Go" || result[0].Score != 10 {
		t.Errorf("expected first to be Go(10), got %v", result[0])
	}
	if result[1].Keyword != "Java" || result[1].Score != 7 {
		t.Errorf("expected second to be Java(7), got %v", result[1])
	}
}

func TestRankKeywordsByScore_LimitZero(t *testing.T) {
	scoreMap := map[string]int{"a": 1, "b": 2}
	originalMap := map[string]string{"a": "A", "b": "B"}
	result := RankKeywordsByScore(scoreMap, originalMap, 0)
	if len(result) != 2 {
		t.Errorf("expected 2 results, got %d", len(result))
	}
}
