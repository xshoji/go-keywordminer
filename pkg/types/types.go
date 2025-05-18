package types

import (
	"time"

	"github.com/PuerkitoBio/goquery"
)

// PageData はウェブページから抽出したさまざまなコンテンツを保持します
type PageData struct {
	Title       string
	MetaTags    map[string]string
	MainContent string
}

// KeywordWithScore はキーワードとスコアの構造体
type KeywordWithScore struct {
	Keyword string `json:"keyword"`
	Score   int    `json:"score"`
}

// AnalysisResult はウェブページの解析結果を表す構造体
type AnalysisResult struct {
	Title    string             `json:"title,omitempty"`
	MetaTags map[string]string  `json:"meta_tags,omitempty"`
	Keywords []KeywordWithScore `json:"keywords,omitempty"`
}

// PageFetcher: ページ取得のインターフェース
type PageFetcher interface {
	Fetch(url string, timeout time.Duration) ([]byte, error)
}

// KeywordExtractor: キーワード抽出のインターフェース
type KeywordExtractor interface {
	Extract(text string) ([]KeywordWithScore, error)
}

// DocumentParser: 文書解析のインターフェース
type DocumentParser interface {
	ParseTitle(doc *goquery.Document) (string, error)
	ParseMetaTags(doc *goquery.Document) (map[string]string, error)
	ParseMainContent(doc *goquery.Document) (string, error)
}
