package domain

import (
	"strings"
	"unicode"
)

func NormalizeDocument(value string) string {
	var b strings.Builder
	for _, r := range value {
		if unicode.IsDigit(r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func MaskDocument(document string) string {
	normalized := NormalizeDocument(document)
	switch len(normalized) {
	case 11:
		return "***." + normalized[3:6] + "." + normalized[6:9] + "-**"
	default:
		return ""
	}
}
