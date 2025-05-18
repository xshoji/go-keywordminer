package config

import (
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Timeout != 10*time.Second {
		t.Errorf("expected Timeout 10s, got %v", cfg.Timeout)
	}
	if cfg.UserAgent == "" {
		t.Error("UserAgent should not be empty")
	}
	if cfg.ScoreWeights.Title != 5 || cfg.ScoreWeights.MetaKeyword != 8 || cfg.ScoreWeights.Description != 3 || cfg.ScoreWeights.MainContent != 1 {
		t.Errorf("unexpected ScoreWeights: %+v", cfg.ScoreWeights)
	}
	if cfg.MaxKeywords != 20 {
		t.Errorf("expected MaxKeywords 20, got %d", cfg.MaxKeywords)
	}
}
