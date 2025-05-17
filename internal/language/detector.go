package language

import (
	"unicode"
)

// ContainsJapanese はテキストが日本語（ひらがな・カタカナ・漢字）を含むか判定
func ContainsJapanese(text string) bool {
	for _, r := range text {
		if unicode.In(r, unicode.Hiragana, unicode.Katakana, unicode.Han) {
			return true
		}
	}
	return false
}

// IsHiragana はテキストがひらがなのみで構成されているか判定
func IsHiragana(text string) bool {
	for _, r := range text {
		if !unicode.In(r, unicode.Hiragana) {
			return false
		}
	}
	return len(text) > 0 // 空文字列はfalse
}

// IsSymbolOrPunctuation はテキストが記号や特殊文字のみかどうかを判定
func IsSymbolOrPunctuation(text string) bool {
	if text == "" {
		return false
	}
	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsDigit(r) ||
			unicode.In(r, unicode.Hiragana, unicode.Katakana, unicode.Han) {
			return false
		}
	}
	return true
}
