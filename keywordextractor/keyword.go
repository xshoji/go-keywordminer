package keywordextractor

import (
	"regexp"
	"sort"
	"strings"
)

// 英語のストップワードリスト
var stopWords = map[string]int{
	"the": 0, "is": 0, "are": 0, "was": 0, "were": 0, "be": 0, "been": 0, "being": 0,
	"am": 0, "i": 0, "you": 0, "he": 0, "she": 0, "it": 0, "we": 0, "they": 0,
	"of": 0, "and": 0, "or": 0, "to": 0, "in": 0, "that": 0, "have": 0, "has": 0,
	"had": 0, "with": 0, "for": 0, "on": 0, "at": 0, "by": 0, "an": 0, "a": 0,
	"as": 0, "from": 0, "but": 0, "not": 0, "this": 0, "which": 0, "will": 0,
	"would": 0, "can": 0, "could": 0, "should": 0, "do": 0, "does": 0, "did": 0,
	"so": 0, "if": 0, "about": 0, "into": 0, "than": 0, "then": 0, "them": 0,
	"their": 0, "there": 0, "these": 0, "those": 0, "such": 0, "also": 0, "just": 0,
	"up": 0, "out": 0, "over": 0, "after": 0, "before": 0, "between": 0, "because": 0,
	"while": 0, "where": 0, "when": 0, "who": 0, "whom": 0, "what": 0,
	"how": 0, "why": 0, "all": 0, "any": 0, "each": 0, "few": 0, "more": 0, "most": 0,
	"other": 0, "some": 0, "no": 0, "nor": 0, "only": 0, "own": 0, "same": 0, "too": 0,
	"very": 0, "s": 0, "t": 0, "don": 0, "now": 0,
	"here": 0, "my": 0, "your": 0, "his": 0, "her": 0, "its": 0, "our": 0,
	"mine": 0, "yours": 0, "hers": 0, "ours": 0, "theirs": 0,
	"thing": 0, "things": 0, "something": 0, "anything": 0, "everything": 0, "nothing": 0,
	"anyone": 0, "someone": 0, "everyone": 0, "none": 0, "one": 0, "ones": 0, "another": 0, "others": 0,
	"again": 0, "always": 0, "never": 0, "sometimes": 0, "often": 0, "maybe": 0, "perhaps": 0,
	"really": 0, "quite": 0, "even": 0, "still": 0, "yet": 0, "already": 0, "soon": 0,
	"today": 0, "tomorrow": 0, "yesterday": 0, "lot": 0, "lots": 0, "bit": 0, "bits": 0,
	"kind": 0, "kinds": 0, "type": 0, "types": 0, "way": 0, "ways": 0, "part": 0, "parts": 0,
	"place": 0, "places": 0, "area": 0, "areas": 0, "case": 0, "cases": 0, "example": 0, "examples": 0,
	"etc": 0, "etc.": 0,
	"well": 0, "oh": 0, "hey": 0, "hi": 0, "hello": 0, "hmm": 0, "uh": 0, "um": 0, "ah": 0, "like": 0, "okay": 0, "ok": 0, "alright": 0, "right": 0, "yeah": 0, "nope": 0, "yep": 0, "huh": 0, "hurray": 0, "oops": 0, "wow": 0, "gee": 0, "gosh": 0, "whoa": 0,
}

const UserAgent = "Mozilla/5.0 (compatible; KeywordBot/1.0)"

// 頻出する単語の単数形・複数形のマッピング辞書
var pluralSingularMap = map[string]string{
	"men": "man", "women": "woman", "children": "child",
	"teeth": "tooth", "feet": "foot", "geese": "goose", "mice": "mouse",
	"people": "person", "oxen": "ox", "leaves": "leaf", "knives": "knife",
	"lives": "life", "wolves": "wolf", "shelves": "shelf", "selves": "self",
	"sequences": "sequence", "trees": "tree", "cars": "car", "hospitals": "hospital",
	"books": "book", "phones": "phone", "houses": "house", "homes": "home",
	"schools": "school", "games": "game", "names": "name", "words": "word",
	"times": "time", "years": "year", "days": "day", "weeks": "week",
	"months": "month", "hours": "hour", "minutes": "minute", "seconds": "second",
	"businesses": "business", "companies": "company", "products": "product",
	"services": "service", "customers": "customer", "users": "user",
	"applications": "application", "systems": "system", "files": "file",
	"databases": "database", "servers": "server", "networks": "network",
	"devices": "device", "computers": "computer", "technologies": "technology",
	"industries": "industry", "markets": "market", "countries": "country",
	"cities": "city", "universities": "university", "colleges": "college",
	"students": "student", "teachers": "teacher", "doctors": "doctor",
	"patients": "patient", "engineers": "engineer", "scientists": "scientist",
	"researchers": "researcher", "developers": "developer", "designers": "designer",
	"artists": "artist", "writers": "writer", "readers": "reader",
	"viewers": "viewer", "listeners": "listener", "speakers": "speaker",
	"managers": "manager", "leaders": "leader", "employees": "employee",
}

