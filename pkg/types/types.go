package types

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
