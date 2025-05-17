package parser

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// HTMLDocument は goquery.Document のラッパー
// 今後の拡張やテスト容易性のための型
// （必要に応じて拡張可）
type HTMLDocument struct {
	Doc *goquery.Document
}

// ParseHTMLDocument はHTML文字列から goquery.Document を生成します
func ParseHTMLDocument(html string) (*HTMLDocument, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, fmt.Errorf("Failed to parse HTML document: %w", err)
	}
	return &HTMLDocument{Doc: doc}, nil
}

// FetchTags 指定タグの内容をすべて抜き出します
func (h *HTMLDocument) FetchTags(tag string) []string {
	var result []string
	h.Doc.Find(tag).Each(func(i int, s *goquery.Selection) {
		result = append(result, s.Text())
	})
	return result
}