// 複数形にしても末尾が変わらない単語
var invariantWords = map[string]bool{
	"series": true, "species": true, "deer": true, "fish": true,
	"sheep": true, "moose": true, "aircraft": true, "news": true,
	"information": true, "equipment": true, "furniture": true,
	"rice": true, "sugar": true, "water": true, "oil": true,
	"advice": true, "knowledge": true, "research": true, "data": true,
}

// normalizeKeyword は単語を正規化します（より控えめなアプローチ）
func normalizeKeyword(word string) string {
	word = strings.ToLower(word)

	// 不変の単語はそのまま
	if invariantWords[word] {
		return word
	}

	// 辞書に登録されている単語はマッピングを使用
	if singular, exists := pluralSingularMap[word]; exists {
		return singular
	}

	// 特殊なパターン変換（確実なケースのみ）
	if len(word) > 3 {
		// -ies -> -y 変換（確実なパターン）
		if strings.HasSuffix(word, "ies") && !strings.HasSuffix(word, "eies") && !strings.HasSuffix(word, "aies") {
			return word[:len(word)-3] + "y"
		}
	}

	// それ以外はそのまま返す
	return word
}

// KeywordNormalizer はキーワードの正規化戦略を提供するインターフェース
type KeywordNormalizer interface {
	NormalizeKeyword(word string) string
}

// DefaultKeywordNormalizer は基本的なキーワード正規化を提供する
type DefaultKeywordNormalizer struct{}

// NormalizeKeyword は単語を正規化します
func (n DefaultKeywordNormalizer) NormalizeKeyword(word string) string {
	return normalizeKeyword(word)
}

// KeywordWithScore はキーワードとスコアの構造体
type KeywordWithScore struct {
	Keyword string
	Score   int
}

// KeywordGroup は関連するキーワードのグループを表す
type KeywordGroup struct {
	MainKeyword string
	Related     []string
	TotalScore  int
}

// NormalizedKeyword は正規化されたキーワードの情報を保持
type NormalizedKeyword struct {
	Original string
	Normal   string
	Score    int
}

// ExtractKeywordsWithFrequency テキストからキーワードとその頻度を抽出します
func ExtractKeywordsWithFrequency(text string) []KeywordWithScore {
	clean := strings.ToLower(text)
	clean = regexp.MustCompile(`[^\w\s-]`).ReplaceAllString(clean, " ")
	clean = regexp.MustCompile(`-{2,}`).ReplaceAllString(clean, "-")
	words := strings.Fields(clean)

	// キーワードとその頻度をカウント
	wordFreq := make(map[string]int)
	normalizedWords := make(map[string][]string) // 正規化→元の単語のマッピング

	for _, w := range words {
		if _, skip := stopWords[w]; skip || len(w) <= 1 || w == "-" {
			continue
		}

		// 単語を正規化
		norm := normalizeKeyword(w)

		// ストップワードでなければカウント
		if _, skip := stopWords[norm]; !skip && len(norm) > 1 && norm != "-" {
			wordFreq[w]++

			// 正規化された形で同じになる単語をグループ化
			if norm != w {
				normalizedWords[norm] = append(normalizedWords[norm], w)
			}
		}
	}

	// 正規化された単語のスコアを集計
	normalizedScores := make(map[string]int)
	for word, freq := range wordFreq {
		norm := normalizeKeyword(word)
		normalizedScores[norm] += freq
	}

	// KeywordWithScoreのリストに変換
	var result []KeywordWithScore
	for norm, score := range normalizedScores {
		// 元の単語の中で最も頻度の高いものを代表として使用
		bestWord := norm
		bestScore := 0

		// 正規化形に対応する元の単語がある場合
		if originals, exists := normalizedWords[norm]; exists && len(originals) > 0 {
			for _, original := range originals {
				if wordFreq[original] > bestScore {
					bestWord = original
					bestScore = wordFreq[original]
				}
			}
		}

		result = append(result, KeywordWithScore{
			Keyword: bestWord,
			Score:   score,
		})
	}

	// スコアでソート
	sort.Slice(result, func(i, j int) bool {
		return result[i].Score > result[j].Score
	})

	return result
}
