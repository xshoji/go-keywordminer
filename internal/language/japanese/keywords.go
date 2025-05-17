package japanese

import (
	"strings"
	"unicode"

	"github.com/ikawaha/kagome-dict/ipa"
	"github.com/ikawaha/kagome/v2/tokenizer"
)

// ExtractJapaneseKeywords 日本語テキストからキーワードを抽出
func ExtractJapaneseKeywords(text string) []string {
	t, err := tokenizer.New(ipa.Dict(), tokenizer.OmitBosEos())
	if err != nil {
		return []string{}
	}
	tokens := t.Tokenize(text)
	keywordMap := make(map[string]bool)
	normalizedMap := make(map[string]string)
	for _, token := range tokens {
		features := token.Features()
		if len(features) == 0 || features[0] != "名詞" {
			continue
		}
		if len(features) <= 1 || !(features[1] == "一般" || features[1] == "固有名詞" || features[1] == "サ変接続" || features[1] == "形容動詞語幹") {
			continue
		}
		surface := token.Surface
		runes := []rune(surface)
		charLength := len(runes)
		if charLength == 1 && !unicode.In(runes[0], unicode.Han) {
			continue
		}
		if isSymbolOrPunctuation(surface) {
			continue
		}
		normalized := strings.ToLower(surface)
		keywordMap[normalized] = true
		if existing, ok := normalizedMap[normalized]; !ok || len(surface) > len(existing) {
			normalizedMap[normalized] = surface
		}
	}
	result := make([]string, 0, len(keywordMap))
	for norm := range keywordMap {
		if original, ok := normalizedMap[norm]; ok {
			result = append(result, original)
		} else {
			result = append(result, norm)
		}
	}
	return result
}

// isSymbolOrPunctuation 日本語用: 記号や特殊文字のみか判定
func isSymbolOrPunctuation(text string) bool {
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
