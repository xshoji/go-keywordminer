package parser

import (
	"testing"
)

func TestParseHTMLDocument_Success(t *testing.T) {
	html := `<html><head><title>Test</title></head><body><h1>Hello</h1></body></html>`
	doc, err := ParseHTMLDocument(html)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if doc == nil || doc.Doc == nil {
		t.Fatal("doc or doc.Doc is nil")
	}

	tags := doc.FetchTags("h1")
	if len(tags) != 1 || tags[0] != "Hello" {
		t.Errorf("expected [Hello], got %v", tags)
	}
}

func TestParseHTMLDocument_InvalidHTML(t *testing.T) {
	_, err := ParseHTMLDocument("")
	if err != nil {
		t.Error("goquery should not error on empty string")
	}
}
