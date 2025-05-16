package keywordextractor

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode"

	"github.com/PuerkitoBio/goquery"
	"github.com/ikawaha/kagome-dict/ipa"
	"github.com/ikawaha/kagome/v2/tokenizer"
)

// AnalyzerはURLと取得済みレスポンスボディを内部に持ち、解析・キーワード抽出を行う構造体です
type Analyzer struct {
	URL          string
	responseBody []byte
	document     *goquery.Document // キャッシュ用
}

// PageData はウェブページから抽出したさまざまなコンテンツを保持します
type PageData struct {
	Title       string
	MetaTags    map[string]string
	MainContent string
}

// NewAnalyzer はURLへアクセスし、Analyzer構造体を生成します
func NewAnalyzer(url string, timeoutSeconds int) (*Analyzer, error) {
	client := &http.Client{
		Timeout: time.Duration(timeoutSeconds) * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:       10,
			IdleConnTimeout:    30 * time.Second,
			DisableCompression: true,
		},
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("リクエストの作成に失敗しました: %w", err)
	}

	req.Header.Set("User-Agent", UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("URLにアクセスできませんでした: %w", err)
	}
	defer resp.Body.Close()

	finalURL := resp.Request.URL.String()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("レスポンスボディの読み込みに失敗しました: %w", err)
	}

	return &Analyzer{
		URL:          finalURL,
		responseBody: body,
	}, nil
}

// parseDocument はレスポンスボディからgoqueryドキュメントを取得します（キャッシュあり）
func (a *Analyzer) parseDocument() (*goquery.Document, error) {
	if a.document != nil {
		return a.document, nil
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(a.responseBody)))
	if err != nil {
		return nil, fmt.Errorf("HTMLドキュメントのパースに失敗しました: %w", err)
	}
	a.document = doc
	return doc, nil
}

// FetchTags 指定タグの内容をすべて抜き出します
func (a *Analyzer) FetchTags(tag string) ([]string, error) {
	doc, err := a.parseDocument()
	if err != nil {
		return nil, err
	}
	var result []string
	doc.Find(tag).Each(func(i int, s *goquery.Selection) {
		result = append(result, s.Text())
	})
	return result, nil
}

// FetchTitle titleタグの最初の内容を返します
func (a *Analyzer) FetchTitle() (string, error) {
	titles, err := a.FetchTags("title")
	if err != nil {
		return "", err
	}
	if len(titles) == 0 {
		return "", nil
	}
	return titles[0], nil
}

// FetchDescription descriptionとog:descriptionから、より情報量の多い方を返します
func (a *Analyzer) FetchDescription() (string, error) {
	metaTags, err := a.FetchMetaTags()
	if err != nil {
		return "", fmt.Errorf("メタタグ取得エラー: %w", err)
	}

	description, hasDesc := metaTags["description"]
	ogDescription, hasOgDesc := metaTags["og:description"]

	if !hasDesc && !hasOgDesc {
		return "", nil
	}

	if !hasDesc {
		return ogDescription, nil
	}

	if !hasOgDesc {
		return description, nil
	}

	// より情報量の多い（長い）方を採用
	if len(ogDescription) > len(description) {
		return ogDescription, nil
	}
	return description, nil
}

// FetchKeywords keywordsとog:keywordsからキーワードをユニークにマージして返します
func (a *Analyzer) FetchKeywords() ([]string, error) {
	metaTags, err := a.FetchMetaTags()
	if err != nil {
		return nil, fmt.Errorf("メタタグ取得エラー: %w", err)
	}

	keywords, hasKeywords := metaTags["keywords"]
	ogKeywords, hasOgKeywords := metaTags["og:keywords"]

	if !hasKeywords && !hasOgKeywords {
		return []string{}, nil
	}

	// キーワードをユニークにマージする
	keywordMap := make(map[string]bool)

	if hasKeywords {
		// カンマまたはセミコロンで区切られていると想定
		for _, kw := range strings.FieldsFunc(keywords, func(r rune) bool {
			return r == ',' || r == ';'
		}) {
			keywordMap[strings.TrimSpace(kw)] = true
		}
	}

	if hasOgKeywords {
		// カンマまたはセミコロンで区切られていると想定
		for _, kw := range strings.FieldsFunc(ogKeywords, func(r rune) bool {
			return r == ',' || r == ';'
		}) {
			keywordMap[strings.TrimSpace(kw)] = true
		}
	}

	// マップをスライスに変換
	result := make([]string, 0, len(keywordMap))
	for kw := range keywordMap {
		if kw != "" {
			result = append(result, kw)
		}
	}

	return result, nil
}

