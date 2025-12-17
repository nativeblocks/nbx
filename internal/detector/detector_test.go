package detector

import "testing"

func TestDetectFormat(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name:     "Empty string",
			content:  "",
			expected: FormatUnknown,
		},
		{
			name:     "XML with declaration",
			content:  `<frame name="test" route="/test"></frame>`,
			expected: FormatXML,
		},
		{
			name:     "XML without declaration",
			content:  `<frame name="test" route="/test"></frame>`,
			expected: FormatXML,
		},
		{
			name:     "XML with leading whitespace",
			content:  `  <frame name="test" route="/test"></frame>`,
			expected: FormatXML,
		},
		{
			name:     "DSL format",
			content:  `frame(name = "test", route = "/test") {}`,
			expected: FormatDSL,
		},
		{
			name:     "DSL with leading whitespace",
			content:  `  frame(name = "test", route = "/test") {}`,
			expected: FormatDSL,
		},
		{
			name:     "DSL with newlines",
			content:  "\nframe(name = \"test\", route = \"/test\") {}",
			expected: FormatDSL,
		},
		{
			name:     "Unknown format",
			content:  "random text",
			expected: FormatUnknown,
		},
		{
			name:     "Only whitespace",
			content:  "   \n\t  ",
			expected: FormatUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DetectFormat(tt.content)
			if result != tt.expected {
				t.Errorf("DetectFormat() = %v, expected %v", result, tt.expected)
			}
		})
	}
}
