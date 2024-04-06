package utils

import (
	"strings"

	"github.com/google/uuid"
	"github.com/microcosm-cc/bluemonday"
)

func TrimAndSanitizeStrict(s string) string {
	p := bluemonday.StrictPolicy()
	s = strings.TrimSpace(s)
	s = p.Sanitize(s)
	return s
}

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func GenerateUUID() (string, error) {
	uid, uiderr := uuid.NewRandom()
	if uiderr != nil {
		return "", uiderr
	}
	return uid.String(), nil
}
