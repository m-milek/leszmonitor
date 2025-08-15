package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIdFromString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello World", "hello-world"},
		{"already-hyphenated", "already-hyphenated"},
		{"Mix-ed Hyphen", "mix-ed-hyphen"},
		{"Special@Chars#Here!", "specialcharshere"},
		{"  Trim  Spaces  ", "trim-spaces"},
		{"Multiple   Spaces", "multiple-spaces"},
		{"UPPERCASE", "uppercase"},
		{"123-numbers-456", "123-numbers-456"},
		{"", ""},
		{"---", ""},
		{"Test - With - Dashes", "test-with-dashes"},
		{"consecutive---hyphens", "consecutive-hyphens"},
		{"@#$%^&*()", ""},
		{"   ", ""},
		{"-leading-hyphen", "leading-hyphen"},
		{"trailing-hyphen-", "trailing-hyphen"},
		{"one-two--three---four", "one-two-three-four"},
		{"CamelCase", "camelcase"},
		{"snake_case", "snakecase"},
		{"dot.separated.words", "dotseparatedwords"},
		{"email@example.com", "emailexamplecom"},
		{"tabs\there\ttoo", "tabs-here-too"},
		{"newline\ntest", "newline-test"},
		{"!start-middle-end!", "start-middle-end"},
		{"a", "a"},
		{"1", "1"},
		{"-", ""},
		{"test—with—em—dash", "testwithemdash"}, // em dash is removed
		{"test–with–en–dash", "testwithendash"}, // en dash is removed
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := IdFromString(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
