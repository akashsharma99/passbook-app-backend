package utils

import (
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

func TrimAndSanitizeStrict(s string) string {
	p := bluemonday.StrictPolicy()
	s = strings.TrimSpace(s)
	s = p.Sanitize(s)
	return s
}
