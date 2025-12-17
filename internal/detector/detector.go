package detector

import (
	"regexp"
	"strings"
)

const (
	FormatDSL     = "dsl"
	FormatXML     = "xml"
	FormatUnknown = "unknown"
)

var (
	xmlPattern = regexp.MustCompile(`^\s*<\?xml`)
	dslPattern = regexp.MustCompile(`^\s*frame\s*\(`)
)

// DetectFormat detects whether the input is DSL, XML, or unknown format.
// It returns one of: "dsl", "xml", or "unknown"
func DetectFormat(content string) string {
	if content == "" {
		return FormatUnknown
	}

	trimmed := strings.TrimSpace(content)

	if xmlPattern.MatchString(trimmed) {
		return FormatXML
	}
	if strings.HasPrefix(trimmed, "<frame") {
		return FormatXML
	}

	if dslPattern.MatchString(trimmed) {
		return FormatDSL
	}

	return FormatUnknown
}
