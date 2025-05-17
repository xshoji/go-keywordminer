package language

import "testing"

func TestContainsJapanese(t *testing.T) {
	jp := "これはテストです"
	en := "This is a test"
	if !ContainsJapanese(jp) {
		t.Error("expected true for Japanese text")
	}
	if ContainsJapanese(en) {
		t.Error("expected false for English text")
	}
}

func TestIsHiragana(t *testing.T) {
	if !IsHiragana("あいうえお") {
		t.Error("expected true for hiragana")
	}
	if IsHiragana("アイウエオ") {
		t.Error("expected false for katakana")
	}
	if IsHiragana("abc") {
		t.Error("expected false for latin")
	}
	if IsHiragana("") {
		t.Error("expected false for empty string")
	}
}

func TestIsSymbolOrPunctuation(t *testing.T) {
	if !IsSymbolOrPunctuation("!@#") {
		t.Error("expected true for symbols")
	}
	if IsSymbolOrPunctuation("あいう") {
		t.Error("expected false for hiragana")
	}
	if IsSymbolOrPunctuation("abc123") {
		t.Error("expected false for alphanum")
	}
	if IsSymbolOrPunctuation("") {
		t.Error("expected false for empty string")
	}
}
