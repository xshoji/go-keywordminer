package parser

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestFetchMetaTags(t *testing.T) {
	html := `<html><head>
	<meta name="description" content="desc1">
	<meta name="keywords" content="go, test">
	<meta property="og:description" content="ogdesc">
	<meta property="og:site_name" content="sitename">
	</head></html>`
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	h := &HTMLDocument{Doc: doc}
	meta := h.FetchMetaTags()
	if meta["description"] != "desc1" {
		t.Errorf("expected desc1, got %s", meta["description"])
	}
	if meta["keywords"] != "go, test" {
		t.Errorf("expected 'go, test', got %s", meta["keywords"])
	}
	if meta["og:description"] != "ogdesc" {
		t.Errorf("expected ogdesc, got %s", meta["og:description"])
	}
	if meta["og:site_name"] != "sitename" {
		t.Errorf("expected sitename, got %s", meta["og:site_name"])
	}
}
