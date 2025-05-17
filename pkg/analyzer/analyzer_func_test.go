package analyzer

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestFetchPage(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	}))
	defer ts.Close()

	resp, err := FetchPage(ts.URL, 2*time.Second)
	if err != nil {
		t.Fatalf("FetchPage error: %v", err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if string(body) != "hello world" {
		t.Errorf("expected 'hello world', got '%s'", string(body))
	}
}

func TestParseDocument(t *testing.T) {
	html := []byte(`<html><head><title>abc</title></head><body><h1>hoge</h1></body></html>`)
	doc, err := ParseDocument(html)
	if err != nil {
		t.Fatalf("ParseDocument error: %v", err)
	}
	title := doc.Find("title").Text()
	if title != "abc" {
		t.Errorf("expected 'abc', got '%s'", title)
	}
}

func TestExtractKeywords_English(t *testing.T) {
	text := "Go is awesome. Go is fast."
	stopWords := map[string]int{"is": 0}
	keywords, err := ExtractKeywords(text, false, stopWords, func(s string) string { return s })
	if err != nil {
		t.Fatalf("ExtractKeywords error: %v", err)
	}
	found := false
	for _, k := range keywords {
		if k == "go" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected 'go' in keywords, got %v", keywords)
	}
}

func TestExtractKeywords_Japanese(t *testing.T) {
	text := "これはテストです。Go言語を使います。"
	keywords, err := ExtractKeywords(text, true, map[string]int{}, func(s string) string { return s })
	if err != nil {
		t.Fatalf("ExtractKeywords error: %v", err)
	}
	if len(keywords) == 0 {
		t.Error("expected some keywords, got none")
	}
}