// FetchMetaTags metaタグを抽出します
func (a *Analyzer) FetchMetaTags() (map[string]string, error) {
	doc, err := a.parseDocument()
	if err != nil {
		return nil, fmt.Errorf("goqueryでのパースに失敗しました: %w", err)
	}
	result := make(map[string]string)

	// メタタグの属性と処理したい項目のマップ
	metaNameTargets := []string{"description", "pubdate", "keywords"}
	metaPropTargets := []string{"og:description", "og:site_name"}

	doc.Find("meta").Each(func(i int, s *goquery.Selection) {
		if name, exists := s.Attr("name"); exists {
			name = strings.ToLower(name)
			for _, target := range metaNameTargets {
				if name == target {
					if content, ok := s.Attr("content"); ok {
						result[name] = content
					}
					break
				}
			}
		}
		if prop, exists := s.Attr("property"); exists {
			prop = strings.ToLower(prop)
			for _, target := range metaPropTargets {
				if prop == target {
					if content, ok := s.Attr("content"); ok {
						result[prop] = content
					}
					break
				}
			}
		}
	})
	return result, nil
}

// FetchMainContent 見出しと本文から重要なテキストを抽出します
func (a *Analyzer) FetchMainContent() (string, error) {
	doc, err := a.parseDocument()
	if err != nil {
		return "", err
	}

	var content strings.Builder

	doc.Find("h1, h2, h3").Each(func(i int, s *goquery.Selection) {
		headingText := strings.TrimSpace(s.Text())
		if headingText != "" {
			content.WriteString(headingText + " " + headingText + " " + headingText + " ")
		}
	})

	// 本文はノイズが入ることがあるので、取得しない
	// doc.Find("p, li").Each(func(i int, s *goquery.Selection) {
	// 	content.WriteString(strings.TrimSpace(s.Text()) + " ")
	// })

	return content.String(), nil
}

// テキストが日本語を含むかどうかを判定します
func containsJapanese(text string) bool {
	for _, r := range text {
		if unicode.In(r, unicode.Hiragana, unicode.Katakana, unicode.Han) {
			return true
		}
	}
	return false
}

// 日本語テキストからキーワードを抽出します
func extractJapaneseKeywords(text string) []string {
	t, err := tokenizer.New(ipa.Dict(), tokenizer.OmitBosEos())
	if err != nil {
		return []string{}
	}

	tokens := t.Tokenize(text)
	keywordMap := make(map[string]bool)

	// 表記ゆれ対策のマップ（小文字をキーに、実際の表記を値に）
	normalizedMap := make(map[string]string)

	for _, token := range tokens {
		// 品詞情報を取得
		features := token.Features()
		if len(features) > 0 && (features[0] == "名詞") {
			// 一般名詞、固有名詞、サ変接続名詞などを抽出
			if len(features) > 1 && (features[1] == "一般" || features[1] == "固有名詞" ||
				features[1] == "サ変接続" || features[1] == "形容動詞語幹") {
				surface := token.Surface
				if len(surface) > 1 { // 1文字の名詞は除外
					// 小文字化した文字列をキーにする
					normalized := strings.ToLower(surface)
					keywordMap[normalized] = true

					// 最初に出現した表記を保持（または長さが最も長い表記を保持）
					if existing, ok := normalizedMap[normalized]; !ok || len(surface) > len(existing) {
						normalizedMap[normalized] = surface
					}
				}
			}
		}
	}

	// マップをスライスに変換（元の表記を使用）
	result := make([]string, 0, len(keywordMap))
	for norm := range keywordMap {
		if original, ok := normalizedMap[norm]; ok {
			result = append(result, original)
		} else {
			// 万が一マッピングがない場合は正規化したものを使う
			result = append(result, norm)
		}
	}

	return result
}

// ExtractKeywords テキストからキーワードを抽出し、単数・複数をまとめて返します
func ExtractKeywords(text string) []string {
	// 日本語テキストの場合は日本語向け処理
	if containsJapanese(text) {
		return extractJapaneseKeywords(text)
	}

	// 英語などの言語向け処理
	clean := strings.ToLower(text)
	clean = regexp.MustCompile(`[^\w\s-]`).ReplaceAllString(clean, " ")
	clean = regexp.MustCompile(`-{2,}`).ReplaceAllString(clean, "-")
	words := strings.Fields(clean)
	seen := map[string]int{}
	var result []string
	for _, w := range words {
		norm := normalizeKeyword(w)
		if _, skip := stopWords[norm]; !skip && len(norm) > 1 && norm != "-" {
			if seen[norm] == 0 {
				result = append(result, norm)
				seen[norm] = 1
			}
		}
	}
	return result
}

