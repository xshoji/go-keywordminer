package analyzer

import (
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/xshoji/go-keywordminer/config"
	"github.com/xshoji/go-keywordminer/internal/fetcher"
	"github.com/xshoji/go-keywordminer/internal/language"
	"github.com/xshoji/go-keywordminer/internal/language/english"
	"github.com/xshoji/go-keywordminer/internal/language/japanese"
	"github.com/xshoji/go-keywordminer/internal/parser"
	"github.com/xshoji/go-keywordminer/internal/scoring"
)

type PageData struct {
	Title       string
	MetaTags    map[string]string
	MainContent string
}

type Analyzer struct {
	URL          string
	responseBody []byte
	doc          *parser.HTMLDocument
	Config       config.Config
}

func NewAnalyzer(url string, cfg config.Config) (*Analyzer, error) {
	res, err := fetcher.FetchURL(url, int(cfg.Timeout.Seconds()))
	if err != nil {
		return nil, err
	}
	doc, err := parser.ParseHTMLDocument(string(res.Body))
	if err != nil {
		return nil, err
	}
	return &Analyzer{
		URL:          res.URL,
		responseBody: res.Body,
		doc:          doc,
		Config:       cfg,
	}, nil
}

// テスト用: HTML文字列からAnalyzerを生成
func NewAnalyzerFromHTML(html string, cfg config.Config) (*Analyzer, error) {
	doc, err := parser.ParseHTMLDocument(html)
	if err != nil {
		return nil, err
	}
	return &Analyzer{
		URL:          "dummy",
		responseBody: []byte(html),
		doc:          doc,
		Config:       cfg,
	}, nil
}

func (a *Analyzer) FetchTitle() (string, error) {
	titles := a.doc.FetchTags("title")
	if len(titles) == 0 {
		return "", nil
	}
	return titles[0], nil
}

func (a *Analyzer) FetchMetaTags() (map[string]string, error) {
	return a.doc.FetchMetaTags(), nil
}

func (a *Analyzer) FetchMainContent() (string, error) {
	var content string
	hTags := a.doc.FetchTags("h1")
	hTags = append(hTags, a.doc.FetchTags("h2")...)
	hTags = append(hTags, a.doc.FetchTags("h3")...)
	for _, headingText := range hTags {
		if headingText != "" {
			content += headingText + " " + headingText + " " + headingText + " "
		}
	}
	return content, nil
}

func (a *Analyzer) CollectPageData() (*PageData, error) {
	title, _ := a.FetchTitle()
	meta := a.doc.FetchMetaTags()
	content, _ := a.FetchMainContent()
	return &PageData{
		Title:       title,
		MetaTags:    meta,
		MainContent: content,
	}, nil
}

func (a *Analyzer) GetTopKeywords(n int, stopWords map[string]int, normalizeKeyword func(string) string) ([]scoring.KeywordWithScore, error) {
	cfg := a.Config
	weightMetaKeyword := cfg.ScoreWeights.MetaKeyword
	weightTitle := cfg.ScoreWeights.Title
	weightDesc := cfg.ScoreWeights.Description
	weightMain := cfg.ScoreWeights.MainContent
	if n <= 0 {
		n = cfg.MaxKeywords
	}
	scoreMap := map[string]int{}
	originalMap := map[string]string{}

	// タイトル
	title, _ := a.FetchTitle()
	if title != "" {
		for _, k := range extractKeywords(title, stopWords, normalizeKeyword) {
			normKey := k
			scoreMap[normKey] += weightTitle
			if existing, ok := originalMap[normKey]; !ok || len(k) > len(existing) {
				originalMap[normKey] = k
			}
		}
	}

	// メタキーワード
	meta := a.doc.FetchMetaTags()
	if keywords, ok := meta["keywords"]; ok {
		for _, k := range extractKeywords(keywords, stopWords, normalizeKeyword) {
			normKey := k
			scoreMap[normKey] += weightMetaKeyword
			if existing, ok := originalMap[normKey]; !ok || len(k) > len(existing) {
				originalMap[normKey] = k
			}
		}
	}

	// 説明文
	desc := ""
	if d, ok := meta["description"]; ok {
		desc = d
	}
	if d, ok := meta["og:description"]; ok && len(d) > len(desc) {
		desc = d
	}
	if desc != "" {
		for _, k := range extractKeywords(desc, stopWords, normalizeKeyword) {
			normKey := k
			scoreMap[normKey] += weightDesc
			if existing, ok := originalMap[normKey]; !ok || len(k) > len(existing) {
				originalMap[normKey] = k
			}
		}
	}

	// メインコンテンツ
	mainContent, _ := a.FetchMainContent()
	if mainContent != "" {
		for _, k := range extractKeywords(mainContent, stopWords, normalizeKeyword) {
			normKey := k
			scoreMap[normKey] += weightMain
			if existing, ok := originalMap[normKey]; !ok || len(k) > len(existing) {
				originalMap[normKey] = k
			}
		}
	}

	return scoring.RankKeywordsByScore(scoreMap, originalMap, n), nil
}

// extractKeywords: 言語自動判定して適切な抽出関数を呼ぶ
func extractKeywords(text string, stopWords map[string]int, normalizeKeyword func(string) string) []string {
	if language.ContainsJapanese(text) {
		return japanese.ExtractJapaneseKeywords(text)
	}
	return english.ExtractEnglishKeywords(text, stopWords, normalizeKeyword)
}

// ページ取得の分離
func FetchPage(url string, timeout time.Duration) (*http.Response, error) {
	client := &http.Client{Timeout: timeout}
	return client.Get(url)
}

// 文書解析の分離
func ParseDocument(body []byte) (*goquery.Document, error) {
	return goquery.NewDocumentFromReader(strings.NewReader(string(body)))
}

// キーワード抽出の分離
func ExtractKeywords(content string, isJapanese bool, stopWords map[string]int, normalizeKeyword func(string) string) ([]string, error) {
	if isJapanese || language.ContainsJapanese(content) {
		return japanese.ExtractJapaneseKeywords(content), nil
	}
	return english.ExtractEnglishKeywords(content, stopWords, normalizeKeyword), nil
}
