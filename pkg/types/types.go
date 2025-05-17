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
	Keyword string
	Score   int
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
