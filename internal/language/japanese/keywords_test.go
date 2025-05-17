package japanese

import (
	"testing"
)

func TestExtractJapaneseKeywords(t *testing.T) {
	text := "これはテスト用の文章です。Go言語と形態素解析を使います。"
	keywords := ExtractJapaneseKeywords(text)
	if len(keywords) == 0 {
		t.Error("expected some keywords, got none")
	}
	found := false
	for _, k := range keywords {
		if k == "言語" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected '言語' in keywords, got %v", keywords)
	}
}

func TestIsSymbolOrPunctuation(t *testing.T) {
	if !isSymbolOrPunctuation("！＠＃") {
		t.Error("expected true for symbols")
	}
	if isSymbolOrPunctuation("テスト") {
		t.Error("expected false for Japanese text")
	}
	if isSymbolOrPunctuation("abc123") {
		t.Error("expected false for alphanum")
	}
	if isSymbolOrPunctuation("") {
		t.Error("expected false for empty string")
	}
}
