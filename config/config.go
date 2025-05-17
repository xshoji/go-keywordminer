package config

import "time"

type Config struct {
	Timeout         time.Duration
	UserAgent       string
	ScoreWeights    ScoreWeightConfig
	MaxKeywords     int
	IgnoreStopWords bool
}

type ScoreWeightConfig struct {
	Title       int
	MetaKeyword int
	Description int
	MainContent int
}

// DefaultConfig はデフォルト設定を返します
func DefaultConfig() Config {
	return Config{
		Timeout:   10 * time.Second,
		UserAgent: "Mozilla/5.0 (compatible; KeywordBot/1.0)",
		ScoreWeights: ScoreWeightConfig{
			Title:       5,
			MetaKeyword: 8,
			Description: 3,
			MainContent: 1,
		},
		MaxKeywords:     20,
		IgnoreStopWords: false,
	}
}
