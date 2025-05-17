package utils

import "strings"

// NormalizeSpace は文字列の余分な空白を1つにまとめてトリムします
func NormalizeSpace(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

// UniqueStrings は重複を除いた文字列スライスを返します
func UniqueStrings(input []string) []string {
	m := make(map[string]struct{})
	var result []string
	for _, v := range input {
		if _, ok := m[v]; !ok {
			m[v] = struct{}{}
			result = append(result, v)
		}
	}
	return result
}
