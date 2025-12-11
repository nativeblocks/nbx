package errors

import (
	"strings"
	"testing"

	"github.com/nativeblocks/nbx/internal/lexer"
)

func TestErrorFormat(t *testing.T) {
	err := &NBXError{
		Severity:   SeverityError,
		Message:    "Expected '=' after property name",
		Line:       10,
		Column:     15,
		SourceLine: "    fontSize \"24\"",
		Suggestion: "Add '=' between the property name and value",
	}

	formatted := err.Format()

	if !strings.Contains(formatted, "Error:") {
		t.Error("Formatted error should contain 'Error:'")
	}
	if !strings.Contains(formatted, "line 10, column 15") {
		t.Error("Formatted error should contain line and column")
	}
	if !strings.Contains(formatted, "fontSize") {
		t.Error("Formatted error should contain source line")
	}
	if !strings.Contains(formatted, "Suggestion:") {
		t.Error("Formatted error should contain suggestion")
	}
}

func TestErrorCollector(t *testing.T) {
	source := "var x: INT = 10\nvar y: STRING = hello"
	ec := NewErrorCollector(source)

	ec.AddSimpleError("Test error", 1, 5)

	if !ec.HasErrors() {
		t.Error("ErrorCollector should have errors")
	}

	if len(ec.Errors()) != 1 {
		t.Errorf("Expected 1 error, got %d", len(ec.Errors()))
	}

	ec.AddWarning("Test warning", 2, 10, "Fix this")

	if !ec.HasWarnings() {
		t.Error("ErrorCollector should have warnings")
	}

	if len(ec.Warnings()) != 1 {
		t.Errorf("Expected 1 warning, got %d", len(ec.Warnings()))
	}

	err := ec.Errors()[0]
	if err.SourceLine == "" {
		t.Error("Error should have source line populated")
	}
}

func TestUndefinedVariableError(t *testing.T) {
	availableVars := []string{"count", "counter", "value"}
	err := UndefinedVariableError("counts", 10, 5, availableVars)

	if !strings.Contains(err.Message, "Undefined variable") {
		t.Error("Error message should mention undefined variable")
	}

	if !strings.Contains(err.Suggestion, "count") {
		t.Errorf("Should suggest 'count', got: %s", err.Suggestion)
	}
}

func TestTypeMismatchError(t *testing.T) {
	err := TypeMismatchError("INT", "STRING", 15, 20)

	if !strings.Contains(err.Message, "Type mismatch") {
		t.Error("Error message should mention type mismatch")
	}
	if !strings.Contains(err.Message, "INT") || !strings.Contains(err.Message, "STRING") {
		t.Error("Error message should mention both types")
	}
}

func TestDuplicateDeclarationError(t *testing.T) {
	err := DuplicateDeclarationError("myVar", 20, 5, 10)

	if !strings.Contains(err.Message, "Duplicate") {
		t.Error("Error message should mention duplicate")
	}
	if len(err.RelatedInfo) == 0 {
		t.Error("Should have related info about first declaration")
	}
}

func TestUnknownAttributeError(t *testing.T) {
	validAttrs := []string{"name", "route", "type"}
	err := UnknownAttributeError("nam", "frame", 5, 10, validAttrs)

	if !strings.Contains(err.Message, "Unknown attribute") {
		t.Error("Error message should mention unknown attribute")
	}

	// Should suggest similar attribute
	if !strings.Contains(err.Suggestion, "name") {
		t.Errorf("Should suggest 'name', got: %s", err.Suggestion)
	}
}

func TestLevenshtein(t *testing.T) {
	tests := []struct {
		s1       string
		s2       string
		expected int
	}{
		{"", "", 0},
		{"a", "", 1},
		{"", "a", 1},
		{"abc", "abc", 0},
		{"abc", "abd", 1},
		{"abc", "adc", 1},
		{"kitten", "sitting", 3},
	}

	for _, tt := range tests {
		result := levenshtein(tt.s1, tt.s2)
		if result != tt.expected {
			t.Errorf("levenshtein(%q, %q) = %d, want %d", tt.s1, tt.s2, result, tt.expected)
		}
	}
}

func TestFindSimilar(t *testing.T) {
	candidates := []string{"fontSize", "fontWeight", "fontStyle", "backgroundColor"}

	tests := []struct {
		target   string
		expected string
	}{
		{"fontSiz", "fontSize"},
		{"fontsize", "fontSize"},
		{"fontwieght", "fontWeight"},
	}

	for _, tt := range tests {
		similar := findSimilar(tt.target, candidates)
		if len(similar) == 0 {
			t.Errorf("findSimilar(%q) returned no results", tt.target)
			continue
		}
		if similar[0] != tt.expected {
			t.Errorf("findSimilar(%q) = %q, want %q", tt.target, similar[0], tt.expected)
		}
	}
}

func TestUnexpectedTokenError(t *testing.T) {
	expected := lexer.Token{Type: lexer.TOKEN_ASSIGN, Literal: "=", Line: 1, Column: 10}
	got := lexer.Token{Type: lexer.TOKEN_COMMA, Literal: ",", Line: 1, Column: 10}

	err := UnexpectedTokenError(expected, got)

	if !strings.Contains(err.Message, "Expected") {
		t.Error("Error message should say 'Expected'")
	}
	if !strings.Contains(err.Suggestion, "Try replacing") {
		t.Error("Should provide replacement suggestion")
	}
}
