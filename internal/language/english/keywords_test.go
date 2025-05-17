package english

import (
	"testing"
)

func dummyNormalize(word string) string { return word }

func TestExtractEnglishKeywords(t *testing.T) {
	stopWords := map[string]int{"the": 0, "is": 0}
	text := "The quick brown fox jumps over the lazy dog."
	keywords := ExtractEnglishKeywords(text, stopWords, dummyNormalize)
	found := false
	for _, k := range keywords {
		if k == "quick" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected 'quick' in keywords, got %v", keywords)
	}
	for _, k := range keywords {
		if k == "the" || k == "is" {
			t.Errorf("stopword '%s' should not be included", k)
		}
	}
}