// CollectPageData はページからすべてのテキストデータを収集します
func (a *Analyzer) CollectPageData() (*PageData, error) {
	title, err := a.FetchTitle()
	if err != nil {
		return nil, fmt.Errorf("タイトル取得エラー: %w", err)
	}

	meta, err := a.FetchMetaTags()
	if err != nil {
		return nil, fmt.Errorf("メタタグ取得エラー: %w", err)
	}

	content, err := a.FetchMainContent()
	if err != nil {
		return nil, fmt.Errorf("コンテンツ取得エラー: %w", err)
	}

	return &PageData{
		Title:       title,
		MetaTags:    meta,
		MainContent: content,
	}, nil
}

// GetTopKeywords 上位N個のキーワードを取得します
func (a *Analyzer) GetTopKeywords(n int) ([]KeywordWithScore, error) {
	// 重み設定
	const (
		weightMetaKeyword = 8
		weightTitle       = 5
		weightDesc        = 3
		weightMain        = 1
	)

	// スコアマップでは小文字化したキーワードを使用
	scoreMap := map[string]int{}
	// 元の表記を保持するマップ
	originalMap := map[string]string{}

	// タイトルの処理
	title, err := a.FetchTitle()
	if err != nil {
		return nil, fmt.Errorf("タイトル取得エラー: %w", err)
	}
	if title != "" {
		for _, k := range ExtractKeywords(title) {
			normKey := strings.ToLower(k)
			scoreMap[normKey] += weightTitle
			if existing, ok := originalMap[normKey]; !ok || len(k) > len(existing) {
				originalMap[normKey] = k
			}
		}
	}

	// メタキーワードの処理
	keywords, err := a.FetchKeywords()
	if err != nil {
		return nil, fmt.Errorf("キーワード取得エラー: %w", err)
	}
	for _, k := range keywords {
		// キーワードはそのままではなく、標準化する
		for _, extractedKw := range ExtractKeywords(k) {
			normKey := strings.ToLower(extractedKw)
			scoreMap[normKey] += weightMetaKeyword
			if existing, ok := originalMap[normKey]; !ok || len(extractedKw) > len(existing) {
				originalMap[normKey] = extractedKw
			}
		}
	}

	// 説明文の処理
	description, err := a.FetchDescription()
	if err != nil {
		return nil, fmt.Errorf("説明文取得エラー: %w", err)
	}
	if description != "" {
		for _, k := range ExtractKeywords(description) {
			normKey := strings.ToLower(k)
			scoreMap[normKey] += weightDesc
			if existing, ok := originalMap[normKey]; !ok || len(k) > len(existing) {
				originalMap[normKey] = k
			}
		}
	}

	// メインコンテンツの処理
	mainContent, err := a.FetchMainContent()
	if err != nil {
		return nil, fmt.Errorf("コンテンツ取得エラー: %w", err)
	}
	if mainContent != "" {
		for _, k := range ExtractKeywords(mainContent) {
			normKey := strings.ToLower(k)
			scoreMap[normKey] += weightMain
			if existing, ok := originalMap[normKey]; !ok || len(k) > len(existing) {
				originalMap[normKey] = k
			}
		}
	}

	return rankKeywordsByScore(scoreMap, originalMap, n), nil
}

// rankKeywordsByScore はキーワードをスコア順にランク付けします
func rankKeywordsByScore(scoreMap map[string]int, originalMap map[string]string, limit int) []KeywordWithScore {
	type kv struct {
		Key   string
		Value int
	}
	var sorted []kv
	for k, v := range scoreMap {
		sorted = append(sorted, kv{k, v})
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Value > sorted[j].Value
	})

	var result []KeywordWithScore
	for _, kv := range sorted {
		// 元の表記（大文字小文字を保持）を使用
		originalKey := kv.Key
		if original, ok := originalMap[kv.Key]; ok {
			originalKey = original
		}

		result = append(result, KeywordWithScore{Keyword: originalKey, Score: kv.Value})
		if limit > 0 && len(result) >= limit {
			break
		}
	}
	return result
}

// ExtractSiteKeywords すべての情報源からキーワードを抽出します
func (a *Analyzer) ExtractSiteKeywords() ([]string, error) {
	keywords, err := a.GetTopKeywords(0) // 0は全てを取得する意味
	if err != nil {
		return nil, err
	}

	result := make([]string, len(keywords))
	for i, kw := range keywords {
		result[i] = kw.Keyword
	}
	return result, nil
}
