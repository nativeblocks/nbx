package errors

import (
	"fmt"
	"strings"

	"github.com/nativeblocks/nbx/internal/lexer"
)

type ErrorSeverity int

const (
	SeverityError ErrorSeverity = iota
	SeverityWarning
)

func (s ErrorSeverity) String() string {
	switch s {
	case SeverityError:
		return "Error"
	case SeverityWarning:
		return "Warning"
	default:
		return "Unknown"
	}
}

type NBXError struct {
	Severity    ErrorSeverity
	Message     string
	Line        int
	Column      int
	SourceLine  string
	Suggestion  string
	RelatedInfo []string
	Token       *lexer.Token
}

func (e *NBXError) Format() string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("%s: %s\n", e.Severity, e.Message))

	if e.Line > 0 {
		b.WriteString(fmt.Sprintf("  --> line %d, column %d\n", e.Line, e.Column))
	}

	if e.SourceLine != "" {
		b.WriteString(fmt.Sprintf("    |\n"))
		b.WriteString(fmt.Sprintf("%3d | %s\n", e.Line, e.SourceLine))

		if e.Column > 0 {
			pointer := strings.Repeat(" ", e.Column-1) + "^"
			if e.Token != nil && len(e.Token.Literal) > 1 {
				pointer += strings.Repeat("~", len(e.Token.Literal)-1)
			}
			b.WriteString(fmt.Sprintf("    | %s\n", pointer))
		}
		b.WriteString(fmt.Sprintf("    |\n"))
	}

	if e.Suggestion != "" {
		b.WriteString(fmt.Sprintf("\nSuggestion: %s\n", e.Suggestion))
	}

	if len(e.RelatedInfo) > 0 {
		b.WriteString("\n")
		for _, info := range e.RelatedInfo {
			b.WriteString(fmt.Sprintf("Note: %s\n", info))
		}
	}

	return b.String()
}

func (e *NBXError) String() string {
	return e.Format()
}

type ErrorCollector struct {
	errors   []*NBXError
	warnings []*NBXError
	source   string
}

func NewErrorCollector(source string) *ErrorCollector {
	return &ErrorCollector{
		errors:   make([]*NBXError, 0),
		warnings: make([]*NBXError, 0),
		source:   source,
	}
}

func (ec *ErrorCollector) AddError(err *NBXError) {
	if err.SourceLine == "" && ec.source != "" && err.Line > 0 {
		err.SourceLine = ec.getSourceLine(err.Line)
	}

	if err.Severity == SeverityWarning {
		ec.warnings = append(ec.warnings, err)
	} else {
		ec.errors = append(ec.errors, err)
	}
}

func (ec *ErrorCollector) AddSimpleError(message string, line, column int) {
	ec.AddError(&NBXError{
		Severity: SeverityError,
		Message:  message,
		Line:     line,
		Column:   column,
	})
}

func (ec *ErrorCollector) AddTokenError(message string, token lexer.Token, suggestion string) {
	ec.AddError(&NBXError{
		Severity:   SeverityError,
		Message:    message,
		Line:       token.Line,
		Column:     token.Column,
		Token:      &token,
		Suggestion: suggestion,
	})
}

func (ec *ErrorCollector) AddWarning(message string, line, column int, suggestion string) {
	ec.AddError(&NBXError{
		Severity:   SeverityWarning,
		Message:    message,
		Line:       line,
		Column:     column,
		Suggestion: suggestion,
	})
}

func (ec *ErrorCollector) HasErrors() bool {
	return len(ec.errors) > 0
}

func (ec *ErrorCollector) HasWarnings() bool {
	return len(ec.warnings) > 0
}

func (ec *ErrorCollector) Errors() []*NBXError {
	return ec.errors
}

func (ec *ErrorCollector) Warnings() []*NBXError {
	return ec.warnings
}

func (ec *ErrorCollector) AllIssues() []*NBXError {
	all := make([]*NBXError, 0, len(ec.errors)+len(ec.warnings))
	all = append(all, ec.errors...)
	all = append(all, ec.warnings...)
	return all
}

func (ec *ErrorCollector) FormatAll() string {
	var b strings.Builder

	if len(ec.errors) > 0 {
		b.WriteString(fmt.Sprintf("Found %d error(s):\n\n", len(ec.errors)))
		for i, err := range ec.errors {
			b.WriteString(fmt.Sprintf("[%d] %s", i+1, err.Format()))
			if i < len(ec.errors)-1 {
				b.WriteString("\n")
			}
		}
	}

	if len(ec.warnings) > 0 {
		if len(ec.errors) > 0 {
			b.WriteString("\n\n")
		}
		b.WriteString(fmt.Sprintf("Found %d warning(s):\n\n", len(ec.warnings)))
		for i, warn := range ec.warnings {
			b.WriteString(fmt.Sprintf("[%d] %s", i+1, warn.Format()))
			if i < len(ec.warnings)-1 {
				b.WriteString("\n")
			}
		}
	}

	return b.String()
}

func (ec *ErrorCollector) getSourceLine(lineNum int) string {
	if ec.source == "" || lineNum < 1 {
		return ""
	}

	lines := strings.Split(ec.source, "\n")
	if lineNum > len(lines) {
		return ""
	}

	return lines[lineNum-1]
}

