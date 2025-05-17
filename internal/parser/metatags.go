package parser

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// FetchMetaTags はmetaタグから指定属性の値を抽出します
func (h *HTMLDocument) FetchMetaTags() map[string]string {
	result := make(map[string]string)

	metaNameTargets := []string{"description", "pubdate", "keywords"}
	metaPropTargets := []string{"og:description", "og:site_name"}

	h.Doc.Find("meta").Each(func(i int, s *goquery.Selection) {
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
	return result
}
