package types

import (
	"strings"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// --- ダミー実装 ---
type dummyFetcher struct{}

func (d dummyFetcher) Fetch(url string, timeout time.Duration) ([]byte, error) {
	return []byte("dummy body"), nil
}

type dummyExtractor struct{}

func (d dummyExtractor) Extract(text string) ([]KeywordWithScore, error) {
	return []KeywordWithScore{{Keyword: "go", Score: 1}}, nil
}

type dummyParser struct{}

func (d dummyParser) ParseTitle(doc *goquery.Document) (string, error) { return "title", nil }
func (d dummyParser) ParseMetaTags(doc *goquery.Document) (map[string]string, error) {
	return map[string]string{"k": "v"}, nil
}
func (d dummyParser) ParseMainContent(doc *goquery.Document) (string, error) { return "main", nil }

func TestPageFetcherInterface(t *testing.T) {
	var f PageFetcher = dummyFetcher{}
	b, err := f.Fetch("http://example.com", 1*time.Second)
	if err != nil || string(b) != "dummy body" {
		t.Errorf("PageFetcher interface not working: %v, %s", err, string(b))
	}
}

func TestKeywordExtractorInterface(t *testing.T) {
	var e KeywordExtractor = dummyExtractor{}
	kws, err := e.Extract("text")
	if err != nil || len(kws) != 1 || kws[0].Keyword != "go" {
		t.Errorf("KeywordExtractor interface not working: %v, %+v", err, kws)
	}
}

func TestDocumentParserInterface(t *testing.T) {
	var p DocumentParser = dummyParser{}
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader("<html><head><title>t</title></head><body></body></html>"))
	title, _ := p.ParseTitle(doc)
	if title != "title" {
		t.Errorf("DocumentParser.ParseTitle not working: %s", title)
	}
	meta, _ := p.ParseMetaTags(doc)
	if meta["k"] != "v" {
		t.Errorf("DocumentParser.ParseMetaTags not working: %+v", meta)
	}
	main, _ := p.ParseMainContent(doc)
	if main != "main" {
		t.Errorf("DocumentParser.ParseMainContent not working: %s", main)
	}
}
