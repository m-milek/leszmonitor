package util

import (
	"regexp"
	"strings"
)

var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9\-\s]+`)
var whitespaceRegex = regexp.MustCompile(`\s+`)
var multipleHyphensRegex = regexp.MustCompile(`-+`)

func SlugFromString(s string) string {
	var slug string

	slug = strings.ToLower(s)

	slug = whitespaceRegex.ReplaceAllString(slug, "-")
	slug = nonAlphanumericRegex.ReplaceAllString(slug, "")
	slug = multipleHyphensRegex.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")

	return slug
}
