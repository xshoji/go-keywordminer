package types

import "testing"

func TestPageDataStruct(t *testing.T) {
	meta := map[string]string{"description": "desc"}
	pd := PageData{
		Title:       "title",
		MetaTags:    meta,
		MainContent: "main",
	}
	if pd.Title != "title" || pd.MetaTags["description"] != "desc" || pd.MainContent != "main" {
		t.Errorf("PageData struct fields not set or retrieved correctly: %+v", pd)
	}
}

func TestKeywordWithScoreStruct(t *testing.T) {
	kws := KeywordWithScore{Keyword: "go", Score: 10}
	if kws.Keyword != "go" || kws.Score != 10 {
		t.Errorf("KeywordWithScore struct fields not set or retrieved correctly: %+v", kws)
	}
}
