package utils

import "testing"

func TestNormalizeSpace(t *testing.T) {
	in := "  foo   bar\t baz\nqux  "
	expected := "foo bar baz qux"
	out := NormalizeSpace(in)
	if out != expected {
		t.Errorf("expected '%s', got '%s'", expected, out)
	}
}

func TestUniqueStrings(t *testing.T) {
	in := []string{"a", "b", "a", "c", "b"}
	expected := []string{"a", "b", "c"}
	out := UniqueStrings(in)
	if len(out) != len(expected) {
		t.Fatalf("expected len %d, got %d", len(expected), len(out))
	}
	m := make(map[string]bool)
	for _, v := range out {
		m[v] = true
	}
	for _, v := range expected {
		if !m[v] {
			t.Errorf("missing value '%s' in output", v)
		}
	}
}
