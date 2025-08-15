package util

import (
	"github.com/m-milek/leszmonitor/logging"
	"regexp"
	"strings"
)

var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9\-\s]+`)
var whitespaceRegex = regexp.MustCompile(`\s+`)
var multipleHyphensRegex = regexp.MustCompile(`-+`)

func IdFromString(s string) string {
	var id string

	id = strings.ToLower(s)

	id = whitespaceRegex.ReplaceAllString(id, "-")
	id = nonAlphanumericRegex.ReplaceAllString(id, "")
	id = multipleHyphensRegex.ReplaceAllString(id, "-") // Collapse multiple hyphens
	id = strings.Trim(id, "-")

	logging.Init.Trace().Str("id", id).Str("source", s).Msg("Generated ID from string")
	return id
}