func UnexpectedTokenError(expected, got lexer.Token) *NBXError {
	return &NBXError{
		Severity:   SeverityError,
		Message:    fmt.Sprintf("Expected %s, but got '%s'", tokenTypeToString(expected.Type), got.Literal),
		Line:       got.Line,
		Column:     got.Column,
		Token:      &got,
		Suggestion: fmt.Sprintf("Try replacing '%s' with %s", got.Literal, tokenTypeToString(expected.Type)),
	}
}

func UndefinedVariableError(varName string, line, column int, availableVars []string) *NBXError {
	err := &NBXError{
		Severity: SeverityError,
		Message:  fmt.Sprintf("Undefined variable '%s'", varName),
		Line:     line,
		Column:   column,
	}

	if len(availableVars) > 0 {
		similar := findSimilar(varName, availableVars)
		if len(similar) > 0 {
			err.Suggestion = fmt.Sprintf("Did you mean '%s'?", similar[0])
		} else {
			err.RelatedInfo = []string{
				fmt.Sprintf("Available variables: %s", strings.Join(availableVars, ", ")),
			}
		}
	} else {
		err.Suggestion = "No variables have been declared yet. Use 'var name: TYPE = value' to declare a variable"
	}

	return err
}

func TypeMismatchError(expected, got string, line, column int) *NBXError {
	return &NBXError{
		Severity:   SeverityError,
		Message:    fmt.Sprintf("Type mismatch: expected %s, got %s", expected, got),
		Line:       line,
		Column:     column,
		Suggestion: fmt.Sprintf("Convert the value to %s or change the variable type to %s", expected, got),
	}
}

func DuplicateDeclarationError(name string, line, column, firstLine int) *NBXError {
	return &NBXError{
		Severity:   SeverityError,
		Message:    fmt.Sprintf("Duplicate declaration of '%s'", name),
		Line:       line,
		Column:     column,
		Suggestion: fmt.Sprintf("Variable '%s' is already declared", name),
		RelatedInfo: []string{
			fmt.Sprintf("First declaration at line %d", firstLine),
		},
	}
}

func UnknownAttributeError(attrName, context string, line, column int, validAttrs []string) *NBXError {
	err := &NBXError{
		Severity: SeverityError,
		Message:  fmt.Sprintf("Unknown attribute '%s' in %s", attrName, context),
		Line:     line,
		Column:   column,
	}

	if len(validAttrs) > 0 {
		similar := findSimilar(attrName, validAttrs)
		if len(similar) > 0 {
			err.Suggestion = fmt.Sprintf("Did you mean '%s'?", similar[0])
		}
		err.RelatedInfo = []string{
			fmt.Sprintf("Valid attributes for %s: %s", context, strings.Join(validAttrs, ", ")),
		}
	}

	return err
}

func tokenTypeToString(tokenType lexer.TokenType) string {
	switch tokenType {
	case lexer.TOKEN_IDENT:
		return "an identifier"
	case lexer.TOKEN_STRING:
		return "a string"
	case lexer.TOKEN_ASSIGN:
		return "'='"
	case lexer.TOKEN_COLON:
		return "':'"
	case lexer.TOKEN_COMMA:
		return "','"
	case lexer.TOKEN_DOT:
		return "'.'"
	case lexer.TOKEN_LPAREN:
		return "'('"
	case lexer.TOKEN_RPAREN:
		return "')'"
	case lexer.TOKEN_LBRACE:
		return "'{'"
	case lexer.TOKEN_RBRACE:
		return "'}'"
	case lexer.TOKEN_KEYWORD:
		return "a keyword"
	default:
		return "a token"
	}
}

func findSimilar(target string, candidates []string) []string {
	type scoredCandidate struct {
		name  string
		score int
	}

	scored := make([]scoredCandidate, 0)
	for _, candidate := range candidates {
		score := levenshtein(strings.ToLower(target), strings.ToLower(candidate))
		if score <= 3 {
			scored = append(scored, scoredCandidate{candidate, score})
		}
	}

	if len(scored) == 0 {
		return nil
	}

	for i := 0; i < len(scored)-1; i++ {
		for j := i + 1; j < len(scored); j++ {
			if scored[j].score < scored[i].score {
				scored[i], scored[j] = scored[j], scored[i]
			}
		}
	}

	result := make([]string, 0, 3)
	for i := 0; i < len(scored) && i < 3; i++ {
		result = append(result, scored[i].name)
	}

	return result
}

func levenshtein(s1, s2 string) int {
	if len(s1) == 0 {
		return len(s2)
	}
	if len(s2) == 0 {
		return len(s1)
	}

	matrix := make([][]int, len(s1)+1)
	for i := range matrix {
		matrix[i] = make([]int, len(s2)+1)
		matrix[i][0] = i
	}
	for j := range matrix[0] {
		matrix[0][j] = j
	}

	for i := 1; i <= len(s1); i++ {
		for j := 1; j <= len(s2); j++ {
			cost := 1
			if s1[i-1] == s2[j-1] {
				cost = 0
			}

			matrix[i][j] = min(
				matrix[i-1][j]+1,      // deletion
				matrix[i][j-1]+1,      // insertion
				matrix[i-1][j-1]+cost, // substitution
			)
		}
	}

	return matrix[len(s1)][len(s2)]
}
